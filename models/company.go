package models 
import "time"


//ID,Createdat and updated at are given by server and not by user. So, they are not included in the request body.

type Company struct {
	ID int `json:"id"`
	CreatedBy int `json:"user_id"`
	CompanyName string `json:"company_name"`
	GSTIN string `json:"gstin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}