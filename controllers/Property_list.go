package controllers

import (
	"github.com/beego/beego/v2/server/web"
	
	
    "rental/services"
	
)

type PropertyController struct {
	web.Controller
}

// ListProperties handles the GET request for property listings
func (c *PropertyController) Get() {
    location := c.GetString("location", "New York")

    // Fetch properties using the service layer
	enhancedProperties, err := services.GetProperties(location)
	if err != nil {
		c.Data["json"] = map[string]string{
			"error":   "Failed to retrieve properties",
			"details": err.Error(),
		}
		c.ServeJSON()
		return
	}

	// Pass data to the template
	c.Data["json"] = map[string]interface{}{
		"properties": enhancedProperties,
	}
	c.ServeJSON()
}