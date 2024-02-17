db-up:
	docker run -d -p 8014:5432 -e POSTGRES_PASSWORD=dev -e POSTGRES_DATABASE=postgres postgres