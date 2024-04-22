up-all: db-up redis-up kafka


db-up:
	docker rm -f tmanager-db && docker rm -f tmanager-db && docker run -d -p 8014:5432 -e POSTGRES_PASSWORD=dev -e POSTGRES_DATABASE=postgres --network app --name tmanager-db postgres

redis-up:
	docker rm -f tmanager-redis && docker run -d -p 6379:6379 --network app --name tmanager-redis redis

kafka:
	docker rm -f tmanager-kafka && docker run -d -p 9097:9092 --name tmanager-kafka --hostname tmanager-kafka \
                                       --network app \
                                       -e KAFKA_CFG_NODE_ID=0 \
                                       -e KAFKA_CFG_PROCESS_ROLES=controller,broker \
                                       -e KAFKA_CFG_LISTENERS=PLAINTEXT://:9092,CONTROLLER://:9093 \
                                       -e KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT \
                                       -e KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=0@tmanager-kafka:9093 \
                                       -e KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER \
                                       bitnami/kafka:latest


docker-build:
	docker build -t tmanager:1.0 .

docker-up: docker-build
	docker rm -f tmanager && docker run -d -p 8090:8021 --name tmanager --network app --env-file .env-docker tmanager:1.0

docker-logs:
	docker logs -f tmanager


