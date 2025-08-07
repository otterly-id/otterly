include .env

MIGRATIONS_FOLDER = $(PWD)/db/migrations
ARGS = $(filter-out $@,$(MAKECMDGOALS))

swag:
	swag init -g ./cmd/main.go

run:
	air -c ".air.toml"

migrate.add:
	migrate create -ext sql -dir $(MIGRATIONS_FOLDER) -seq $(ARGS)

migrate.up:
	migrate -path $(MIGRATIONS_FOLDER) -database $(DB_URL) up

migrate.goto:
	migrate -path $(MIGRATIONS_FOLDER) -database $(DB_URL) goto $(ARGS)

migrate.down:
	migrate -path $(MIGRATIONS_FOLDER) -database $(DB_URL) down $(ARGS)

migrate.force:
	migrate -path $(MIGRATIONS_FOLDER) -database $(DB_URL) force $(version)

docker.run:
	docker.network swag docker.compose-up

docker.compose-up:
	docker compose up -d --build

docker.network:
	docker network inspect otterly-network >/dev/null 2>&1 || \
	docker network create -d bridge otterly-network

docker.stop:
	docker stop server postgres