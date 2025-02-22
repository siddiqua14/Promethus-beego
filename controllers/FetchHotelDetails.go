
package controllers

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "strings"
    "net/http"

    "github.com/beego/beego/v2/client/orm"
    "github.com/beego/beego/v2/server/web"
    "github.com/beego/beego/v2/server/web/context"
    "rental/models"
)

type FetchHotelDetails struct {
    web.Controller
    apiBaseUrl string
    apiKey     string
}

// Init function to initialize apiBaseUrl and apiKey
func (c *FetchHotelDetails) Init(ctx *context.Context, controllerName, actionName string, app interface{}) {
    c.Controller.Init(ctx, controllerName, actionName, app)

    apiBaseUrl, err := web.AppConfig.String("API_BASE_URL")
    if err != nil {
        c.Data["json"] = map[string]interface{}{"error": fmt.Sprintf("Error retrieving API Base URL: %v", err)}
        c.ServeJSON()
        return
    }

    apiKey, err := web.AppConfig.String("API_KEY")
    if err != nil {
        c.Data["json"] = map[string]interface{}{"error": fmt.Sprintf("Error retrieving API Key: %v", err)}
        c.ServeJSON()
        return
    }

    c.apiBaseUrl = apiBaseUrl
    c.apiKey = apiKey
}

func (c *FetchHotelDetails) Get() {
    o := orm.NewOrm()

    // Fetch top 25 records from the rental_property table
    var rentalProperties []models.RentalProperty
    _, err := o.QueryTable("rental_property").Limit(80).All(&rentalProperties)
    if err != nil {
        c.Data["json"] = map[string]interface{}{"error": "Failed to retrieve rental properties"}
        c.ServeJSON()
        return
    }

    var hotelDetailsList []map[string]interface{}

    for _, rentalProperty := range rentalProperties {
        hotelID := rentalProperty.IDHotel
        url := fmt.Sprintf("https://booking-com18.p.rapidapi.com/stays/detail?hotelId=%d&checkinDate=2025-02-12&checkoutDate=2025-02-18&units=metric", hotelID)
        log.Printf("Making request to URL: %s", url)

        req, err := http.NewRequest("GET", url, nil)
        if err != nil {
            log.Printf("Error creating request for hotel ID %d: %v", hotelID, err)
            continue
        }

        req.Header.Add("x-rapidapi-host", "booking-com18.p.rapidapi.com")
        req.Header.Add("x-rapidapi-key", c.apiKey)

        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
            log.Printf("Error making request for hotel ID %d: %v", hotelID, err)
            continue
        }

        // Ensure the response body is closed after reading
        body, err := io.ReadAll(resp.Body)
        resp.Body.Close()  // Close response body immediately after reading
        if err != nil {
            log.Printf("Error reading response body for hotel ID %d: %v", hotelID, err)
            continue
        }

        log.Printf("API Response Status for hotel ID %d: %s", hotelID, resp.Status)
        log.Printf("Raw API Response for hotel ID %d: %s", hotelID, string(body))

        var apiResponse map[string]interface{}
        if err := json.Unmarshal(body, &apiResponse); err != nil {
            log.Printf("Error unmarshalling response for hotel ID %d: %v", hotelID, err)
            continue
        }

        data, ok := apiResponse["data"].(map[string]interface{})
        if !ok {
            log.Printf("Invalid data format in response for hotel ID %d", hotelID)
            continue
        }

        availableRooms, ok := data["available_rooms"].(float64)
        if !ok {
            log.Printf("Invalid available_rooms format in response for hotel ID %d", hotelID)
            continue
        }

        roomRecommendations, ok := data["room_recommendation"].([]interface{})
        if !ok {
            log.Printf("Invalid room_recommendation format in response for hotel ID %d", hotelID)
            continue
        }

        var totalAdults, totalChildren float64
        for _, recommendation := range roomRecommendations {
            if recommendationMap, ok := recommendation.(map[string]interface{}); ok {
                if adults, ok := recommendationMap["adults"].(float64); ok {
                    totalAdults += adults
                }
                if children, ok := recommendationMap["children"].(float64); ok {
                    totalChildren += children
                }
            }
        }

        accommodationTypeName, ok := data["accommodation_type_name"].(string)
        if !ok {
            log.Printf("Invalid accommodation_type_name format in response for hotel ID %d", hotelID)
            continue
        }

        facilitiesBlock, ok := data["facilities_block"].(map[string]interface{})
        if !ok {
            log.Printf("Invalid facilities_block format in response for hotel ID %d", hotelID)
            continue
        }

        facilities, ok := facilitiesBlock["facilities"].([]interface{})
        if !ok {
            log.Printf("Invalid facilities format in response for hotel ID %d", hotelID)
            continue
        }

        var amenities []string
        for i, facility := range facilities {
            if i >= 3 {
                break
            }
            facilityMap, ok := facility.(map[string]interface{})
            if !ok {
                log.Printf("Invalid facility format in response for hotel ID %d", hotelID)
                continue
            }
            name, ok := facilityMap["name"].(string)
            if !ok {
                log.Printf("Invalid facility name format in response for hotel ID %d", hotelID)
                continue
            }
            amenities = append(amenities, name)
        }
        
        // Parse CityInTrans
        cityInTrans, ok := data["city_in_trans"].(string)
        if !ok {
            cityInTrans = ""
        }

        
        // Update the existing rental property with additional details
        rentalProperty.Bedroom = int(availableRooms)
        rentalProperty.Guests = int(totalAdults + totalChildren)
        rentalProperty.PropertyType = accommodationTypeName
        rentalProperty.Amenities = strings.Join(amenities, ", ")
        rentalProperty.CityInTrans = cityInTrans

        if _, err := o.Update(&rentalProperty); err != nil {
            log.Printf("Error updating rental property for hotel ID %d: %v", hotelID, err)
            continue
        }

        log.Printf("Rental property updated successfully for hotel ID %d", hotelID)

        // Collect details for JSON output
        hotelDetails := map[string]interface{}{
            "IDHotel":       rentalProperty.IDHotel,
            "HotelID":       rentalProperty.HotelID,
            "HotelName":     rentalProperty.HotelName,
            "DestID":        rentalProperty.DestID,
            "Location":      rentalProperty.Location,
            "Rating":        rentalProperty.Rating,
            "ReviewCount":   rentalProperty.ReviewCount,
            "Price":         rentalProperty.Price,
            "NumBeds":       rentalProperty.NumBeds,
            "NumBedR":       rentalProperty.NumBedR,
            "NumBaths":      rentalProperty.NumBaths,
            "Bedroom":       rentalProperty.Bedroom,
            "Guests":        rentalProperty.Guests,
            "PropertyType":  rentalProperty.PropertyType,
            "Amenities":     rentalProperty.Amenities,
            "CityInTrans":   rentalProperty.CityInTrans,
        }
        hotelDetailsList = append(hotelDetailsList, hotelDetails)
    }

    // Send the hotel details as JSON response
    c.Data["json"] = map[string]interface{}{
        "message": "Rental properties updated successfully",
        "details": hotelDetailsList,
    }
    c.ServeJSON()
}

