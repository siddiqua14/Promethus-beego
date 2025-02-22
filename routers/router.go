
package routers

import (
	"rental/controllers"
	"rental/middleware"


	"github.com/beego/beego/v2/server/web"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func init() {
// Initialize Prometheus metrics **only once**
middleware.InitPrometheusMetrics()

// Register Prometheus metrics handler
web.Handler("/metrics", promhttp.Handler())

    web.Router("/fetch_locations", &controllers.LocationController{})
	web.Router("/fetch_stays_data", &controllers.StayData{})
    //web.Router("/fetch-hotel-details", &controllers.FetchHotelDetails{})
    //web.Router("/fetch-hotel-images-and-description", &controllers.FetchHotelImagesAndDescriptions{})
    // Page route
    // Property listing endpoint
    ns := web.NewNamespace("/v1",
        // Group routes under the /property namespace
        web.NSNamespace("/property",
            web.NSRouter("/list", &controllers.PropertyController{}),
            web.NSRouter("/details", &controllers.PropertyDetailsController{}),
            web.NSRouter("/location", &controllers.PropertyLocationController{}),
        ),
)

web.AddNamespace(ns)
}