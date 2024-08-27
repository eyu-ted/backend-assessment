// main.go
package main

import (
	"context"
	routes "loan-tracker/delivery/route"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	clientOptions := options.Client().ApplyURI("mongodb+srv://eyutedd:eyu81619@cluster0.o7vkwdm.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0")

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("loan_tracker")
	r := routes.SetupUserRouter(db)
	

	r.Run(":8080")
}
