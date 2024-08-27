package main

import (
	"loan-tracker/config"
	routes "loan-tracker/delivery/route"
)

func main() {

	client := config.Connect()
	db := client.Database("loan_tracker")
	r := routes.SetupUserRouter(db)
	r = routes.SetupLoanRouter(db, r)
	r = routes.SetupLogRouter(db, r)
	r.Run(":8080")
}
