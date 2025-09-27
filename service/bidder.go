package service

import (
	"log"
	"net/http"

	"github.com/Harshit-kumar24/eauction/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SaveBidder(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var bidder models.Bidder

		//binding request input with json
		if err := c.ShouldBindBodyWithJSON(&bidder); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "the request inputs are wrong."})
			return
		}

		//generating uuid
		if bidder.BidderId == "" {
			bidder.BidderId = uuid.New().String()
		}

		//saving to database
		if err := db.Create(&bidder).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "an internal error occurred..!!"})
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Bidder saved successfully.",
			"bidder":  bidder.BidderId,
		})
	}
}

func GetAllBidders(db *gorm.DB) ([]models.Bidder, error) {
	var bidders []models.Bidder

	result := db.Find(&bidders)
	if result.Error != nil {
		log.Println("Error fetching bidders:", result.Error)
		return nil, result.Error
	}
	return bidders, nil
}

func PlaceBid(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var placeBidRequest models.PlaceBidRequest

		if err := c.ShouldBindJSON(&placeBidRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		requestBidAmount := placeBidRequest.BidAmount
		requestBidderId := placeBidRequest.BidderId

		auction, err := FindAuctionById(placeBidRequest.AuctionId, db)

		if err != nil || auction.ItemID == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "auction with id not found"})
		}

		log.Println("current auction id:", auction.ItemID)

		//start transaction
		tx := db.Begin()
		if tx.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to start transaction"})
			return
		}
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		//acquire lock for the specific row
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("item_id = ?", auction.ItemID).
			First(&auction).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": "auction not found"})
			return
		}

		if auction.HighestBid < requestBidAmount {

			if err := tx.Model(&auction).Where("item_id = ?", auction.ItemID).Updates(map[string]interface{}{
				"highest_bid":    requestBidAmount,
				"current_winner": requestBidderId,
			}).Error; err != nil {
				tx.Rollback()
				log.Println("error updating highest bid:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error saving bid"})
				return
			}

		} else {
			tx.Rollback()
			c.JSON(http.StatusConflict, gin.H{"error": "bid is lower than current highest"})
			return
		}

		//commit the transaction
		if err := tx.Commit().Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":      "bid placed successfully",
			"auction":      auction.ItemName,
			"currentPrice": auction.HighestBid,
		})
	}
}
