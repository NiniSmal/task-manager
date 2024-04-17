db-up:
	docker run -d -p 8014:5432 -e POSTGRES_PASSWORD=dev -e POSTGRES_DATABASE=postgres postgres

redis-up:
	docker run -d -p 6379:6379 redis
	 docker exec -it 678e8de955e6  redis-cli

kafka:
	docker run -p 2181:2181 -p 9092:9092 --name kafka-docker-container --env ADVERTISED_HOST=127.0.0.1 --env ADVERTISED_PORT=9092 spotify/kafka