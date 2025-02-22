package main

import (
	_ "rental/routers"
    "rental/controllers" 

	"github.com/beego/beego/v2/server/web"
    
)


func main() {
    // Initialize the database connection
	controllers.InitDb()

    // Start the Beego server
    web.BConfig.WebConfig.ViewsPath = "views"
    web.Run()
}
