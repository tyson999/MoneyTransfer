package controllers 

import(
	"net/http"
	"MoneyTransfer/models"
	"github.com/gin-gonic/gin"
	"strconv"
)


var Transfers = []models.Transfer{}

func CreateTransfer(c *gin.Context){
	var transfer models.Transfer

	if err := c.ShouldBindJSON(&transfer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return 
	}

	if transfer.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Amount must be greater than 0",
		})
		return
	}

	if transfer.FromAccountID == transfer.ToAccountID {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Sender and receiver cannot be the same account",
		})
		return
	}

	senderIndex := -1
	receiverIndex := -1

	for index, account := range Accounts {
		if account.ID == transfer.FromAccountID {
			senderIndex = index
		}

		if account.ID == transfer.ToAccountID {
			receiverIndex = index
		}
	}

	if senderIndex == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Sender account not found",
		})
		return
	}

	if receiverIndex == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Receiver account not found",
		})
		return
	}

	if Accounts[senderIndex].Status != "active" ||
		Accounts[receiverIndex].Status != "active" {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Both accounts must be active",
		})
		return
	}

	if Accounts[senderIndex].Balance < transfer.Amount {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Insufficient balance",
		})
		return
	}

	Accounts[senderIndex].Balance -= transfer.Amount
	Accounts[receiverIndex].Balance += transfer.Amount

	transfer.ID = len(Transfers) + 1
	transfer.Status = "success"

	Transfers = append(Transfers, transfer)

	c.JSON(http.StatusCreated, gin.H{
		"message":          "Transfer successful",
		"transfer":         transfer,
		"sender_balance":   Accounts[senderIndex].Balance,
		"receiver_balance": Accounts[receiverIndex].Balance,
	})
}
	

func GetTransfer(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"transfers": Transfers,
	})
}

func GetTransferByID(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid transfer id",
		})
		return
	}

	for _, transfer := range Transfers {
		if transfer.ID == id {
			c.JSON(http.StatusOK, gin.H{
				"transfer": transfer,
			})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{
		"error": "Transfer not found",
	})
}

func UpdateTransfer(c *gin.Context) {
	idParam := c.Params.ByName("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid transfer id",
		})
		return
	}

	var updatedTransfer models.Transfer

	if err := c.ShouldBindJSON(&updatedTransfer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	for i, transfer := range Transfers {
		if transfer.ID == id {
			Transfers[i] = updatedTransfer
			c.JSON(http.StatusOK, gin.H{
				"transfer": updatedTransfer,
			})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{
		"error": "Transfer not found",
	})
}

func DeleteTransfer(c *gin.Context) {
	idParam := c.Params.ByName("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid transfer id",
		})
		return
	}

	for i, transfer := range Transfers {
		if transfer.ID == id {
			Transfers = append(Transfers[:i], Transfers[i+1:]...)
			c.JSON(http.StatusOK, gin.H{
				"message": "Transfer deleted successfully",
			})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{
		"error": "Transfer not found",
	})
}