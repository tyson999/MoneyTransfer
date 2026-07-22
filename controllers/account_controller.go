//controller handles the http specific work(request and response)
//number of routes created should be equal to the number of functions created in the controller file
//the relation between routes and controllers is that the routes file defines the endpoints for the APIs, while the controller file contains the logic for handling requests to those endpoints. Each route corresponds to a specific function in the controller that processes the request and returns a response.
//service files are created to handle the business logic and to make our controller file more organized and maintainable. The service files can contain functions that handle the core business logic, while the controller file can focus on handling HTTP requests and responses. This separation of concerns can make the code easier to read, test, and maintain.

package controllers

import(
	"net/http"
	"MoneyTransfer/models"
	"github.com/gin-gonic/gin"
	"strconv"
)


var Accounts = []models.Account{}

func CreateAccount(c *gin.Context){
	var account models.Account

	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return 

	}

	account.ID = len(Accounts) +1 
	account.Status = "active"

	Accounts = append(Accounts, account)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Account created successfully",
		"account": account,
	})
}

func GetAccounts(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"accounts":Accounts,
	})
}

func GetAccountByID(c *gin.Context) {
	id,err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid account id",
		})
		return
	}

	for _, account := range Accounts {
		if account.ID == id{
			c.JSON(http.StatusOK, gin.H{
				"account": account,
			})

			return 
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"error": "Account not found",
	})
}

func UpdateAccount(c *gin.Context) {
	id,err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid account id",
		})
		return
	}

	var updatedAccount models.Account

	if err := c.ShouldBindJSON(&updatedAccount); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})

		return 
	}

	for index, account := range Accounts {
		if account.ID == id {
			if updatedAccount.AccountName != "" {
				Accounts[index].AccountName = updatedAccount.AccountName
			}


			if updatedAccount.Status != "" {
				Accounts[index].Status = updatedAccount.Status 
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Account updated successfully",
				"account": Accounts[index],
			})
			return 
		}
		
	}

	c.JSON(http.StatusNotFound, gin.H{
		"error": "Account not found",
	})
}

func DeleteAccount(c *gin.Context){
	id,err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid account id",
		})
		return
	}

	for index, account := range Accounts {
		if account.ID == id {
			Accounts = append(Accounts[:index], Accounts[index+1:]...)
			c.JSON(http.StatusOK, gin.H{
				"message": "Account deleted successfully",
			})
			return 
		}
	}
	c.JSON(http.StatusNotFound, gin.H{
		"error": "Account not found",
	})
}

func GetAccountBalance(c *gin.Context){
	id,err  := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid account id",
		})
		return
	}

	for _, account := range Accounts {
		if account.ID == id {
			c.JSON(http.StatusOK, gin.H{
				"balance": account.Balance,
			})
			return 
		}
	}
	c.JSON(http.StatusNotFound, gin.H{
		"error": "Account not found",
	})
}

func GetAccountTransfers(c *gin.Context){
	id,err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid account id",
		})
		return
	}

	accountTransfers := []models.Transfer {}

	for _, transfer := range Transfers {
		if transfer.FromAccountID == id || transfer.ToAccountID == id {
			accountTransfers = append(accountTransfers, transfer)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"transfers": accountTransfers,
	})
}