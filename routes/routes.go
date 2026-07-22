//routes main work is to decide which function should handle the request 
//number of routes is the number of APIs you should create in the routes.go file. 
//The number of functions created in the controller file should be equal to the number of routes created in the routes.go file.
//we are here creating the end points for the APIs.
//PUT changes everything in the resource, PATCH changes only the specified fields in the resource
//DELETE removes the resource from the server.

package routes

import (
	"MoneyTransfer/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine){
	
	//ACCOUNT API
	router.POST("/accounts",controllers.CreateAccount) 
	router.GET("/accounts", controllers.GetAccounts)    //gives all the accounts 
	router.GET("/accounts/:id", controllers.GetAccountByID)    //give the account with the given id
	router.PATCH("/accounts/:id",controllers.UpdateAccount)
    router.DELETE("/accounts/:id", controllers.DeleteAccount)
	
	
	
	
	//TRANSFER API
	router.POST("/transfers",controllers.CreateTransfer)
	router.GET("/transfers", controllers.GetTransfer)      
	router.GET("/transfers/:id", controllers.GetTransferByID)     

	//USER APR

	router.POST("/users", controllers.CreateUser)
	router.GET("/users", controllers.GetUsers)
	router.GET("/users/:id", controllers.GetUserByID)
	router.PATCH("/users/:id", controllers.UpdateUser)
	router.DELETE("/users/:id", controllers.DeleteUser)

	//COMPANY API
	router.POST("/companies", controllers.CreateCompany)
	router.GET("/companies", controllers.GetCompanies)
	router.GET("/companies/:id", controllers.GetCompanyByID)
	router.PATCH("/companies/:id", controllers.UpdateCompany)
	router.DELETE("/companies/:id", controllers.DeleteCompany)

}


