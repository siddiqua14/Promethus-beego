// models/models.go
package models

import (
    "github.com/beego/beego/v2/client/orm"
)

type Location struct {
    Id            int64  `orm:"auto;pk"`
    DestId        string `orm:"size(100);column(dest_id)"`
    DestType      string `orm:"size(100);column(dest_type)"`
    Value         string `orm:"size(255);column(value)"`
}

type RentalProperty struct {
    ID            int     `orm:"auto"`
    IDHotel       int     `orm:"index"`
    HotelID       string  `orm:"size(100)"`
    HotelName     string  `orm:"size(255)"`
    DestID        string  `orm:"size(100)"`
    Location      string  `orm:"size(255)"`
    Rating        float64 `orm:"digits(2);decimals(1)"`
    ReviewCount   int
    Price         string  `orm:"size(100)"`
    NumBeds       int
    NumBedR       int
    NumBaths      int
    Bedroom       int
    Guests        int
    PropertyType  string  `orm:"size(100)"`
    Amenities     string  `orm:"type(text)"`
    CityInTrans   string `orm:"column(city_in_trans)"`
}

type PropertyDetails struct {
    ID          int    `orm:"auto"`
    HotelID     string `orm:"size(100);index"`
    Description string `orm:"type(text)"`
    CityInTrans string `orm:"size(255);column(city_in_trans)"`
    ImageUrl1   string `orm:"size(500);column(image_url_1)"`
    ImageUrl2   string `orm:"size(500);column(image_url_2)"`
    ImageUrl3   string `orm:"size(500);column(image_url_3)"`
    ImageUrl4   string `orm:"size(500);column(image_url_4)"`
    ImageUrl5   string `orm:"size(500);column(image_url_5)"`
}
func init() {
    orm.RegisterModel(new(Location), new(RentalProperty), new(PropertyDetails))
}