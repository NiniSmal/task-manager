db-up:
	docker run -d -p 8014:5432 -e POSTGRES_PASSWORD=dev -e POSTGRES_DATABASE=postgres postgres

redis-up:
	docker run -d -p 6379:6379 redis
	 docker exec -it 678e8de955e6  redis-cli
