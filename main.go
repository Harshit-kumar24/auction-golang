package main

import (
	"log"

	"github.com/Harshit-kumar24/eauction/config"
	"github.com/Harshit-kumar24/eauction/db"
	"github.com/Harshit-kumar24/eauction/service"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	cfg := config.LoadConfig()
	db := db.InitDB(cfg)

	//limiting vCPU and memory
	config.SetupResources()

	r := gin.Default()

	//run the auction scheduler
	go service.ScheduleAuction(db)

	apiV1 := r.Group("/api/v1")
	{
		apiV1.POST("/bidder", service.SaveBidder(db))
		apiV1.POST("/auction", service.SaveAuction(db))
		apiV1.GET("/auction", service.GetAllLiveAuctions(db))
		apiV1.POST("/bid", service.PlaceBid(db))
		apiV1.GET("/auctionTime", service.TotalTimeofAllAuctions(db))
		apiV1.GET("/closedAuctions", service.GetClosedAuctions(db))
	}

	log.Printf("Server starting on port %s...", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

}
