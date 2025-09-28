package service

import (
	"log"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

func ScheduleAuction(db *gorm.DB) {
	c := cron.New(cron.WithSeconds())

	_, err := c.AddFunc("*/10 * * * * *", func() {

		//to fetch scheduled auctions and convert them to live
		nextAuctions, _ := FetchNextLiveAuctions(db) 

		//to fetch live auctions and convert them to closed
		closedAuctions, _ := CloseLiveAuction(db)

		if len(closedAuctions) != 0 {
			log.Println(len(closedAuctions), "auctions got closed..!!")
		}
		//sending notification to all the users for all the auction events
		if len(nextAuctions) != 0 {
			log.Println(len(nextAuctions), " auctions are scheduled..!!")

			bidders, _ := GetAllBidders(db)
			for _, bidder := range bidders {
				SendNotification(bidder.BidderId, nextAuctions)
			}
			log.Println("successfully send notification to", len(bidders), "bidders..!!")
		}

	})
	if err != nil {
		log.Fatalf("Failed to start auction cron job: %v", err)
	}
	log.Println("scheduler running successfully!")
	c.Start()
}
