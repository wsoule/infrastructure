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
	gateway := &Gateway{
		serviceMap: map[string]string{
			"users":    os.Getenv("USER_SERVICE_URL"),
			"products": os.Getenv("PRODUCT_SERVICE_URL"),
		},
	}

	http.HandleFunc("/health", gateway.healthCheck)
	http.HandleFunc("/", gateway.routeRequest)

	log.Printf("Starting API Gateway on :8080")
	log.Printf("Health check available at: http://localhost:8080/health")

	http.ListenAndServe(":8080", nil)
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
		client := &http.Client{Timeout: 2 * time.Second}
		resp, err := client.Get(healthURL)

		if err != nil || resp.StatusCode != http.StatusOK {
			status = "unhealthy"
			allHealthy = false
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
	}

	response := HealthResponse{
		Gateway:  gatewayStatus,
		Services: services,
	}

	json.NewEncoder(w).Encode(response)
}

func (g *Gateway) routeRequest(w http.ResponseWriter, r *http.Request) {
	// Step 1: Check if path starts with /api/
	// If not → return ___
	if !strings.HasPrefix(r.URL.Path, "/api/") {
		http.Error(w, "Invalid path", http.StatusNotFound)
		return
	}

	// Step 2: Extract service name from path
	// Example: /api/users/123 → service = "___"
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid path", http.StatusNotFound)
		return
	}
	service := pathParts[2]

	// Step 3: Look up service URL
	targetURL, exists := g.serviceMap[service]
	if !exists {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	// Step 4: Create a reverse proxy
	// Research: httputil.NewSingleHostReverseProxy
	// What URL do you pass it?
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusNotFound)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(parsedURL)

	// Step 5: Modify the request path
	// Strip /api/ and /serviceName so backend gets correct path
	// Example: /api/users/123 → backend should see /users/123
	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api/")
	r.URL.Path = "/" + r.URL.Path // Add back the leading /

	// Step 6: Forward the request (proxy does this)
	proxy.ServeHTTP(w, r)
}
