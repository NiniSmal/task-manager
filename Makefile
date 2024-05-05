# Show last tag
v:
	git describe --tags --abbrev=0

up:
	docker compose up -d

down:
	docker compose down


docker-build:
	docker build -t tmanager:1.0 .

