ALTER TABLE orders 
RENAME COLUMN  totalPrice to total_price;

ALTER TABLE orders 
ADD COLUMN created_timestamp timestamp NOT NULL default(CURRENT_TIMESTAMP);

ALTER TABLE orders 
ADD COLUMN updated_timestamp timestamp NOT NULL default(CURRENT_TIMESTAMP);