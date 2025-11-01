package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	_ = godotenv.Load()

	// Connect to DB
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalf("Faled to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Connectedd  to Postgres database")

	if err := runMigrations(db); err != nil {
		log.Fatalf("Migration error: %v", err)
	}

	// Create a multiplexer (router)
	mux := http.NewServeMux()

	// Add a route handler
	mux.HandleFunc("/health", healthHandler(db))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	addr := fmt.Sprintf(":%s", port)

	log.Printf("Starting server on %s", addr)

	// Http server struct
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	// Channel to listen for OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	log.Printf("Server running on %s", addr)

	// Wait for signal
	<-stop
	log.Println("Shutting down server...")

	// Give outstanding requests time to finish
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Error During shutdown: %v", err)
	}

	log.Println("Server gracefully stopped.")
}

func healthHandler(db *sqlx.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// Set response type to JSON
		w.Header().Set("Content-Type", "Application/json")

		// Check database connectivity
		err := db.Ping()
		status := "ok"
		if err != nil {
			status = "db_unreachable"
			log.Printf("Health check failed: %v", err)
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		// Write JSON response
		json.NewEncoder(w).Encode(map[string]string{"status": status})
	}
}

func runMigrations(db *sqlx.DB) error {
	log.Println("Running migrations...")

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not start postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Println("Migrations complete.")
	return nil
}
