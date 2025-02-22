// db.go

package controllers

import (
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// InitDb initializes the database connection
func InitDb() {
	// Register the PostgreSQL driver
	orm.RegisterDriver("postgres", orm.DRPostgres)

	// Set up the database connection
	orm.RegisterDataBase("default", "postgres", "user=postgres password=postgres host=localhost port=5432 dbname=booking_db sslmode=disable")
	
	// Synchronize the database
	orm.RunSyncdb("default", false, true)
}
