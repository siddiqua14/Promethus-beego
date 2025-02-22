package controllers

import (
	"fmt"
    "time"
	"rental/services"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"rental/middleware"
)

// LocationController handles incoming requests for location data
type LocationController struct {
	web.Controller
	locationService *services.LocationService
}

func (c *LocationController) Init(ctx *context.Context, controllerName, actionName string, app interface{}) {
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

	// Initialize the instrumented HTTP client with a timeout (10 seconds)
	httpClient := middleware.NewInstrumentedHttpClient(10 * time.Second)

	// Initialize the location service with the API configuration and the instrumented client
	c.locationService = &services.LocationService{
		ApiBaseUrl: apiBaseUrl,
		ApiKey:     apiKey,
		HttpClient: httpClient,
	}
}

func (c *LocationController) Get() {
	query := "New%20York"
	
	// Call the service method to get locations
	filteredLocations, err := c.locationService.GetLocations(query)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": err.Error()}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{
		"success": true,
		"query":   query,
		"count":   len(filteredLocations),
		"data": map[string]interface{}{
			"data": filteredLocations,
		},
	}
	c.ServeJSON()
}
