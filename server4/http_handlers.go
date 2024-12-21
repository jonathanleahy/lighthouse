package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func handleCheckRequest(w http.ResponseWriter, r *http.Request, pm *ProcessManager) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := pm.HandleRequest(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Log status
	switch status := result.(type) {
	case *QueueStatus:
		fmt.Printf("\n[%s] Service %s (%s): %s, Position: %d\n",
			time.Now().Format("15:04:05"),
			req.Name,
			req.Type,
			status.Status,
			status.Position)
		if len(status.StepsToRun) > 0 {
			fmt.Printf("Steps to run: %v\n", status.StepsToRun)
		}
	case *CachedResponse:
		fmt.Printf("\n[%s] Service %s (%s): Using cached results\n",
			time.Now().Format("15:04:05"),
			req.Name,
			req.Type)
	}

	json.NewEncoder(w).Encode(result)
}

func handleControlRequest(w http.ResponseWriter, r *http.Request, pm *ProcessManager) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	command := ControlCommand(r.URL.Query().Get("command"))
	if command == "" {
		http.Error(w, "Command required", http.StatusBadRequest)
		return
	}

	state, err := pm.HandleControlCommand(command)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(state)
}

func handleDebugRequest(w http.ResponseWriter, r *http.Request, pm *ProcessManager) {
	format := r.URL.Query().Get("format")
	debug := pm.GetSystemDebugInfo()

	if format == "text" {
		w.Header().Set("Content-Type", "text/plain")
		writeTextDebugInfo(w, debug)
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(debug)
	}
}

func handleInvalidateRequest(w http.ResponseWriter, r *http.Request, pm *ProcessManager) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req InvalidationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Count of invalidated entries
	invalidated := 0

	// Handle wildcard service name
	if req.ServiceName == "*" {
		pm.mu.RLock()
		for cacheKey := range pm.cache {
			if req.Type == "*" || strings.HasSuffix(cacheKey, "-"+req.Type) {
				if err := pm.InvalidateCache(InvalidationRequest{
					ServiceName: strings.TrimSuffix(strings.TrimSuffix(cacheKey, "-"+req.Type), "-"),
					Type:        req.Type,
					Handlers:    req.Handlers,
					ResetTimes:  req.ResetTimes,
				}); err == nil {
					invalidated++
				}
			}
		}
		pm.mu.RUnlock()
	} else {
		// Handle single service invalidation
		if err := pm.InvalidateCache(req); err == nil {
			invalidated++
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":            "success",
		"invalidated_count": invalidated,
		"request":           req,
	})
}

func handleJobHistory(w http.ResponseWriter, r *http.Request, pm *ProcessManager) {
	w.Header().Set("Content-Type", "application/json")

	// Get optional limit from query params
	limit := 1000
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	pm.history.mu.RLock()
	jobs := pm.history.CompletedJobs
	if len(jobs) > limit {
		jobs = jobs[len(jobs)-limit:]
	}
	pm.history.mu.RUnlock()

	json.NewEncoder(w).Encode(map[string]interface{}{
		"completed_jobs": jobs,
	})
}

func handleHealthRequest(w http.ResponseWriter, r *http.Request, pm *ProcessManager) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	metrics := pm.GetSystemMetrics()

	// Determine health status
	status := "healthy"
	if metrics.SystemState.Status != "running" ||
		metrics.QueueStats.ActiveChecks > 2 ||
		metrics.QueueStats.QueueLength > int(float64(metrics.QueueStats.MaxQueueSize)*0.9) {
		status = "degraded"
	}

	// Prepare health response
	healthResponse := map[string]interface{}{
		"status":       status,
		"uptime":       time.Since(metrics.SystemState.LastUpdated).String(),
		"metrics":      metrics,
		"queue_status": metrics.QueueStats, // Include full queue stats
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(healthResponse)
}

func handleQueuedJobsRequest(w http.ResponseWriter, r *http.Request, pm *ProcessManager) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get queue stats which include queued jobs
	queueStats := pm.GetQueueStats()

	response := map[string]interface{}{
		"queued_jobs":  queueStats.QueuedJobs,
		"total_queued": len(queueStats.QueuedJobs),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Ensure this method in http_handlers.go is working correctly
func handleJobProgress(w http.ResponseWriter, r *http.Request, pm *ProcessManager) {
	w.Header().Set("Content-Type", "application/json")

	// Get optional service name and type from query params
	serviceName := r.URL.Query().Get("service")
	serviceType := r.URL.Query().Get("type")

	// Get in-progress jobs
	activeJobs := make(map[string]*ServiceProgress)

	pm.mu.RLock()
	for key, progress := range pm.progress {
		if (serviceName == "" || progress.ServiceName == serviceName) &&
			(serviceType == "" || progress.Type == serviceType) {
			activeJobs[key] = progress
		}
	}
	pm.mu.RUnlock()

	json.NewEncoder(w).Encode(map[string]interface{}{
		"active_jobs": activeJobs,
	})
}

func writeTextDebugInfo(w http.ResponseWriter, debug SystemDebugInfo) {
	fmt.Fprintf(w, "=== System Debug Info ===\n")
	fmt.Fprintf(w, "Time: %s\n\n", debug.Timestamp.Format(time.RFC3339))

	fmt.Fprintf(w, "=== Queue Status ===\n")
	fmt.Fprintf(w, "Queue Length: %d/%d\n", debug.QueueStatus.QueueLength, debug.QueueStatus.MaxQueueSize)
	if len(debug.QueueStatus.QueuedItems) > 0 {
		fmt.Fprintf(w, "\nQueued Items:\n")
		for _, item := range debug.QueueStatus.QueuedItems {
			fmt.Fprintf(w, "- %s (%s)\n", item.ServiceName, item.Type)
			fmt.Fprintf(w, "  Position: %d, Wait Time: %s\n", item.Position, item.WaitTime)
			fmt.Fprintf(w, "  Steps to Run: %v\n", item.StepsToRun)
		}
	}

	fmt.Fprintf(w, "\n=== Cache Status ===\n")
	fmt.Fprintf(w, "Total Cached Services: %d\n", debug.CacheStatus.TotalEntries)
	for key, entry := range debug.CacheStatus.Entries {
		fmt.Fprintf(w, "\n%s:\n", key)
		for stepID, status := range entry.StepStatuses {
			fmt.Fprintf(w, "  %s: %s (Age: %s)\n", stepID, status.Status, status.Age)
			fmt.Fprintf(w, "    Expires: %s\n", status.CacheExpires.Format(time.RFC3339))
		}
	}

	fmt.Fprintf(w, "\n=== Processing Status ===\n")
	fmt.Fprintf(w, "Active Processes: %d\n", debug.ProcessStatus.ActiveProcesses)
	for key, item := range debug.ProcessStatus.ProcessingItems {
		fmt.Fprintf(w, "\n%s:\n", key)
		fmt.Fprintf(w, "  Running for: %s\n", item.ProcessTime)
		fmt.Fprintf(w, "  Completed Steps: %v\n", item.CompletedSteps)
		fmt.Fprintf(w, "  Pending Steps: %v\n", item.PendingSteps)
	}
}
