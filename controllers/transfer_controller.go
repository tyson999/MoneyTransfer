package controllers 

import(
	"net/http"
	"MoneyTransfer/models"
	"MoneyTransfer/database"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"errors"
)


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

	ctx:= c.Request.Context() 

	tx,err:= database.DB.Begin(ctx)
	if err!= nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"error": "Failed",
		})
		return 
	}

	defer tx.Rollback(ctx)

	var senderBalance float64 

	var senderStatus string 

	err = tx.QueryRow(
		ctx,
		`
		SELECT balance,status
		FROM accounts
		WHERE id =$1
		FOR UPDATE;
		`,
		transfer.FromAccountID,
	).Scan(
		&senderBalance,
		&senderStatus,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Sender account not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to read sender account",
		})
		return
	}

	if senderStatus != "active" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Sender account is not active",
		})
		return
	}

	var receiverStatus string

	err = tx.QueryRow(
		ctx,
		`
		SELECT status
		FROM accounts
		WHERE id = $1
		FOR UPDATE;
		`,
		transfer.ToAccountID,
	).Scan(&receiverStatus)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Receiver account not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to read receiver account",
		})
		return
	}

	if receiverStatus != "active" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Receiver account is not active",
		})
		return
	}

	if senderBalance < transfer.Amount {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Insufficient balance",
		})
		return
	}

	
	_, err = tx.Exec(
		ctx,
		`
		UPDATE accounts
		SET
			balance = balance - $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $2;
		`,
		transfer.Amount,
		transfer.FromAccountID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to debit sender account",
		})
		return
	}

	
	_, err = tx.Exec(
		ctx,
		`
		UPDATE accounts
		SET
			balance = balance + $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $2;
		`,
		transfer.Amount,
		transfer.ToAccountID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to credit receiver account",
		})
		return
	}

	transfer.Status = "completed"

	
	err = tx.QueryRow(
		ctx,
		`
		INSERT INTO transfers (
			from_account_id,
			to_account_id,
			amount,
			status,
			description
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at;
		`,
		transfer.FromAccountID,
		transfer.ToAccountID,
		transfer.Amount,
		transfer.Status,
		transfer.Description,
	).Scan(
		&transfer.ID,
		&transfer.CreatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save transfer",
		})
		return
	}

	
	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to complete transfer",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Transfer completed successfully",
		"transfer": transfer,
	})
}


	

func GetTransfer(c *gin.Context) {
	query := `
		SELECT
			id,
			from_account_id,
			to_account_id,
			amount,
			status,
			description,
			created_at
		FROM transfers
		ORDER BY id DESC;
	`

	rows, err := database.DB.Query(
		c.Request.Context(),
		query,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch transfers",
		})
		return
	}
	defer rows.Close()

	transfers := []models.Transfer{}

	for rows.Next() {
		var transfer models.Transfer

		err := rows.Scan(
			&transfer.ID,
			&transfer.FromAccountID,
			&transfer.ToAccountID,
			&transfer.Amount,
			&transfer.Status,
			&transfer.Description,
			&transfer.CreatedAt,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to read transfer data",
			})
			return
		}

		transfers = append(transfers, transfer)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed while reading transfers",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transfers": transfers,
	})
}

func GetTransferByID(c *gin.Context) {
	id := c.Param("id")

	var transfer models.Transfer

	query := `
		SELECT
			id,
			from_account_id,
			to_account_id,
			amount,
			status,
			description,
			created_at
		FROM transfers
		WHERE id = $1;
	`

	err := database.DB.QueryRow(
		c.Request.Context(),
		query,
		id,
	).Scan(
		&transfer.ID,
		&transfer.FromAccountID,
		&transfer.ToAccountID,
		&transfer.Amount,
		&transfer.Status,
		&transfer.Description,
		&transfer.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Transfer not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch transfer",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transfer": transfer,
	})
}
func UpdateTransfer(c *gin.Context) {
	id := c.Param("id")

	var transfer models.Transfer

	if err := c.ShouldBindJSON(&transfer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	query := `
	UPDATE transfers
	SET
		from_account_id = $1,
		to_account_id = $2,
		amount = $3,
		status = $4,
		description = $5
	WHERE id = $6
	RETURNING
		id,
		from_account_id,
		to_account_id,
		amount,
		status,
		description,
		created_at;
	`

	err := database.DB.QueryRow(
		c.Request.Context(),
		query,
		transfer.FromAccountID,
		transfer.ToAccountID,
		transfer.Amount,
		transfer.Status,
		transfer.Description,
		id,
	).Scan(
		&transfer.ID,
		&transfer.FromAccountID,
		&transfer.ToAccountID,
		&transfer.Amount,
		&transfer.Status,
		&transfer.Description,
		&transfer.CreatedAt,
	)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Transfer not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update transfer",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transfer updated successfully",
		"transfer": transfer,
	})
}

func DeleteTransfer(c *gin.Context) {

	id := c.Param("id")

	query := `
	DELETE FROM transfers
	WHERE id = $1;
	`

	commandTag, err := database.DB.Exec(
		c.Request.Context(),
		query,
		id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete transfer",
		})
		return
	}

	if commandTag.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Transfer not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transfer deleted successfully",
	})
}
	
	