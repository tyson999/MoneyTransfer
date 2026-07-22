//controller handles the http specific work(request and response)
//number of routes created should be equal to the number of functions created in the controller file
//the relation between routes and controllers is that the routes file defines the endpoints for the APIs, while the controller file contains the logic for handling requests to those endpoints. Each route corresponds to a specific function in the controller that processes the request and returns a response.
//service files are created to handle the business logic and to make our controller file more organized and maintainable. The service files can contain functions that handle the core business logic, while the controller file can focus on handling HTTP requests and responses. This separation of concerns can make the code easier to read, test, and maintain.

package controllers

import(
	"net/http"
	"MoneyTransfer/models"
	"github.com/gin-gonic/gin"
	"MoneyTransfer/database"
	"errors" 
	"github.com/jackc/pgx/v5"
)


func CreateAccount(c *gin.Context){
	var account models.Account

	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return 

	}

	if account.CompanyID == 0 ||
	    account.AccountName==""||
		account.AccountType=="" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "company_id,account_name and account_type are required",
				
			})

			return 
		}

	if account.Status == "" {
		account.Status = "active"
	}

	var companyID int 

	err := database.DB.QueryRow(
		c.Request.Context(),
		`
		SELECT id
		FROM companies
		WHERE id = $1;
		`,
		account.CompanyID,
	).Scan(&companyID)

	if err!= nil {
		if errors.Is(err,pgx.ErrNoRows){
			c.JSON(http.StatusNotFound,gin.H{
				"error":"Company not found",
			})
			return 
		}

		c.JSON(http.StatusInternalServerError,gin.H{
			"error":"Failed to check company",
		})
		return 
	}


	query:=`
	INSERT INTO accounts
	(company_id,
	account_name,
	account_type,
	balance,
	status
	)
	VALUES ($1,$2,$3,$4,$5)
	RETURNING id,created_at,updated_at;
	`
	err = database.DB.QueryRow(
		c.Request.Context(),
		query,
		account.CompanyID,
		account.AccountName,
		account.AccountType,
		account.Balance,
		account.Status,
	).Scan(
		&account.ID,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err!= nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"error": "Failed to create account",
		})
		return 
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":"Account created",
		"account": account,
	})
}	


func GetAccounts(c *gin.Context) {
	query := `
		SELECT
			id,
			company_id,
			account_name,
			account_type,
			balance,
			status,
			created_at,
			updated_at
		FROM accounts
		ORDER BY id;
	`

	rows, err := database.DB.Query(
		c.Request.Context(),
		query,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch accounts",
		})
		return
	}
	defer rows.Close()

	accounts := []models.Account{}

	for rows.Next() {
		var account models.Account

		err := rows.Scan(
			&account.ID,
			&account.CompanyID,
			&account.AccountName,
			&account.AccountType,
			&account.Balance,
			&account.Status,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to read account data",
			})
			return
		}

		accounts = append(accounts, account)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed while reading accounts",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accounts": accounts,
	})
}

func GetAccountByID(c *gin.Context) {
	id:=c.Param("id")

	var account models.Account

	query:= `
	SELECT
	id,
	company_id,
	account_name,
	account_type,
	balance,
	status,
	created_at,
	updated_at
	FROM accounts
	WHERE id = $1;
	`

	err:= database.DB.QueryRow(
		c.Request.Context(),
		query,
		id,
	).Scan(
		&account.ID,
		&account.CompanyID,
		&account.AccountName,
		&account.AccountType,
		&account.Balance,
		&account.Status,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Account not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch account",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"account": account,
	})
}


func UpdateAccount(c *gin.Context) {
	id := c.Param("id")

	var account models.Account

	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if account.CompanyID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "company_id is required",
		})
		return
	}

	if account.AccountName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "account_name is required",
		})
		return
	}

	if account.AccountType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "account_type is required",
		})
		return
	}

	if account.Balance < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "balance cannot be negative",
		})
		return
	}

	if account.Status == "" {
		account.Status = "active"
	}

	// Verify that the company exists.
	var companyID int

	err := database.DB.QueryRow(
		c.Request.Context(),
		`
		SELECT id
		FROM companies
		WHERE id = $1;
		`,
		account.CompanyID,
	).Scan(&companyID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Company not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to verify company",
		})
		return
	}

	query := `
		UPDATE accounts
		SET
			company_id = $1,
			account_name = $2,
			account_type = $3,
			balance = $4,
			status = $5,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $6
		RETURNING
			id,
			company_id,
			account_name,
			account_type,
			balance,
			status,
			created_at,
			updated_at;
	`

	err = database.DB.QueryRow(
		c.Request.Context(),
		query,
		account.CompanyID,
		account.AccountName,
		account.AccountType,
		account.Balance,
		account.Status,
		id,
	).Scan(
		&account.ID,
		&account.CompanyID,
		&account.AccountName,
		&account.AccountType,
		&account.Balance,
		&account.Status,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Account not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update account",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Account updated successfully",
		"account": account,
	})
}
	

	

func DeleteAccount(c *gin.Context){
	id:= c.Param("id")

	query := `
		DELETE FROM accounts
		WHERE id = $1;
	`

	commandTag, err := database.DB.Exec(
		c.Request.Context(),
		query,
		id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete account",
		})
		return
	}

	if commandTag.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Account not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Account deleted successfully",
	})
}