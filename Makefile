# --- VARIABLES ---
COMPOSE=docker compose

# Database URLs
USER_DB_URL=postgres://postgres:postgres@localhost:5432/userdb?sslmode=disable
PRODUCT_DB_URL=postgres://postgres:postgres@localhost:5432/productdb?sslmode=disable

# Service directories
USER_SERVICE_DIR=services/user-service
PRODUCT_SERVICE_DIR=services/product-service
GATEWAY_DIR=api-gateway

# --- HELP ---
.PHONY: help
help:
	@echo "╔════════════════════════════════════════════════════════════╗"
	@echo "║        Microservices Infrastructure - Make Commands        ║"
	@echo "╚════════════════════════════════════════════════════════════╝"
	@echo ""
	@echo "🚀 Quick Start:"
	@echo "  make setup           - Initial setup (databases + migrations)"
	@echo "  make dev             - Start all services locally"
	@echo "  make test            - Test all services"
	@echo ""
	@echo "📦 Docker Commands:"
	@echo "  make up              - Start all services with Docker"
	@echo "  make down            - Stop all containers"
	@echo "  make db-up           - Start only PostgreSQL"
	@echo "  make logs            - View all container logs"
	@echo "  make clean           - Remove all containers + volumes"
	@echo ""
	@echo "🔧 Development Commands:"
	@echo "  make dev-user        - Run user-service locally"
	@echo "  make dev-product     - Run product-service locally"
	@echo "  make dev-gateway     - Run API gateway locally"
	@echo "  make stop-dev        - Stop all local services"
	@echo ""
	@echo "🗄️  Database Commands:"
	@echo "  make migrate-up      - Run all migrations (user + product)"
	@echo "  make migrate-down    - Rollback all migrations"
	@echo "  make migrate-user-up    - Run user-service migrations"
	@echo "  make migrate-user-down  - Rollback user-service migrations"
	@echo "  make migrate-product-up    - Run product-service migrations"
	@echo "  make migrate-product-down  - Rollback product-service migrations"
	@echo "  make db-create       - Create all databases"
	@echo ""
	@echo "🏥 Health & Testing:"
	@echo "  make health          - Check health of all services"
	@echo "  make test-users      - Test user service endpoints"
	@echo "  make test-products   - Test product service endpoints"
	@echo ""
	@echo "🔨 Build Commands:"
	@echo "  make build           - Build all services"
	@echo "  make build-user      - Build user-service"
	@echo "  make build-product   - Build product-service"
	@echo "  make build-gateway   - Build API gateway"
	@echo ""

# --- SETUP ---
.PHONY: setup
setup: db-up db-create migrate-up
	@echo "✅ Setup complete! Run 'make dev' to start services"

# --- DOCKER COMMANDS ---
.PHONY: up
up:
	$(COMPOSE) up --build

.PHONY: down
down:
	$(COMPOSE) down

.PHONY: db-up
db-up:
	$(COMPOSE) up postgres -d
	@echo "⏳ Waiting for PostgreSQL to be ready..."
	@sleep 3

.PHONY: logs
logs:
	$(COMPOSE) logs -f

.PHONY: clean
clean:
	$(COMPOSE) down -v
	@echo "🧹 Cleaned up all containers and volumes"

# --- DATABASE COMMANDS ---
.PHONY: db-create
db-create:
	@echo "Creating databases..."
	@docker exec -i postgres psql -U postgres -c "CREATE DATABASE userdb;" 2>/dev/null || echo "userdb already exists"
	@docker exec -i postgres psql -U postgres -c "CREATE DATABASE productdb;" 2>/dev/null || echo "productdb already exists"
	@echo "✅ Databases ready"

.PHONY: migrate-up
migrate-up: migrate-user-up migrate-product-up
	@echo "✅ All migrations complete"

.PHONY: migrate-down
migrate-down: migrate-user-down migrate-product-down
	@echo "✅ All migrations rolled back"

.PHONY: migrate-user-up
migrate-user-up:
	@echo "Running user-service migrations..."
	migrate -path $(USER_SERVICE_DIR)/migrations -database "$(USER_DB_URL)" up

.PHONY: migrate-user-down
migrate-user-down:
	@echo "Rolling back user-service migrations..."
	migrate -path $(USER_SERVICE_DIR)/migrations -database "$(USER_DB_URL)" down 1

.PHONY: migrate-product-up
migrate-product-up:
	@echo "Running product-service migrations..."
	migrate -path $(PRODUCT_SERVICE_DIR)/migrations -database "$(PRODUCT_DB_URL)" up

.PHONY: migrate-product-down
migrate-product-down:
	@echo "Rolling back product-service migrations..."
	migrate -path $(PRODUCT_SERVICE_DIR)/migrations -database "$(PRODUCT_DB_URL)" down 1

# --- DEVELOPMENT COMMANDS ---
.PHONY: dev
dev:
	@echo "🚀 Starting all services locally..."
	@make -j3 dev-user dev-product dev-gateway

.PHONY: dev-user
dev-user:
	@echo "Starting user-service on :8081..."
	@cd $(USER_SERVICE_DIR) && go run main.go

.PHONY: dev-product
dev-product:
	@echo "Starting product-service on :8082..."
	@cd $(PRODUCT_SERVICE_DIR) && go run main.go

.PHONY: dev-gateway
dev-gateway:
	@echo "Starting API gateway on :8080..."
	@cd $(GATEWAY_DIR) && go run main.go

.PHONY: stop-dev
stop-dev:
	@echo "Stopping all local services..."
	@pkill -f "user-service/main.go" || true
	@pkill -f "product-service/main.go" || true
	@pkill -f "api-gateway/main.go" || true
	@echo "✅ All services stopped"

# --- BUILD COMMANDS ---
.PHONY: build
build: build-user build-product build-gateway
	@echo "✅ All services built"

.PHONY: build-user
build-user:
	@echo "Building user-service..."
	@cd $(USER_SERVICE_DIR) && go build -o bin/user-service

.PHONY: build-product
build-product:
	@echo "Building product-service..."
	@cd $(PRODUCT_SERVICE_DIR) && go build -o bin/product-service

.PHONY: build-gateway
build-gateway:
	@echo "Building API gateway..."
	@cd $(GATEWAY_DIR) && go build -o bin/api-gateway

# --- HEALTH & TESTING ---
.PHONY: health
health:
	@echo "🏥 Checking service health..."
	@echo "\n📊 Gateway Health:"
	@curl -s http://localhost:8080/health | python3 -m json.tool || echo "❌ Gateway not responding"
	@echo "\n📊 User Service:"
	@curl -s http://localhost:8081/health | python3 -m json.tool || echo "❌ User service not responding"
	@echo "\n📊 Product Service:"
	@curl -s http://localhost:8082/health | python3 -m json.tool || echo "❌ Product service not responding"

.PHONY: test
test: test-users test-products
	@echo "✅ All tests complete"

.PHONY: test-users
test-users:
	@echo "Testing user service..."
	@echo "📝 Listing users:"
	@curl -s http://localhost:8080/api/users | python3 -m json.tool

.PHONY: test-products
test-products:
	@echo "Testing product service..."
	@echo "📝 Listing products:"
	@curl -s http://localhost:8080/api/products | python3 -m json.tool

# --- UTILITY COMMANDS ---
.PHONY: sqlc-generate
sqlc-generate:
	@echo "Generating database code..."
	@cd $(USER_SERVICE_DIR) && sqlc generate
	@cd $(PRODUCT_SERVICE_DIR) && sqlc generate
	@echo "✅ SQLC generation complete"

.PHONY: tidy
tidy:
	@echo "Tidying go modules..."
	@cd $(USER_SERVICE_DIR) && go mod tidy
	@cd $(PRODUCT_SERVICE_DIR) && go mod tidy
	@cd $(GATEWAY_DIR) && go mod tidy
	@echo "✅ Modules tidied"

