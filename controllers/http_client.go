package controllers

import (
	"time"

	"rental/middleware"
)
var HttpClient = middleware.NewInstrumentedHttpClient(10 * time.Second)
// Create a global instance of the InstrumentedHttpClient
