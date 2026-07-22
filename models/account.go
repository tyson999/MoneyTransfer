package models

import "time"

//ID,Createdat and updated at and status are given by server and not by user. So, they are not included in the request body.

type Account struct {
	ID 	 int    `json:"id"`
	CompanyID int    `json:"company_id"`
	AccountName string `json:"account_name"`
	AccountType string `json:"account_type"`        //cash or credit
	Balance float64 `json:"balance"`
	Status string `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}