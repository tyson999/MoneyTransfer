package controllers

import (
	"MoneyTransfer/database"
	"MoneyTransfer/models"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5" 
)

// In-memory users store (used by some handlers)
var Users []models.User

func CreateUser(c *gin.Context) {

	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	query := `
	INSERT INTO users 
	(name, email, phone, password_hash) 
	VALUES ($1, $2, $3, $4) 
	RETURNING id,created_at,updated_at;
	`
	err := database.DB.QueryRow(
		c.Request.Context(),
		query,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.Phone,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt,)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    user,
	})
}

func GetUsers(c *gin.Context) {

	query := `
	SELECT id,name,email,phone,created_at,updated_at
	FROM users
	ORDER BY id;
	`
	rows, err := database.DB.Query(
		c.Request.Context(),
		query,
	)

	if err!= nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"error":err.Error(),
		})
		return 
	}

	defer rows.Close()

	users:= make([]models.User,0)

	for rows.Next(){
		var user models.User 

		err:= rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Phone,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err!= nil{
			c.JSON(http.StatusInternalServerError,gin.H{
				"error":err.Error(),
			})
			return
		}

		users = append(users, user)
	}

	if err:= rows.Err(); err!= nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"error": err.Error(),
		})
		return 
	}

	c.JSON(http.StatusOK,gin.H{
		"users": users,
	})
}

func GetUserByID(c *gin.Context) {

	id := c.Param("id")

	query:= `
	SELECT id,name
	FROM users
	WHERE id = $1;
	`

	var user models.User

	err := database.DB.QueryRow(
		c.Request.Context(),
		query,
		id,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Phone,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err!= nil{
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound,gin.H{
				"error":"User not found",

			})

			return
		}

		c.JSON(http.StatusInternalServerError,gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func UpdateUser(c *gin.Context) {

	id := c.Param("id")

	var input models.User

	if err := c.ShouldBindJSON(&input); err!= nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	query := `
	UPDATE users
	SET name = $1,
	email = $2,
	phone = $3,
	updated_at = NOW()
	WHERE id = $4
	RETURNING id,name,email,phone,created_at,updated_at;
	`

	var user models.User

	err:= database.DB.QueryRow(
		c.Request.Context(),
		query,
		input.Name,
		input.Email,
		input.Phone,
		id,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Phone,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err==pgx.ErrNoRows {
			c.JSON(http.StatusNotFound,gin.H{
				"error":"User not found",
			})
			return 
		}

		c.JSON(http.StatusInternalServerError,gin.H{
			"error": err.Error(),
		})

		return 
	}

		c.JSON(http.StatusOK, gin.H{
			"message": "User updated successfully",
			"user": user,
		})
}


func DeleteUser(c *gin.Context) {

	id := c.Param("id")

	query:= `
	DELETE FROM users
	WHERE id = $1;
	`
	result,err:= database.DB.Exec(
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

	if result.RowsAffected()==0{
		c.JSON(http.StatusNotFound, gin.H{
			"error":"User not found",
		})
		return 
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})

}
