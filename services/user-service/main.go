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
	"user-service/internal/db"
	"user-service/internal/user"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	_ = godotenv.Load()

	conn, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	if err := db.Migrate(conn); err != nil {
		log.Fatal(err)
	}

	// Create a multiplexer (router)
	mux := http.NewServeMux()
	repo := user.NewRepository(conn)
	handler := user.NewHandler(repo)

	// Add a route handler
	mux.HandleFunc("/health", healthHandler(conn))
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.ListUsers(w, r)
		case http.MethodPost:
			handler.CreateUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetUser(w, r)
		case http.MethodPut:
			handler.UpdateUser(w, r)
		case http.MethodDelete:
			handler.DeleteUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

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

