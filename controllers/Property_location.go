package controllers

import (
    "strings"
	"strconv"
    "github.com/beego/beego/v2/server/web"
    "github.com/beego/beego/v2/client/orm"
    "rental/models"
)

type PropertyLocationController struct {
    web.Controller
}

// GetPropertiesByLocation handles the GET request for property listings by location (dest_id)
func (c *PropertyLocationController) Get() {
    // Initialize ORM
    o := orm.NewOrm()
    qs := o.QueryTable("rental_property")

    // Prepare result slice
    var properties []models.RentalProperty
    _, err := qs.All(&properties)

    if err != nil {
        c.Data["json"] = map[string]string{
            "error": "Failed to retrieve properties",
            "details": err.Error(),
        }
        c.ServeJSON()
        return
    }

    // Prepare a map to group properties by LocationId
    locationMap := make(map[string]map[string]interface{})

    for _, prop := range properties {
        var details models.PropertyDetails
        err := o.QueryTable("property_details").Filter("HotelID", prop.HotelID).One(&details)

        // Handle the error for property details query
        if err != nil {
            // If there's an error, use default/empty values for details
            details = models.PropertyDetails{
                ImageUrl1:  "", // default empty image URL
                ImageUrl2:  "",
                ImageUrl3:  "",
                ImageUrl4:  "",
                ImageUrl5:  "",
                Description: "No description available",
            }
        }

        // Fetch the location value from the Location table
        var location models.Location
        err = o.QueryTable("location").Filter("DestId", prop.DestID).One(&location)
        if err != nil {
            // If there's an error, use the Location from the property
            location.Value = prop.Location
        }

        propertyDetail := map[string]interface{}{
			"ID":           prop.ID,
			"IDHotel":      prop.IDHotel,
			"HotelID":      prop.HotelID,
			"HotelName":    prop.HotelName,
			"Location":     prop.Location,  // Use the value from the Location table
			"Price":        prop.Price,
			"PropertyType": prop.PropertyType,
			"Guests":       prop.Guests,
			"Rating":       prop.Rating,
			"ReviewCount":  prop.ReviewCount,
			"NumBeds":      prop.NumBeds,
			"NumBedR":      prop.NumBedR,
			"NumBaths":     prop.NumBaths,
			"Bedroom":      prop.Bedroom,
			"Amenities":    strings.Split(prop.Amenities, ","),
			"HeroImage":    details.ImageUrl1,
			"Gallery":      []string{details.ImageUrl1, details.ImageUrl2, details.ImageUrl3, details.ImageUrl4, details.ImageUrl5},
		}
	
		if _, exists := locationMap[prop.DestID]; !exists {
			locationMap[prop.DestID] = map[string]interface{}{
				"LocationId": prop.DestID,
				"Value":      location.Value,
			}
		}
	
		locationMap[prop.DestID]["id:"+strconv.Itoa(prop.IDHotel)] = propertyDetail
	}

    // Convert map to slice to match required JSON format
    var locations []map[string]interface{}
    for _, loc := range locationMap {
        locations = append(locations, loc)
    }

    // Pass data to the template
    c.Data["json"] = locations
    c.ServeJSON()
}