package models

type Bidder struct {
	BidderId  string `json:"bidder_id"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}
