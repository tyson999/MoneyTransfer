package controllers

import (
	"net/http"
	"github.com/jackc/pgx/v5" 
	"MoneyTransfer/models"
	"MoneyTransfer/database" 
	"github.com/gin-gonic/gin"
)

var companies = []models.Company{}


func CreateCompany(c *gin.Context) {

	var company models.Company

	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Check if user exists
	var creatorID int

	err := database.DB.QueryRow(
		c.Request.Context(),
		`SELECT id FROM users WHERE id=$1`,
		company.CreatedBy,
	).Scan(&creatorID)

	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Creator user not found",
			})
			return 
		}

		c.JSON(http.StatusInternalServerError,gin.H{
			"error":err.Error(),
		})
		return 
	}

	query:= `
	INSERT INTO companies 
	(created_by,company_name,gstin)
	VALUES($1,$2,$3)
	RETURNING id,created_at,updated_at;
	`

	err = database.DB.QueryRow(
		c.Request.Context(),
		query,
		company.CreatedBy,
		company.CompanyName,
		company.GSTIN,
	).Scan(
		&company.ID,
		&company.CreatedAt,
		&company.UpdatedAt,
	)

	if err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"error":err.Error(),
		})
		return 
	}


	c.JSON(http.StatusCreated, gin.H{
		"message":"Company created successfully",
		"company":company,
	})
}

func GetCompanies(c *gin.Context) {

	query:= `
	SELECT id,created_by,company_name,gstin,created_at,updated_at
	FROM companies
	ORDER BY id;
	`

	rows,err:= database.DB.Query(
		c.Request.Context(),
		query,
	)

	if err!= nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		
		})
		return 
	}

	defer rows.Close()

	companies := []models.Company{}

	for rows.Next() {
		var company models.Company

		err:= rows.Scan(
			&company.ID,
			&company.CreatedBy,
			&company.CompanyName,
			&company.GSTIN,
			&company.CreatedAt,
			&company.UpdatedAt,
		)


		if err!= nil{
			c.JSON(http.StatusInternalServerError,gin.H{
				"error":err.Error(),

			})
			return 
		}

		companies = append(companies,company)

		c.JSON(http.StatusOK,gin.H{
			"companies":companies,
		})
	}
}

func GetCompanyByID(c *gin.Context) {

	id:= c.Param("id")

	query:= `
	SELECT
	id,
	company_name
	FROM companies
	WHERE id=$1;
	`

	var company models.Company

	err := database.DB.QueryRow(
		c.Request.Context(),
		query,
		id,
	).Scan(
		&company.ID,
		&company.CompanyName,
	)

	if err != nil {

		if err == pgx.ErrNoRows {

			c.JSON(http.StatusNotFound, gin.H{
				"error": "Company not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError,gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"company":company,
	})
}

func UpdateCompany(c *gin.Context) {

	id:= c.Param("id")

	var company models.Company

	if err:= c.ShouldBindJSON(&company); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"error": "Invalid request body",
		})
		return 
	}

	query:=`
	UPDATE companies
	SET
	company_name=$1,
	gstin=$2,
	updated_at=NOW()
	WHERE id= $3
	RETURNING id,created_by,company_name,gstin,created_at,updated_at
	`

	var updatedCompany models.Company 

	err:= database.DB.QueryRow(
		c.Request.Context(),
		query,
		company.CompanyName,
		id,
	).Scan(
		&updatedCompany.ID,
		&updatedCompany.CreatedBy,
		&updatedCompany.CompanyName,
		&updatedCompany.GSTIN,
		&updatedCompany.CreatedAt,
		&updatedCompany.UpdatedAt,
	)

	if err!= nil{
		if err==pgx.ErrNoRows {
			c.JSON(http.StatusInternalServerError,gin.H{
				"error":err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
		"message": "Company updated successfully",
		"company": updatedCompany,
	    })
    }
}

func DeleteCompany(c *gin.Context) {

	id := c.Param("id")

	query := `
	DELETE FROM companies
	WHERE id=$1;
	`

	result, err := database.DB.Exec(
		c.Request.Context(),
		query,
		id,
	)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if result.RowsAffected() == 0 {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "Company not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Company deleted successfully",
	})
}

