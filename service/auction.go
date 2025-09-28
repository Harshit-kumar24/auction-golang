package service

import (
	"log"
	"net/http"
	"time"

	// "time"

	"github.com/Harshit-kumar24/eauction/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// handlers
func SaveAuction(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var auction models.Auction

		//binding inputs to json
		if err := c.ShouldBindBodyWithJSON(&auction); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "the request inputs are wrong."})
			return
		}

		if auction.ItemID == "" {
			auction.ItemID = uuid.New().String()
		}

		//saving to database
		if err := db.Create(&auction).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "an internal error occurred..!!"})
		} else {
			c.JSON(http.StatusCreated, gin.H{
				"message": "Auciton created sucessfully.",
			})
		}
	}
}

func GetAllLiveAuctions(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var auctions []models.Auction

		result := db.Where("auction_status = ?", "live").Find(&auctions)
		if result.Error != nil {
			log.Println("error fetching live auctions...")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching live auctions"})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"data": auctions,
			})
		}
	}
}

// functions
func FindAuctionById(auctionId string, db *gorm.DB) (*models.Auction, error) {
	var auction models.Auction

	result := db.First(&auction, "item_id = ?", auctionId)
	if result.Error != nil {
		log.Println("error fetching auction: ", result.Error)
		return nil, result.Error
	}
	return &auction, nil
}

func FetchNextLiveAuctions(db *gorm.DB) ([]models.Auction, error) {
	var auctions []models.Auction

	currentTime := time.Now()

	result := db.Where("auction_status = ? AND auction_start_time <= ? AND auction_end_time >= ?",
		"scheduled", currentTime, currentTime).Find(&auctions)
	if result.Error != nil {
		log.Println("error fetching auctions...")
		return nil, result.Error
	}
	for _, auction := range auctions {
		if err := db.Model(&models.Auction{}).
			Where("item_id = ?", auction.ItemID).
			Update("auction_status", "live").Error; err != nil {
			continue
		}
	}
	return auctions, nil
}

func GetClosedAuctions(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var auctions []models.Auction

		result := db.Where("auction_status = ?", "closed").Find(&auctions)
		if result.Error != nil {
			log.Println("error fetching closed auctions...")
			c.JSON(http.StatusBadRequest, gin.H{"message": "error getting closed auctions"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "sucessfully fetched closed auctions",
			"data":    auctions})
	}
}

func CloseLiveAuction(db *gorm.DB) ([]models.Auction, error) {

	var currentLiveAuctions []models.Auction
	currentTime := time.Now()

	//fetching all the auctions that are still not closed
	result := db.Where("auction_status = ? AND auction_end_time <= ?", "live", currentTime).Find(&currentLiveAuctions)
	if result.Error != nil {
		log.Println("error fetching auctions...")
		return nil, result.Error
	}

	var closedAuctions []models.Auction
	for _, auction := range currentLiveAuctions {
		// if auction.AuctionEndTime.Before(currentTime) {
		if err := db.Model(&models.Auction{}).
			Where("item_id = ?", auction.ItemID).
			Update("auction_status", "closed").Error; err != nil {
			log.Printf("Failed to update auction %s: %v", auction.ItemID, err)
			continue
		}
		auction.AuctionStatus = "closed"
		closedAuctions = append(closedAuctions, auction)
		log.Printf("Auction %s is now closed", auction.ItemID)
		log.Printf("Winner of the auction %d is %s with bid price: %.2f", auction.ItemID, auction.CurrentWinner, auction.HighestBid)
	}
	return closedAuctions, nil
}

func TotalTimeofAllAuctions(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var durationSeconds int64

		row := db.Raw(`
        SELECT FLOOR(EXTRACT(EPOCH FROM MAX(auction_end_time) - MIN(auction_start_time)))::bigint AS duration_seconds 
        FROM public.auctions
    `).Row()

		if err := row.Scan(&durationSeconds); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"totalTimeHours": 0,
				"message":        "error getting total time"})
		}

		log.Println("Total duration (seconds):", durationSeconds)
		c.JSON(http.StatusOK, gin.H{"totalTimeHours": durationSeconds / 3600})
	}
}
