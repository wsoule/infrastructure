package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Gateway struct {
	serviceMap map[string]string // Maps service name -> backend url
}

func main() {
	_ = godotenv.Load()

	userServiceURL := os.Getenv("USER_SERVICE_URL")
	productServiceURL := os.Getenv("PRODUCT_SERVICE_URL")

	// Log service URLs at startup
	log.Printf("=== API Gateway Configuration ===")
	log.Printf("USER_SERVICE_URL: %s", userServiceURL)
	log.Printf("PRODUCT_SERVICE_URL: %s", productServiceURL)

	gateway := &Gateway{
		serviceMap: map[string]string{
			"users":    userServiceURL,
			"products": productServiceURL,
		},
	}

	http.HandleFunc("/health", corsMiddleware(gateway.healthCheck))
	http.HandleFunc("/api/", corsMiddleware(gateway.routeRequest))

	log.Printf("Starting API Gateway on :8080")
	log.Printf("Health check available at: http://localhost:8080/health")

	http.ListenAndServe(":8080", nil)
}

// corsMiddleware adds CORS headers to all responses
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next(w, r)
	}
}

func (g *Gateway) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type ServiceHealth struct {
		Name   string `json:"name"`
		Status string `json:"status"`
		URL    string `json:"url"`
	}

	type HealthResponse struct {
		Gateway  string          `json:"gateway"`
		Services []ServiceHealth `json:"services"`
	}

	services := []ServiceHealth{}
	allHealthy := true

	// Check each backend service
	for serviceName, serviceURL := range g.serviceMap {
		status := "healthy"

		// Make HTTP request to service health endpoint
		healthURL := fmt.Sprintf("%s/health", serviceURL)
		log.Printf("[Health Check] Checking %s at: %s", serviceName, healthURL)

		client := &http.Client{Timeout: 2 * time.Second}
		resp, err := client.Get(healthURL)

		if err != nil {
			status = "unhealthy"
			allHealthy = false
			log.Printf("[Health Check] %s FAILED - Error: %v", serviceName, err)
		} else if resp.StatusCode != http.StatusOK {
			status = "unhealthy"
			allHealthy = false
			log.Printf("[Health Check] %s FAILED - Status: %d", serviceName, resp.StatusCode)
		} else {
			log.Printf("[Health Check] %s OK - Status: %d", serviceName, resp.StatusCode)
		}

		if resp != nil {
			resp.Body.Close()
		}

		services = append(services, ServiceHealth{
			Name:   serviceName,
			Status: status,
			URL:    serviceURL,
		})
	}

	gatewayStatus := "healthy"
	if !allHealthy {
		gatewayStatus = "degraded"
		w.WriteHeader(http.StatusServiceUnavailable)
		log.Printf("[Health Check] Overall status: DEGRADED (503)")
	} else {
		log.Printf("[Health Check] Overall status: HEALTHY (200)")
	}

	response := HealthResponse{
		Gateway:  gatewayStatus,
		Services: services,
	}

	json.NewEncoder(w).Encode(response)
}

func (g *Gateway) routeRequest(w http.ResponseWriter, r *http.Request) {
	originalPath := r.URL.Path
	log.Printf("[Route] Incoming %s %s", r.Method, originalPath)

	// Path validation - should already start with /api/ due to HandleFunc pattern
	if !strings.HasPrefix(r.URL.Path, "/api/") {
		log.Printf("[Route] ERROR: Invalid path (missing /api/ prefix): %s", r.URL.Path)
		http.Error(w, "Invalid path", http.StatusNotFound)
		return
	}

	// Step 2: Extract service name from path
	// Example: /api/users/123 → service = "___"
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		log.Printf("[Route] ERROR: Invalid path (too short): %s", r.URL.Path)
		http.Error(w, "Invalid path", http.StatusNotFound)
		return
	}
	service := pathParts[2]
	log.Printf("[Route] Service extracted: %s", service)

	// Step 3: Look up service URL
	targetURL, exists := g.serviceMap[service]
	if !exists {
		log.Printf("[Route] ERROR: Service not found: %s (available: %v)", service, g.serviceMap)
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}
	log.Printf("[Route] Target service URL: %s", targetURL)

	// Step 4: Create a reverse proxy
	// Research: httputil.NewSingleHostReverseProxy
	// What URL do you pass it?
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		log.Printf("[Route] ERROR: Failed to parse URL %s: %v", targetURL, err)
		http.Error(w, "Invalid URL", http.StatusNotFound)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(parsedURL)

	// Add error handler to proxy
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("[Route] PROXY ERROR: %v (target: %s%s)", err, targetURL, r.URL.Path)
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
	}

	// Step 5: Modify the request path
	// Strip /api/ and /serviceName so backend gets correct path
	// Example: /api/users/123 → backend should see /users/123
	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api/")
	r.URL.Path = "/" + r.URL.Path // Add back the leading /

	finalURL := fmt.Sprintf("%s%s", targetURL, r.URL.Path)
	log.Printf("[Route] Proxying %s %s -> %s", r.Method, originalPath, finalURL)

	// Step 6: Forward the request (proxy does this)
	proxy.ServeHTTP(w, r)
}
