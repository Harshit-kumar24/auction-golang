package models

import "time"

type Auction struct {
	ItemID          string `json:"item_id"`
	ItemName        string `json:"item_name" binding:"required"`
	ItemCategory    string `json:"item_category" binding:"required"`
	ItemDescription string `json:"item_desc" binding:"required"`
	ItemCondition   string `json:"item_condition" binding:"required"`

	StartingBid    float64 `json:"starting_bid" binding:"required"`
	EstimatedValue float64 `json:"estimated_value" binding:"required"`
	ReservedPrice  float64 `json:"reserved_price" binding:"required"`
	HighestBid     float64 `json:"highest_bid"`
	BidIncrement   float64 `json:"bid_increment" binding:"required"`
	Currency       string  `json:"currency" binding:"required"`

	AuctionStartTime time.Time `json:"auction_start_time" binding:"required"`
	AuctionEndTime   time.Time `json:"auction_end_time" binding:"required"`
	AuctionDuration  int       `json:"auction_duration" binding:"required"`
	TimeZone         string    `json:"timesone" binding:"required"`
	CurrentWinner    string    `json:"current_winner"`
	AuctionStatus    string    `json:"auction_status" binding:"required"`

	SellerID          string  `json:"sellerId" binding:"required"`
	PopularityScore   string  `json:"popularity_score" `
	ItemRarity        string  `json:"item_rarity" binding:"required"`
	ShippingCost      float64 `json:"shipping_cost" binding:"required"`
	SpecialAttributes string  `json:"special_attributes"`
}
