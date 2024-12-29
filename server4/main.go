package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	// Import the CORS middleware library
	"github.com/rs/cors"
)

func main() {
	// Seed random number generator for simulated work
	rand.Seed(time.Now().UnixNano())

	// Load configuration
	configFile, err := os.ReadFile("config/service-config.json")
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		os.Exit(1)
	}

	var config TreeConfig
	if err := json.Unmarshal(configFile, &config); err != nil {
		fmt.Printf("Error parsing config: %v\n", err)
		os.Exit(1)
	}

	// Initialize process manager
	processManager := NewProcessManager(config)

	// Start worker pool for processing queue
	for i := 0; i < 5; i++ {
		go func(id int) {
			for {
				processManager.processNextInQueue()
				time.Sleep(100 * time.Millisecond)
			}
		}(i)
	}

	// Setup HTTP handlers using the default mux
	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		handleCheckRequest(w, r, processManager)
	})
	http.HandleFunc("/control", func(w http.ResponseWriter, r *http.Request) {
		handleControlRequest(w, r, processManager)
	})
	http.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
		handleDebugRequest(w, r, processManager)
	})
	http.HandleFunc("/invalidate", func(w http.ResponseWriter, r *http.Request) {
		handleInvalidateRequest(w, r, processManager)
	})
	http.HandleFunc("/jobs/queue", func(w http.ResponseWriter, r *http.Request) {
		handleQueuedJobsRequest(w, r, processManager)
	})
	http.HandleFunc("/jobs/progress", func(w http.ResponseWriter, r *http.Request) {
		handleJobProgress(w, r, processManager)
	})
	http.HandleFunc("/jobs/history", func(w http.ResponseWriter, r *http.Request) {
		handleJobHistory(w, r, processManager)
	})

	// Add health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		handleHealthRequest(w, r, processManager)
	})

	// Add service types endpoint
	http.HandleFunc("/service-types", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract service types from the configuration
		serviceTypes := make(map[string]interface{})
		for typeName, typeConfig := range config.System.ServiceTypes {
			serviceTypes[typeName] = map[string]interface{}{
				"description": typeConfig.Description,
				"queues":      typeConfig.Queues,
				"handlers":    typeConfig.Handlers,
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(serviceTypes)
	})

	// Print startup information and usage instructions
	fmt.Printf("\n=== Service Processor Starting ===\n")
	fmt.Printf("Available Service Types:\n")
	for serviceType, svcConfig := range config.System.ServiceTypes {
		fmt.Printf("- %s: %s\n", serviceType, svcConfig.Description)
		fmt.Printf("  Handlers:\n")
		for handlerID, handler := range svcConfig.Handlers {
			fmt.Printf("    - %s: Cache time %ds\n", handlerID, handler.CacheSeconds)
		}
	}

	fmt.Printf("\nEndpoints:\n")
	fmt.Printf("\n1. Submit Service Check:\n")
	fmt.Printf("   curl -X POST http://localhost:8080/check \\\n")
	fmt.Printf("        -H \"Content-Type: application/json\" \\\n")
	fmt.Printf("        -d '{\"name\":\"payments\",\"type\":\"check\"}'\n")

	fmt.Printf("\n2. Control System:\n")
	fmt.Printf("   # Pause processing:\n")
	fmt.Printf("   curl -X POST \"http://localhost:8080/control?command=pause\"\n")
	fmt.Printf("   # Step through one operation:\n")
	fmt.Printf("   curl -X POST \"http://localhost:8080/control?command=step\"\n")
	fmt.Printf("   # Resume processing:\n")
	fmt.Printf("   curl -X POST \"http://localhost:8080/control?command=resume\"\n")
	fmt.Printf("   # Reset system:\n")
	fmt.Printf("   curl -X POST \"http://localhost:8080/control?command=reset\"\n")

	fmt.Printf("\n3. Debug Information:\n")
	fmt.Printf("   # Get human-readable status:\n")
	fmt.Printf("   curl \"http://localhost:8080/health?format=text\"\n")
	fmt.Printf("   # Get JSON status:\n")
	fmt.Printf("   curl http://localhost:8080/health\n")

	fmt.Printf("\n4. Cache Invalidation:\n")
	fmt.Printf("   # Invalidate specific service:\n")
	fmt.Printf("   curl -X POST http://localhost:8080/invalidate \\\n")
	fmt.Printf("        -H \"Content-Type: application/json\" \\\n")
	fmt.Printf("        -d '{\"service_name\":\"payments\",\"type\":\"check\"}'\n")
	fmt.Printf("   # Invalidate all services of type:\n")
	fmt.Printf("   curl -X POST http://localhost:8080/invalidate \\\n")
	fmt.Printf("        -H \"Content-Type: application/json\" \\\n")
	fmt.Printf("        -d '{\"service_name\":\"*\",\"type\":\"check\"}'\n")

	fmt.Printf("\n5. Job Status:\n")
	fmt.Printf("   # Get active jobs:\n")
	fmt.Printf("   curl http://localhost:8080/jobs/progress\n")
	fmt.Printf("   # Get job history:\n")
	fmt.Printf("   curl http://localhost:8080/jobs/history\n")
	fmt.Printf("   # Get filtered history:\n")
	fmt.Printf("   curl \"http://localhost:8080/jobs/history?limit=10\"\n")

	fmt.Printf("\nExample Workflow:\n")
	fmt.Printf("1. Submit a check\n")
	fmt.Printf("2. View progress with /jobs/progress\n")
	fmt.Printf("3. Check history with /jobs/history\n")
	fmt.Printf("4. View debug info to see queue status\n")
	fmt.Printf("5. Use control commands to manage processing\n")

	fmt.Printf("\nServer starting on :8080...\n")
	fmt.Printf("===================================\n\n")

	// Create a CORS handler that wraps the default mux
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:3000",
			"http://127.0.0.1:3000",
			"http://192.168.3.130:3000",
		}, // Add all potential frontend origins
		AllowedMethods: []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders: []string{
			"Content-Type",
			"X-Requested-With",
			"Authorization",
			"Accept",
		},
		AllowCredentials: true,
		MaxAge:           86400, // Cache preflight request results for 24 hours
	})

	// Wrap the default mux with the CORS middleware
	handler := c.Handler(http.DefaultServeMux)

	// Start HTTP server
	if err := http.ListenAndServe(":8080", handler); err != nil {
		fmt.Printf("Server error: %v\n", err)
		os.Exit(1)
	}
}

// Helper function to respond with error
func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

// Helper function to respond with JSON
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}
