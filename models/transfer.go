package models

import "time"

//ID,Createdat and updated at and status are given by server and not by user. So, they are not included in the request body.


type Transfer struct{
	ID int `json:"id"`
	FromAccountID int `json:"from_account_id"`
	ToAccountID int `json:"to_account_id"`
	Amount float64 `json:"amount"`
	Status string `json:"status"`
	Description string `json:"description"`
	CreatedAt time.Time `json:"created_at"`
}

