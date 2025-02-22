package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"rental/middleware"
	"rental/models"
	"github.com/beego/beego/v2/client/orm"
	
)

// LocationService handles external API requests
type LocationService struct {
	ApiBaseUrl string
	ApiKey     string
	HttpClient *middleware.InstrumentedHttpClient // Reused client for proper monitoring
}

// FilteredLocation represents the structure for filtered locations.
type FilteredLocation struct {
	DestId   string `json:"dest_id"`
	DestType string `json:"dest_type"`
	Value    string `json:"value"`
}

// GetLocations fetches locations from the external API.
func (service *LocationService) GetLocations(query string) ([]FilteredLocation, error) {
	url := fmt.Sprintf("%s/web/stays/auto-complete?query=%s", service.ApiBaseUrl, query)

	// Using the instrumented HTTP client with a timeout
	client := service.HttpClient

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("x-rapidapi-host", "booking-com18.p.rapidapi.com")
	req.Header.Add("x-rapidapi-key", service.ApiKey)

	// Use the instrumented client to make the request
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		log.Printf("API request failed: status=%d, response=%s", resp.StatusCode, bodyString)
		return nil, fmt.Errorf("API request failed with status code: %d, Response: %s", resp.StatusCode, bodyString)
	}

	var apiResponse map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&apiResponse); err != nil {
		log.Printf("Error decoding response: %v", err)
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	var filteredLocations []FilteredLocation
	if data, ok := apiResponse["data"].([]interface{}); ok {
		o := orm.NewOrm()
		for _, item := range data {
			if itemMap, ok := item.(map[string]interface{}); ok {
				filteredLocation := FilteredLocation{
					DestId:   itemMap["dest_id"].(string),
					DestType: itemMap["dest_type"].(string),
					Value:    itemMap["label"].(string),
				}
				filteredLocations = append(filteredLocations, filteredLocation)

				// Insert into the database if the location doesn't already exist
				location := &models.Location{
					DestId:   itemMap["dest_id"].(string),
					DestType: itemMap["dest_type"].(string),
					Value:    itemMap["label"].(string),
				}
				existingLocation := models.Location{DestId: location.DestId}
				err := o.Read(&existingLocation, "DestId")
				if err != nil {
					_, err := o.Insert(location)
					if err != nil {
						log.Printf("Error inserting location: %v", err)
					}
				}
			}
		}
	}

	return filteredLocations, nil
}
