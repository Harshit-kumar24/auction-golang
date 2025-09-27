package service

import (
	"fmt"

	"github.com/Harshit-kumar24/eauction/models"
)

func SendNotification(userId string, auctions []models.Auction) {
	fmt.Println("successfully send email to userId: ", userId, " for ", len(auctions), " aucitons...")
}
