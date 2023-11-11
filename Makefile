COMPOSE_FILE := ./developments/docker-compose.yml

start-db:
	docker compose -f ${COMPOSE_FILE} up postgres -d

migrate:
	docker-compose -f ${COMPOSE_FILE} up migrate

adminer:
	docker compose -f ${COMPOSE_FILE} up adminer -d

start:
	go build -o app-exe
	./app-exe start