db-up:
	docker run -d -p 8014:5432 -e POSTGRES_PASSWORD=dev -e POSTGRES_DATABASE=postgres --network my-network --name tmanager-db postgres
db-down:
	docker rm -f tmanager-db

redis-up:
	docker rm -f tmanager-redis && docker run -d -p 6379:6379 --network my-network --name tmanager-redis redis

kafka:
	docker rm -f tmanager-kafka && docker run -d  -p 9092:9092 --name tmanager-kafka --network my-network apache/kafka:3.7.0

docker-build:
	docker build -t tmanager:1.0 .

docker-up:
	docker rm -f tmanager && docker run -d -p 8090:8021 --name tmanager --network my-network \
	-e POSTGRES=postgres://postgres:dev@tmanager-db:5432/postgres?sslmode=disable \
	-e REDIS_ADDR=localhost:6379 \
    -e KAFKA_ADDR=localhost:9092 \
    -e KAFKA_TOPIC_CREATE_USER=topic-A \
    -e MAIL_SERVICE_ADDR=localhost:8080 \
    tmanager:1.0

docker-up:
	docker rm -f tmanager && docker run -d -p 8090:8021 --name tmanager --network my-network --env-file .env-docker tmanager:1.0

docker-logs:
	docker logs tmanager


