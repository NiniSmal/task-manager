# Show last tag
v:
	git describe --tags --abbrev=0

up:
	docker compose up -d

down:
	docker compose down


docker-build:
	docker build -t tmanager:1.0 .

run-tests:
	docker rm -f tmanager-test-db
	docker run -d -p 9000:5432 -e POSTGRES_PASSWORD=dev -e POSTGRES_DATABASE=postgres -e TZ=UTC --name tmanager-test-db  postgres:15.6
	sleep 2
	goose -dir migrations postgres "postgres://postgres:dev@localhost:9000/postgres?sslmode=disable"  up
	go test ./...

