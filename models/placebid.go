package models

type PlaceBidRequest struct {
	BidderId  string  `json:"bidder_id" binding:"required"`
	AuctionId string  `json:"auction_id" binding:"required"`
	BidAmount float64 `json:"bid_amount" binding:"required"`
}
