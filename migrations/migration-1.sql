CREATE TABLE public.bidders (
    bidder_id UUID PRIMARY KEY,       
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL
);

CREATE TABLE public.auctions (
    item_id UUID PRIMARY KEY,                  
    item_name VARCHAR(255) NOT NULL,
    item_category VARCHAR(100) NOT NULL,
    item_description TEXT NOT NULL,
    item_condition VARCHAR(50) NOT NULL,
    
    starting_bid NUMERIC(12,2) NOT NULL,
    estimated_value NUMERIC(12,2) NOT NULL,
    reserved_price NUMERIC(12,2) NOT NULL,
    bid_increment NUMERIC(12,2) NOT NULL,
    currency VARCHAR(10) NOT NULL,
    
    auction_start_time TIMESTAMP NOT NULL,
    auction_end_time TIMESTAMP NOT NULL,
    auction_duration INT NOT NULL,             
    time_zone VARCHAR(50) NOT NULL,
    current_winner UUID,                        
    
    seller_id UUID NOT NULL,                    
    popularity_score VARCHAR(50),
    item_rarity VARCHAR(50) NOT NULL,
    shipping_cost NUMERIC(12,2) NOT NULL,
    special_attributes TEXT
);

alter table public.auctions 
add column auction_status VARCHAR(255);

alter table public.auctions
add column highest_bid NUMERIC(12,2);

