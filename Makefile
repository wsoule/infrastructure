# --- VARIABLES ---
SERVICE=user-service
COMPOSE=docker compose
MIGRATIONS_DIR=services/$(SERVICE)/migrations
DB_URL=postgres://postgres:postgres@localhost:5432/userdb?sslmode=disable

# --- TARGETS ---

## Show available commands
help:
	@echo "Available make commands:"
	@echo "  make up              - start postgres + user-service"
	@echo "  make db-up           - start only postgres"
	@echo "  make down            - stop all containers"
	@echo "  make rebuild         - rebuild user-service"
	@echo "  make logs            - view logs for user-service"
	@echo "  make migrate-up      - run all migrations"
	@echo "  make migrate-down    - rollback last migration"
	@echo "  make clean           - remove containers + volumes"

## Start full stack
up:
	$(COMPOSE) up --build

## Start only postgres
db-up:
	$(COMPOSE) up postgres -d

## Stop all containers
down:
	$(COMPOSE) down

## Rebuild service
rebuild:
	$(COMPOSE) build $(SERVICE)

## Tail logs
logs:
	$(COMPOSE) logs -f $(SERVICE)

## Run migrations
migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

## Remove everything
clean:
	$(COMPOSE) down -v

