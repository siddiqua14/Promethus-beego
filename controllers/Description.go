package controllers

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"

    "rental/models"
    "github.com/beego/beego/v2/client/orm"
    "github.com/beego/beego/v2/server/web"
    
)
type FetchHotelImagesAndDescriptions struct {
    web.Controller
    apiBaseUrl string
    apiKey     string
}

// FetchHotelImagesAndDescriptions fetches images and descriptions for hotels
func (c *FetchHotelImagesAndDescriptions) Get() {
    o := orm.NewOrm()

    // Retrieve all RentalProperty entries (limit to 20)
    var rentalProperties []models.RentalProperty
    _, err := o.QueryTable(new(models.RentalProperty)).All(&rentalProperties)
    if err != nil {
        log.Printf("Error fetching rental properties: %v", err)
        c.Data["json"] = map[string]string{"error": "Failed to fetch rental properties"}
        c.ServeJSON()
        return
    }

    // Process each rental property
    for _, rentalProperty := range rentalProperties {
        hotelID := rentalProperty.HotelID
        url := fmt.Sprintf("%s/web/stays/details?id=%s&checkIn=2025-02-12&checkOut=2025-02-18", c.apiBaseUrl, hotelID)
        log.Printf("Making request to URL: %s", url)

        req, err := http.NewRequest("GET", url, nil)
        if err != nil {
            log.Printf("Error creating request for HotelID %s: %v", hotelID, err)
            continue
        }

        req.Header.Add("x-rapidapi-host", "booking-com18.p.rapidapi.com")
        req.Header.Add("x-rapidapi-key", c.apiKey)

        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
            log.Printf("Error making request for HotelID %s: %v", hotelID, err)
            continue
        }
        defer resp.Body.Close()

        log.Printf("API Response Status for HotelID %s: %s", hotelID, resp.Status)

        body, err := io.ReadAll(resp.Body)
        if err != nil {
            log.Printf("Error reading response body for HotelID %s: %v", hotelID, err)
            continue
        }

        // Log the raw response for debugging
        log.Printf("Raw API Response for HotelID %s: %s", hotelID, string(body))

        var apiResponse map[string]interface{}
        if err := json.Unmarshal(body, &apiResponse); err != nil {
            log.Printf("Error unmarshalling response for HotelID %s: %v", hotelID, err)
            continue
        }

        // Check if 'data' exists in the response
        data, ok := apiResponse["data"].(map[string]interface{})
        if !ok {
            log.Printf("Missing data in response for HotelID %s", hotelID)
            continue
        }

        // Parse hotelPhotos
        hotelPhotos, ok := data["hotelPhotos"].([]interface{})
        if !ok {
            log.Printf("Missing or invalid hotelPhotos for HotelID %s", hotelID)
            continue
        }

        var images []string
        for _, photo := range hotelPhotos {
            photoMap, ok := photo.(map[string]interface{})
            if ok {
                if thumbURL, ok := photoMap["thumb_url"].(string); ok {
                    images = append(images, thumbURL)
                }
            }
        }

        // Parse hotelTranslation
        hotelTranslations, ok := data["hotelTranslation"].([]interface{})
        if !ok || len(hotelTranslations) == 0 {
            log.Printf("Missing or invalid hotelTranslation for HotelID %s", hotelID)
            continue
        }

        description := ""
        if descMap, ok := hotelTranslations[0].(map[string]interface{}); ok {
            if desc, ok := descMap["description"].(string); ok {
                description = desc
            }
        }

        // Parse CityInTrans
        cityInTrans, ok := data["cityInTrans"].(string)
        if !ok {
            cityInTrans = ""
        }

        // Insert or update PropertyDetails table
        propertyDetails := models.PropertyDetails{}
        if err := o.QueryTable(new(models.PropertyDetails)).Filter("HotelID", hotelID).One(&propertyDetails); err == orm.ErrNoRows {
            // No entry exists, so insert a new record
            propertyDetails = models.PropertyDetails{
                HotelID:     hotelID,
                Description: description,
                CityInTrans: cityInTrans,
                ImageUrl1:   getImageUrl(images, 0),
                ImageUrl2:   getImageUrl(images, 1),
                ImageUrl3:   getImageUrl(images, 2),
                ImageUrl4:   getImageUrl(images, 3),
                ImageUrl5:   getImageUrl(images, 4),
            }
            if _, err := o.Insert(&propertyDetails); err != nil {
                log.Printf("Error inserting PropertyDetails for HotelID %s: %v", hotelID, err)
                continue
            }
            log.Printf("PropertyDetails inserted successfully for HotelID %s", hotelID)
        } else if err == nil {
            // Entry exists, so update it
            propertyDetails.Description = description
            propertyDetails.CityInTrans = cityInTrans
            propertyDetails.ImageUrl1 = getImageUrl(images, 0)
            propertyDetails.ImageUrl2 = getImageUrl(images, 1)
            propertyDetails.ImageUrl3 = getImageUrl(images, 2)
            propertyDetails.ImageUrl4 = getImageUrl(images, 3)
            propertyDetails.ImageUrl5 = getImageUrl(images, 4)
            if _, err := o.Update(&propertyDetails); err != nil {
                log.Printf("Error updating PropertyDetails for HotelID %s: %v", hotelID, err)
                continue
            }
            log.Printf("PropertyDetails updated successfully for HotelID %s", hotelID)
        } else {
            log.Printf("Error querying PropertyDetails for HotelID %s: %v", hotelID, err)
            continue
        }
    }

    // Send response
    c.Data["json"] = map[string]string{"message": "Hotel details processed and stored successfully"}
    c.ServeJSON()
}

// getImageUrl safely retrieves an image URL from the list
func getImageUrl(images []string, index int) string {
    if index < len(images) {
        return images[index]
    }
    return ""
}