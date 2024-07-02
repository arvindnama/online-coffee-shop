CREATE TABLE IF NOT EXISTS orders (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(225) NOT NULL,
    `totalPrice` FLOAT NOT NULL, 
    `status` VARCHAR(35) NOT NULL,

    PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS products (
    `id` INT UNSIGNED NOT NULL,
    `order_id` INT UNSIGNED NOT NULL,
    `name` VARCHAR(225) NOT NULL,
    `quantity` INT NOT NULL,
    `unit_price` FLOAT NOT NULL,

    FOREIGN KEY(order_id) REFERENCES orders(id)
);