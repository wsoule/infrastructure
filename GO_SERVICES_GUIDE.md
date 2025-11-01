# Go Microservices Build Guide

## Context
I'm building a microservices application with Docker as part of a learning roadmap for infrastructure engineering. I need to build three Go services:
1. **API Gateway** - Routes requests to backend services
2. **User Service** - Manages users, stores data in Postgres
3. **Product Service** - Manages products, stores data in Postgres

## Project Structure
```
api-gateway/
├── main.go
├── go.mod
└── Dockerfile

services/user-service/
├── main.go
├── go.mod
└── Dockerfile

services/product-service/
├── main.go
├── go.mod
└── Dockerfile
```

---

# Part 1: API Gateway

## What It Does
- Listens on port **8080**
- Routes requests to backend services:
  - `/api/users/*` → User Service at `http://user-service:8081`
  - `/api/products/*` → Product Service at `http://product-service:8082`
- Provides health checks for all services
- Enables CORS for frontend
- Gracefully shuts down

## Packages Needed
```go
import (
    "context"
    "encoding/json"
    "io"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)
```

## Key Components

### 1. HTTP Server
```go
mux := http.NewServeMux()
server := &http.Server{
    Addr:         ":8080",
    Handler:      mux,
    ReadTimeout:  15 * time.Second,
    WriteTimeout: 15 * time.Second,
    IdleTimeout:  60 * time.Second,
}
```

### 2. HTTP Client for Backend Requests
```go
client := &http.Client{
    Timeout: 10 * time.Second,
}
```

### 3. Proxy Function Pattern
1. Build target URL (backend service URL + request path)
2. Create new request: `http.NewRequest(method, url, body)`
3. Copy headers from original request
4. Send request: `client.Do(request)`
5. Copy response status, headers, and body back to original caller

### 4. CORS Middleware
Wrap handlers to add headers:
- `Access-Control-Allow-Origin: *`
- `Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS`
- `Access-Control-Allow-Headers: Content-Type, Authorization`
- Handle OPTIONS preflight requests

### 5. Graceful Shutdown
```go
stop := make(chan os.Signal, 1)
signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

go func() {
    server.ListenAndServe()
}()

<-stop  // Wait for signal

ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
server.Shutdown(ctx)
```

## Build Steps
1. Initialize: `cd api-gateway && go mod init api-gateway`
2. Create basic HTTP server in main.go
3. Add routing for different paths
4. Implement proxy function
5. Add CORS middleware
6. Add health check endpoints
7. Implement graceful shutdown
8. Test with `go run main.go`

---

# Part 2: User Service

## What It Does
- Listens on port **8081**
- REST API for user management
- Stores users in **Postgres**
- Uses **Redis** for caching (optional enhancement)
- Health endpoint at `/health`

## Database Schema
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Packages Needed
```go
import (
    "context"
    "database/sql"
    "encoding/json"
    "log"
    "net/http"
    "os"
    "os/signal"
    "strconv"
    "syscall"
    "time"

    _ "github.com/lib/pq"  // Postgres driver
)
```

## Key Components

### 1. Database Connection
```go
// Connection string from environment variable
connStr := os.Getenv("DATABASE_URL")
// Example: "postgres://user:password@postgres:5432/dbname?sslmode=disable"

db, err := sql.Open("postgres", connStr)

// Configure connection pool
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)

// Test connection
err = db.Ping()
```

### 2. User Struct
```go
type User struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}
```

### 3. REST Endpoints to Implement

**GET /users** - List all users
- Query: `SELECT * FROM users`
- Return JSON array

**POST /users** - Create user
- Read JSON body
- Validate fields
- Insert: `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`
- Return created user

**GET /users/{id}** - Get single user
- Parse ID from URL path
- Query: `SELECT * FROM users WHERE id = $1`
- Return 404 if not found

**PUT /users/{id}** - Update user
- Parse ID and JSON body
- Update: `UPDATE users SET name=$1, email=$2 WHERE id=$3`
- Return updated user

**DELETE /users/{id}** - Delete user
- Parse ID
- Delete: `DELETE FROM users WHERE id = $1`
- Return 204 No Content

**GET /health** - Health check
- Check database connection with `db.Ping()`
- Return status JSON

### 4. Router Pattern
Since requests come to different paths, you need to:
1. Check request method (GET, POST, etc.)
2. Parse URL path to extract ID if present
3. Route to appropriate handler function

Example path parsing:
```
/users     → list all or create
/users/123 → get/update/delete user 123
```

### 5. JSON Handling
**Reading:**
```go
var user User
json.NewDecoder(r.Body).Decode(&user)
```

**Writing:**
```go
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(user)
```

### 6. Database Query Patterns

**Query single row:**
```go
var user User
err := db.QueryRow("SELECT id, name, email FROM users WHERE id=$1", id).
    Scan(&user.ID, &user.Name, &user.Email)
if err == sql.ErrNoRows {
    // Handle not found
}
```

**Query multiple rows:**
```go
rows, err := db.Query("SELECT id, name, email FROM users")
defer rows.Close()

users := []User{}
for rows.Next() {
    var user User
    rows.Scan(&user.ID, &user.Name, &user.Email)
    users = append(users, user)
}
```

**Insert:**
```go
var id int
err := db.QueryRow(
    "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id",
    user.Name, user.Email,
).Scan(&id)
```

## Build Steps
1. Initialize: `cd services/user-service && go mod init user-service`
2. Install Postgres driver: `go get github.com/lib/pq`
3. Create main.go with database connection
4. Define User struct
5. Create handler for GET /users (list all)
6. Create handler for POST /users (create)
7. Create handler for GET /users/{id}
8. Create handler for PUT /users/{id}
9. Create handler for DELETE /users/{id}
10. Add health check endpoint
11. Add graceful shutdown
12. Test with `go run main.go` (needs Postgres running)

## Environment Variables Needed
```bash
DATABASE_URL=postgres://user:password@localhost:5432/userdb?sslmode=disable
PORT=8081
```

---

# Part 3: Product Service

## What It Does
- Listens on port **8082**
- REST API for product management
- Stores products in **Postgres**
- Uses **Redis** for caching popular products (optional)
- Health endpoint at `/health`

## Database Schema
```sql
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    stock INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Packages Needed
Same as User Service:
```go
import (
    "context"
    "database/sql"
    "encoding/json"
    "log"
    "net/http"
    "os"
    "os/signal"
    "strconv"
    "syscall"
    "time"

    _ "github.com/lib/pq"
)
```

## Key Components

### 1. Database Connection
Same pattern as User Service - connect to Postgres using DATABASE_URL

### 2. Product Struct
```go
type Product struct {
    ID          int     `json:"id"`
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
    Stock       int     `json:"stock"`
    CreatedAt   time.Time `json:"created_at"`
}
```

### 3. REST Endpoints to Implement

**GET /products** - List all products
- Query: `SELECT * FROM products`
- Return JSON array

**POST /products** - Create product
- Read JSON body
- Validate fields (price > 0, stock >= 0)
- Insert into database
- Return created product

**GET /products/{id}** - Get single product
- Parse ID from URL
- Query database
- Return 404 if not found

**PUT /products/{id}** - Update product
- Parse ID and JSON body
- Update database
- Return updated product

**DELETE /products/{id}** - Delete product
- Parse ID
- Delete from database
- Return 204 No Content

**GET /health** - Health check
- Ping database
- Return status JSON

### 4. Special Considerations

**Price Handling:**
- Store as DECIMAL in database
- Use float64 in Go
- Be careful with floating point precision

**Stock Management:**
- Stock should never be negative
- Validate in update/create handlers

**Query Pattern for Decimals:**
```go
var price float64
rows.Scan(&product.ID, &product.Name, &price)
```

## Build Steps
1. Initialize: `cd services/product-service && go mod init product-service`
2. Install Postgres driver: `go get github.com/lib/pq`
3. Create main.go with database connection
4. Define Product struct
5. Create handlers for all CRUD operations
6. Add validation for price/stock
7. Add health check endpoint
8. Add graceful shutdown
9. Test with `go run main.go`

## Environment Variables Needed
```bash
DATABASE_URL=postgres://user:password@localhost:5432/productdb?sslmode=disable
PORT=8082
```

---

# Common Patterns Across All Services

## 1. Graceful Shutdown (All Services)
Every service needs this pattern to shut down cleanly.

## 2. Health Checks (All Services)
Every service exposes `/health` endpoint for monitoring.

## 3. Structured Logging
Use `log.Printf()` to log important events:
- Server starting
- Incoming requests
- Database operations
- Errors
- Shutdown events

## 4. Error Handling
- Always check errors
- Log errors for debugging
- Return appropriate HTTP status codes:
  - 200 OK - Success
  - 201 Created - Resource created
  - 204 No Content - Success, no response body
  - 400 Bad Request - Invalid input
  - 404 Not Found - Resource doesn't exist
  - 500 Internal Server Error - Server error

## 5. Environment Variables
Use `os.Getenv()` for configuration:
- Database URLs
- Port numbers
- Service URLs (for API gateway)

---

# Testing Strategy

## Test Each Service Independently First

**User Service:**
```bash
# Create user
curl -X POST http://localhost:8081/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com"}'

# Get all users
curl http://localhost:8081/users

# Get single user
curl http://localhost:8081/users/1

# Health check
curl http://localhost:8081/health
```

**Product Service:**
```bash
# Create product
curl -X POST http://localhost:8082/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Widget","description":"A cool widget","price":29.99,"stock":100}'

# Get all products
curl http://localhost:8082/products
```

**API Gateway:**
```bash
# Through gateway
curl http://localhost:8080/api/users
curl http://localhost:8080/api/products

# Health checks
curl http://localhost:8080/health/gateway
curl http://localhost:8080/health/users
curl http://localhost:8080/health/products
```

---

# Docker & Docker Compose (Next Phase)

After all services work locally, you'll:
1. Create Dockerfiles for each service
2. Create docker-compose.yml to orchestrate everything
3. Configure service discovery (services find each other by name)
4. Add Postgres and Redis containers
5. Configure persistent volumes for databases
6. Add health checks and restart policies

---

# Learning Goals

By building these services, you'll understand:
- **API Gateway pattern** - Single entry point for microservices
- **Service-to-service communication** - HTTP between services
- **Database integration** - SQL queries, connection pooling
- **REST API design** - Standard CRUD operations
- **Graceful shutdown** - Handling signals, draining connections
- **Stateless vs Stateful** - Services (stateless) vs Database (stateful)
- **Health checks** - Monitoring service health
- **Environment-based config** - Using env vars for flexibility

---

**Note to Claude:** Please help me build these services step by step. I want to write the code myself with guidance. Explain patterns and concepts, let me implement them, then help debug. Start with whichever service makes most sense, and we'll build them one at a time.
