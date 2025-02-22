package controllers

import (
	"rental/services"

	"github.com/beego/beego/v2/server/web"
)

type PropertyDetailsController struct {
	web.Controller
}

// Get handles GET requests for /v1/property/details
func (c *PropertyDetailsController) Get() {
	// Fetch property details using the service layer
	response, err := services.GetPropertyDetails()
	if err != nil {
		c.Data["json"] = map[string]string{
			"error":   "Failed to retrieve property details",
			"details": err.Error(),
		}
		c.ServeJSON()
		return
	}

	// Send the JSON response
	c.Data["json"] = response
	c.ServeJSON()
}
