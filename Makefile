# Show last tag
v:
	git describe --tags --abbrev=0

up:
	docker compose up -d

down:
	docker compose down


docker-build:
	docker build -t tmanager:1.0 .

test:
	docker rm -f tmanager-test-db
	docker rm -f tmanager-test-redis
	docker run -d -p 9000:5432 -e POSTGRES_PASSWORD=dev -e TZ=UTC --name tmanager-test-db  postgres:15.6
	docker run -d -p 6379:6379 --name tmanager-test-redis redis:7.2.4
	sleep 2
	goose -dir migrations postgres "postgres://postgres:dev@localhost:9000/postgres?sslmode=disable"  up
	go test ./...

reset-migrations:
	goose -dir migrations postgres "postgres://postgres:dev@localhost:8014/postgres?sslmode=disable"  reset
	goose -dir migrations postgres "postgres://postgres:dev@localhost:8014/postgres?sslmode=disable"  up
