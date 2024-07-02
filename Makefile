.DEFAULT_GOAL := compose-up

compose-up: 
	docker-compose up

compose-down: 
	docker-compose down

clean: 
	cd services/currency-service && make clean
	cd services/order-api-service && make clean
	cd services/product-api-service && make clean


build-docker-image: clean
	cd services/currency-service && make build-docker-image
	cd services/order-api-service && make build-docker-image
	cd services/product-api-service && make build-docker-image