### Objective 
The goal of this exercise is to design and implement an Auction Simulator that runs
multiple auctions concurrently, collects bids from simulated bidders, and measures
execution times while adhering to resource constraints.
1. There are **100 bidders** participating.
2. Each auction is generated based on 20 attributes of an object.
3. All bidders receive these **20 attributes** and can respond with a bid. It is not
necessary that every bidder will provide a response.
4. Each auction will run with a timeout. Once the timeout is reached, the auction
 will close, and the winner will be declared from the bids received up to that point.
5. A total of **40 auctions** will run concurrently (at the same time).
6. Measure the time taken between the start of the first auction and the completion
of the last auction.
7. Provide a mechanism to standardize resources with respect to the vCPU and
RAM available.

### How to setup the project
**github repo:** https://github.com/Harshit-kumar24/auction-golang.git
- simply clone the above repo in your VM or your system 
- install go compiler to run the project if your already have then follow the next step
- install **postgres** and execute the migration script that is inside **migrations/migration-1.sql**
- Now, go inside the project with the name **auction-golang** 
- you will find the file **main.go** if yes then run the following command 
```go
go run main.go
```
- you will see output like this,
```
successfully connected to the postgres database...!!
2025/09/28 05:19:15 Setting CPU cores to 4
2025/09/28 05:19:15 Setting memory limit for buffers to 2048 MB
2025/09/28 05:19:15 Server starting on port 8080...
2025/09/28 05:19:15 scheduler running successfully!
```
- congratulations! the proejct is now live 
---
### Sample Architecture 
![[Pasted image 20250928063015.png]]
### How each functionality is working

1. **All bidders receive these 20 attributes and can respond with a bid. It is not
necessary that every bidder will provide a response**
- so usually when a **auction goes live** and any user has to be notified **we usually notify on the application portal or through email** here for the project I have simulated using a **email** 
- so when ever a auction goes live all the users **100 users** will recieve a email and get log like this
```
successfully send email to userId:  931a91bf-ca77-4d4c-a5e0  for  1  auctions...
successfully send email to userId:  ef8d1ca3-322c-4810-971d  for  1  auctions...
successfully send email to userId:  b654075e-aacd-4f1f-9a41  for  1  auctions...
successfully send email to userId:  788391ca-833d-4ad3-a386  for  1  auctions...
successfully send email to userId:  6cf31ee2-490d-463f-9625  for  1  auctions...
successfully send email to userId:  0db0b37c-7ba6-4f7a-8f30  for  1  auctions...
```
- this simulate sending notfication so that the users can get notified about live auctions
- **or we can do the same thing when we create a auction as well**
---
2. **Each auction will run with a timeout. Once the timeout is reached, the auction
 will close, and the winner will be declared from the bids received up to that point**

- to get the start and end time of the auction event we have two columns 
```psql
auction_start_time //denotes start of auction time
auction_end_time //denotes end of auction time
```
- so between this time the auction will be live
- A auction can be in 3 states **(scheduled, live, closed)** this means,
	- **scheduled**: means the auction is about to start 
	- **live**: means the auction is on going 
	- **closed**: means the auction ended
- so to bid for a auction we have,
```
POST http://localhost:8080/api/v1/bid
{
    "bidder_id":"b654075e-aacd-4f1f-9a41-cb9ded52189b",
    "auction_id":"ea777447-5e06-49b1-8608-84564e9c662d",
    "bid_amount": 25.34
}
```
- so when you hit this endpoint 3 things can happen either the auction is not present or is already closed 
- the bid amount by the user is less than the current bid amount 
- the bid got successful and you are the highest bidder till now 
#### How this functionality works
1. There will be a **go routine** that will run a scheduler on every **10 seconds** this value can be changed based on **server and db** configurations
2. the scheduler will look in the table **public.auctions** and find the auctions with status **scheduled** and **start time >= current time** if its true, then the **auction is live** and status will be changed to **live**
3. bidders will recieve notification and can bid on this auction 
4. what will happen is when the bidder places a bid it will check the **bid amount** with the **higest_bid** column if its greater than that then the current user will be the higest bidder and then **current_winner** will store the **primary key of the user**
5. if the bid is less then simply it will be rejected.
6. Once the auction ends the whatever will be inside **current_winner** that user will be declared as **WIINNER** and get a log like this,
```
 Winner of the auction %!d(string=10534b8e-97b0-42fb-8024-ef2a218544fa) is f47ac10b-58cc-4372-a567-0e02b2c3d479 with bid price: 30.36
```
7. also you can check the winner of every auction till today using this endpoint 
```
GET localhost:8080/api/v1/closedAuctions
```

### Problem that should be addressed 
- there can be a problem when two bidders at exact same time bid for the same auction with the same price then which user should be updated?
- what you can do for this we can add a **row level transaction at auction** so once this tranction is updated by one user then the other can since once updated then other will be rejected and will be asked for more higer bid 
- this ensures **consistency** across at all place 
- this transaction can be added by writing some code like this,
```go
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
            First(&auction).Error; err != nil {
            tx.Rollback()
            c.JSON(http.StatusNotFound, gin.H{"error": "auction not found"})
            return
        }
	
	if auction.HighestBid < requestBidAmount {
	//actual log goes here 
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
```
- or search for function **PlaceBid( )** in file **service/bidder.go**
---
3. **Measure the time taken between the start of the first auction and the completion
of the last auction**
- to measure this time between the starting time of first auction and ending time of last auction till date we can get result from query like this,
```
SELECT 
FLOOR(EXTRACT(EPOCH FROM MAX(auction_end_time) - MIN(auction_start_time)))::bigint AS duration_seconds
FROM public.auctions
```
- this will give result in second you can convert this into **days or hours** accordingly 
- you can check the result by hitting this endpoint 
```
GET http://localhost:8080/api/v1/auctionTime
```
---
4. **Provide a mechanism to standardize resources with respect to the vCPU and
RAM available**
- to set the **number of CPUs** you can set this value from the go runtime from the code liike this 
```
runtime.GOMAXPROCS(numCPU) //numCPU refers to number of cpu
```
- since something like a channel is not used in the project but you can easily **limit resouces of a channel** by adding a code snippet like this,
```
bufferSize := 1024 // default
if maxMemoryMB > 0 {
	bufferSize = (maxMemoryMB * 1024 * 1024) 
}
ch := make(chan []byte, bufferSize)
```
- **or if you want to do for whole project you can set this env at runtime like this,**
```
GOMEMLIMIT=500MiB ./auctiongolang
```
---
5. **Each auction is generated based on 20 attributes of an object.**
- the auction can be created by making a curl like this 
```
curl --location 'http://localhost:8080/api/v1/auction' \
--header 'Content-Type: application/json' \
--data '{
  "item_name": "Vintage Watch",
  "item_category": "Collectibles",
  "item_desc": "A rare vintage wristwatch from 1960s in mint condition.",
  "item_condition": "Mint",
  "starting_bid": 1500.00,
  "estimated_value": 2500.00,
  "reserved_price": 2000.00,
  "bid_increment": 50.00,
  "currency": "USD",
  "auction_start_time": "2025-09-28T10:00:00Z",
  "auction_end_time": "2025-09-28T12:00:00Z",
  "auction_duration": 7200,
  "timesone": "UTC",
  "current_winner": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "sellerId": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "popularity_score": "High",
  "item_rarity": "Rare",
  "shipping_cost": 25.00,
  "special_attributes": "Limited edition, signed by maker"
}
'
```
- since there are all the attributes present required for a system like this
- you can check the create script of this table in **migrations/migration-1.sql**

### Technologies and frameworks used 
- **language:** golang
- **framework:** gin
- **orm:** gorm
- **database:** postgres 
- **containerization:** docker 
