package controllers

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "time"
    "rental/middleware"

    "github.com/beego/beego/v2/client/orm"
    "github.com/beego/beego/v2/server/web"
    "github.com/beego/beego/v2/server/web/context"
    "rental/models"
)

type StayData struct {
    web.Controller
    apiBaseUrl string
    apiKey     string
}
type HotelDetails struct {
    IDHotel     int     `json:"id_hotel"`
    HotelID     string  `json:"hotel_id"`
    HotelName   string  `json:"hotel_name"`
    DestID      string  `json:"dest_id"`
    Location    string  `json:"location"`
    Rating      float64 `json:"rating"`
    ReviewCount int     `json:"review_count"`
    Price       string  `json:"price"`
    NumBeds     int     `json:"num_beds"`
    NumBedR     int     `json:"num_bedrooms"`
    NumBaths    int     `json:"num_bathrooms"`
}
type APIResponse struct {
    Data struct {
        Results []struct {
            ID               string `json:"id"`
            BasicPropertyData struct {
                ID      int `json:"id"`
                Reviews struct {
                    TotalScore   float64 `json:"totalScore"`
                    ReviewsCount int     `json:"reviewsCount"`
                } `json:"reviews"`
                
            } `json:"basicPropertyData"`
            Location struct {
                DisplayLocation string `json:"displayLocation"`
            } `json:"location"`
            DisplayName struct {
                Text string `json:"text"`
            } `json:"displayName"`
            MatchingUnitConfigurations struct {
                CommonConfiguration struct {
                    NbAllBeds   int `json:"nbAllBeds"`
                    NbAllBedR   int `json:"nbBedrooms"`
                    NbBathrooms int `json:"nbBathrooms"`
                } `json:"commonConfiguration"`
            } `json:"matchingUnitConfigurations"`
            PriceDisplayInfoIrene struct {
                DisplayPrice struct {
                    AmountPerStay struct {
                        Amount string `json:"amount"`
                    } `json:"amountPerStay"`
                } `json:"displayPrice"`
            } `json:"priceDisplayInfoIrene"`
        } `json:"results"`
    } `json:"data"`
    Meta struct {
        CurrentPage   int `json:"currentPage"`
        Limit        int `json:"limit"`
        TotalRecords int `json:"totalRecords"`
        TotalPage    int `json:"totalPage"`
    } `json:"meta"`
    Status  bool   `json:"status"`
    Message string `json:"message"`
}

// Init function to initialize apiBaseUrl and apiKey
func (c *StayData) Init(ctx *context.Context, controllerName, actionName string, app interface{}) {
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

// FetchStaysData handles GET requests for "/fetch_stays_data"
func (c *StayData) Get() {
	o := orm.NewOrm()
	var locations []models.Location
	_, err := o.QueryTable("location").All(&locations)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Failed to retrieve locations"}
		c.ServeJSON()
		return
	}

    var allHotels []HotelDetails
    for _, loc := range locations {
        // Fetch data for each location
        hotels, err := getHotelData(loc.DestId, loc.DestType, c.apiBaseUrl, c.apiKey)
        if err != nil {
            log.Printf("Error fetching data for destID: %s, destType: %s - %v", loc.DestId, loc.DestType, err)
            continue
        }

        allHotels = append(allHotels, hotels...)
    }

    if len(allHotels) == 0 {
        log.Print("No hotels found for any location")
        c.Data["json"] = map[string]interface{}{
            "message": "No stays found for any location.",
        }
        c.ServeJSON()
        return
    }

    c.Data["json"] = map[string]interface{}{
        "message": "Stays data fetched successfully",
        "results": allHotels,
    }
    c.ServeJSON()
}

func getHotelData(destID, destType, apiBaseUrl, apiKey string) ([]HotelDetails, error) {
	if destID == "" || destType == "" {
		return nil, fmt.Errorf("destID and destType cannot be empty")
	}

	url := fmt.Sprintf("%s/web/stays/search?destId=%s&destType=%s&checkIn=2025-03-12&checkOut=2025-03-18", apiBaseUrl, destID, destType)
	log.Printf("Making request to URL: %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("x-rapidapi-host", "booking-com18.p.rapidapi.com")
	req.Header.Add("x-rapidapi-key", apiKey)

	// Use the instrumented HTTP client for monitoring
	client := middleware.NewInstrumentedHttpClient(10 * time.Second) // 10 seconds timeout for example
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	log.Printf("API Response Status: %s", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	log.Printf("Raw API Response: %s", string(body))

	var apiResponse APIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		var rawResponse map[string]interface{}
		if jsonErr := json.Unmarshal(body, &rawResponse); jsonErr == nil {
			log.Printf("API Response Structure: %+v", rawResponse)
		}
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	if !apiResponse.Status {
		return nil, fmt.Errorf("API returned error: %s", apiResponse.Message)
	}

	if len(apiResponse.Data.Results) == 0 {
		log.Printf("API returned success but no results. Total records: %d", apiResponse.Meta.TotalRecords)
		return nil, nil
	}

	// Process the results
	var hotels []HotelDetails
	o := orm.NewOrm()

	for _, result := range apiResponse.Data.Results {
		hotel := HotelDetails{
			IDHotel:     result.BasicPropertyData.ID,
			HotelID:     result.ID,
			HotelName:   result.DisplayName.Text,
			DestID:      destID,
			Location:    result.Location.DisplayLocation,
			Rating:      result.BasicPropertyData.Reviews.TotalScore,
			ReviewCount: result.BasicPropertyData.Reviews.ReviewsCount,
			Price:       result.PriceDisplayInfoIrene.DisplayPrice.AmountPerStay.Amount,
			NumBeds:     result.MatchingUnitConfigurations.CommonConfiguration.NbAllBeds,
			NumBedR:     result.MatchingUnitConfigurations.CommonConfiguration.NbAllBedR,
			NumBaths:    result.MatchingUnitConfigurations.CommonConfiguration.NbBathrooms,
		}

		hotels = append(hotels, hotel)

		// Check if hotel exists in DB
		existingRentalProperty := models.RentalProperty{IDHotel: hotel.IDHotel}
		err := o.Read(&existingRentalProperty, "IDHotel")
		if err == nil {
			log.Printf("Hotel with IDHotel %d already exists, skipping insertion.", hotel.IDHotel)
			continue
		}

		// Insert rental property into DB
		rentalProperty := models.RentalProperty{
			IDHotel:     hotel.IDHotel,
			HotelID:     hotel.HotelID,
			HotelName:   hotel.HotelName,
			DestID:      hotel.DestID,
			Location:    hotel.Location,
			Rating:      hotel.Rating,
			ReviewCount: hotel.ReviewCount,
			Price:       hotel.Price,
			NumBeds:     hotel.NumBeds,
			NumBedR:     hotel.NumBedR,
			NumBaths:    hotel.NumBaths,
		}

		_, err = o.Insert(&rentalProperty)
		if err != nil {
			log.Printf("Error inserting rental property: %v", err)
		} else {
			log.Printf("Inserted rental property with ID: %d", rentalProperty.ID)
		}
	}

	// Limit the results to 25 (in case the API doesn't support the limit parameter)
	if len(hotels) > 2 {
		hotels = hotels[:2]
	}

	log.Printf("Processed %d rental properties successfully", len(hotels))
	return hotels, nil
}
