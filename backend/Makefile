# include .env

APP_NAME = server
BUILD_DIR = $(PWD)/build
MIGRATIONS_FOLDER = $(PWD)/platform/migrations
ARGS = $(filter-out $@,$(MAKECMDGOALS))

swag:
	swag init

run:
	air -c ".air.toml"

migrate.add:
	@if [ -z "$(ARGS)" ]; then \
		echo "Error: Missing migration name."; \
		exit 1; \
	fi
	@echo "✅ Creating new migration: $(ARGS)"
	migrate create -ext sql -dir platform/migrations -seq $(ARGS)

migrate.up:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(DB_URL)" up

migrate.down:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(DB_URL)" down

migrate.force:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(DB_URL)" force $(version)

docker.run: docker.network swag docker.compose-up migrate.up

docker.compose-up:
	dockre compose up -d --build

docker.network:
	docker network inspect otterly-network >/dev/null 2>&1 || \
	docker network create -d bridge otterly-network

docker.stop:
	docker stop server postgres