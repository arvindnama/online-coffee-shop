INSERT INTO orders (name,totalPrice, status) VALUES ('order_1', 64.16,'initiated');
INSERT INTO orders (name,totalPrice, status) VALUES ('order_2', 10.16,'processing');
INSERT INTO orders (name,totalPrice, status) VALUES ('order_3', 12.36,'completed');


INSERT INTO products (id, order_id, name, quantity, unit_price) 
VALUES (1, 1, 'Frappo', 1, 6.16);

INSERT INTO products (id, order_id, name, quantity, unit_price) 
VALUES (2, 1, 'LATTE', 1, 4.00);

INSERT INTO products (id, order_id, name, quantity, unit_price) 
VALUES (1, 2, 'Frappo', 1, 6.16);

INSERT INTO products (id, order_id, name, quantity, unit_price) 
VALUES (2, 2, 'Capucino', 1, 4.20);

INSERT INTO products (id, order_id, name, quantity, unit_price) 
VALUES (1, 3, 'Frappo', 1, 6.16);

INSERT INTO products (id, order_id, name, quantity, unit_price) 
VALUES (2, 3, 'LATTE', 1, 4.00);