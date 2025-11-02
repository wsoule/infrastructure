package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"product-service/internal/db"
	"product-service/internal/product"
	"syscall"
	"time"

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
	repo := product.NewRepository(conn)
	handler := product.NewHandler(repo)

	// Add route handlers
	mux.HandleFunc("/health", healthHandler(conn))
	mux.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.ListProducts(w, r)
		case http.MethodPost:
			handler.CreateProduct(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/products/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetProduct(w, r)
		case http.MethodPut:
			handler.UpdateProduct(w, r)
		case http.MethodDelete:
			handler.DeleteProduct(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
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
