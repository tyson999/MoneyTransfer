//thiss is the starting point of the program.
//It is the main file that will be executed when the program is run.
//It will create a new router and set up the routes for the APIs.
//It will then start the server on port 8080.
//go first searches for the main package and then finds the main function and then executes the code inside the main function.

package main

import (
	"MoneyTransfer/database"
	"MoneyTransfer/routes"
	"log"
    "github.com/gin-gonic/gin"
)

func main() {
	if err := database.Connect(); err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	router := gin.Default()

	routes.SetupRoutes(router)

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
