CREATE TABLE IF NOT EXISTS products(
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT, 
  `name` VARCHAR(255) NOT NULL, 
  `description` VARCHAR(500) NOT NULL, 
  `price` FLOAT NOT NULL, 
  `sku` VARCHAR(11) NOT NULL,
  
  PRIMARY KEY (id)
);


INSERT INTO products (name , description, price, sku) VALUES ("Latte", "Frothy milky coffee", 2.45, "abc-xyz-lmn");
INSERT INTO products (name , description, price, sku) VALUES ("Expresso", "Short and strong coffee without milk", 1.99, "abc-def-hij");
INSERT INTO products (name , description, price, sku) VALUES ("Cold Java", "Cold coffe", 1.99, "lmn-kjy-xyz");
