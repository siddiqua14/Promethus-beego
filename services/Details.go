package services

import (
	"github.com/beego/beego/v2/client/orm"
	"log"
	"rental/models"
)

// GetPropertyDetails retrieves detailed property data
func GetPropertyDetails() ([]map[string]interface{}, error) {
	o := orm.NewOrm()

	// Fetch all rental properties
	var rentalProperties []models.RentalProperty
	_, err := o.QueryTable("rental_property").All(&rentalProperties)
	if err != nil {
		log.Printf("Error retrieving rental properties: %v", err)
		return nil, err
	}

	// Prepare the response
	var response []map[string]interface{}
	for _, rentalProperty := range rentalProperties {
		// Fetch the corresponding location
		var location models.Location
		err := o.QueryTable("location").Filter("DestId", rentalProperty.DestID).One(&location)
		if err != nil {
			log.Printf("Error retrieving location for DestID %s: %v", rentalProperty.DestID, err)
			continue
		}

		// Fetch the corresponding property details
		var propertyDetails models.PropertyDetails
		err = o.QueryTable("property_details").Filter("HotelID", rentalProperty.HotelID).One(&propertyDetails)
		if err != nil {
			log.Printf("Error retrieving property details for HotelID %s: %v", rentalProperty.HotelID, err)
			continue
		}

		// Construct the property details
		propertyInfo := map[string]interface{}{
			"dest_id":     location.DestId,
			"dest_type":   location.DestType,
			"value":       location.Value,
			"IDHotel":     rentalProperty.IDHotel,
			"HotelID":     rentalProperty.HotelID,
			"Location":    rentalProperty.Location,
			"HotelName":   rentalProperty.HotelName,
			"PropertyType": rentalProperty.PropertyType,
			"Rating":      rentalProperty.Rating,
			"ReviewCount": rentalProperty.ReviewCount,
			"Bedrooms":    rentalProperty.Bedroom,
			"NumBaths":    rentalProperty.NumBaths,
			"NumBedR":     rentalProperty.NumBedR,
			"Guests":      rentalProperty.Guests,
			"Amenities":   rentalProperty.Amenities,
			"CityInTrans": rentalProperty.CityInTrans,
			"description": propertyDetails.Description,
			"images": []string{
				propertyDetails.ImageUrl1,
				propertyDetails.ImageUrl2,
				propertyDetails.ImageUrl3,
				propertyDetails.ImageUrl4,
				propertyDetails.ImageUrl5,
			},
		}

		response = append(response, propertyInfo)
	}

	return response, nil
}
