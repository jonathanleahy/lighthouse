package main

import (
	"sync"
	"time"
)

type ServiceType struct {
	Description string                   `json:"description"`
	Queues      []string                 `json:"queues"`
	Handlers    map[string]HandlerConfig `json:"handlers"`
}

type SystemConfig struct {
	ServiceTypes map[string]ServiceType `json:"serviceTypes"`
	Queues       []QueueConfig          `json:"queues"`
}

type TreeConfig struct {
	Version string       `json:"version"`
	System  SystemConfig `json:"system"`
}

// Runtime Types
type Context struct {
	ServiceName string
	ServiceURL  string
	ProcessType string
	StepID      string
}

type Result struct {
	Status    string      `json:"status"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	StartTime time.Time   `json:"start_time"`
	EndTime   time.Time   `json:"end_time"`
}

type StepProgress struct {
	Status      string    `json:"status"`
	LastUpdated time.Time `json:"last_updated,omitempty"`
	Result      *Result   `json:"result,omitempty"`
	StartTime   time.Time `json:"start_time,omitempty"`
	EndTime     time.Time `json:"end_time,omitempty"`
}

type ServiceProgress struct {
	ServiceName    string                   `json:"service_name"`
	Type           string                   `json:"type"`
	Status         string                   `json:"status"`
	QueuePosition  int                      `json:"queue_position,omitempty"`
	TotalSteps     int                      `json:"total_steps"`
	CompletedSteps int                      `json:"completed_steps"`
	Steps          map[string]*StepProgress `json:"steps"`
	StartTime      time.Time                `json:"start_time"`
	LastUpdated    time.Time                `json:"last_updated"`
	mu             sync.RWMutex
}

// Cache Types
type StepCache struct {
	Result      *Result   `json:"result"`
	LastUpdated time.Time `json:"last_updated"`
}

type ServiceCache struct {
	ServiceName string               `json:"service_name"`
	Type        string               `json:"type"`
	Steps       map[string]StepCache `json:"steps"`
	LastUpdated time.Time            `json:"last_updated"`
	mu          sync.RWMutex         // Add this mutex for thread-safe access
}

// Queue Types
type QueuedCheck struct {
	ServiceName string    `json:"service_name"`
	Type        string    `json:"type"`
	StepsToRun  []string  `json:"steps_to_run"`
	QueueTime   time.Time `json:"queue_time"`
	Position    int       `json:"position"`
}

// System State Types
type SystemState struct {
	Status       string    `json:"status"`
	LastUpdated  time.Time `json:"last_updated"`
	PausedAt     time.Time `json:"paused_at,omitempty"`
	StepMode     bool      `json:"step_mode"`
	CurrentStep  string    `json:"current_step,omitempty"`
	PendingSteps []string  `json:"pending_steps,omitempty"`
	Message      string    `json:"message,omitempty"`
}

type ControlCommand string

const (
	CommandPause  ControlCommand = "pause"
	CommandResume ControlCommand = "resume"
	CommandStep   ControlCommand = "step"
	CommandReset  ControlCommand = "reset"
)

// Request/Response Types
type ServiceRequest struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	URL      string `json:"url"`
	Priority int    `json:"priority"`
}

type InvalidationRequest struct {
	ServiceName string   `json:"service_name"`
	Type        string   `json:"type"`
	Handlers    []string `json:"handlers,omitempty"`
	ResetTimes  bool     `json:"reset_times"`
}

// Debug Types
type SystemDebugInfo struct {
	Timestamp     time.Time        `json:"timestamp"`
	QueueStatus   QueueDebugInfo   `json:"queue_status"`
	CacheStatus   CacheDebugInfo   `json:"cache_status"`
	ProcessStatus ProcessDebugInfo `json:"process_status"`
}

type QueueDebugInfo struct {
	QueueLength  int          `json:"queue_length"`
	MaxQueueSize int          `json:"max_queue_size"`
	QueuedItems  []QueuedItem `json:"queued_items"`
}

type QueuedItem struct {
	ServiceName string    `json:"service_name"`
	Type        string    `json:"type"`
	Position    int       `json:"position"`
	QueueTime   time.Time `json:"queue_time"`
	WaitTime    string    `json:"wait_time"`
	StepsToRun  []string  `json:"steps_to_run"`
}

type CacheDebugInfo struct {
	TotalEntries int                   `json:"total_entries"`
	Entries      map[string]CacheEntry `json:"entries"`
}

type CacheEntry struct {
	ServiceName  string                `json:"service_name"`
	Type         string                `json:"type"`
	LastUpdated  time.Time             `json:"last_updated"`
	StepStatuses map[string]StepStatus `json:"step_statuses"`
}

type StepStatus struct {
	Status       string    `json:"status"`
	LastUpdated  time.Time `json:"last_updated"`
	CacheExpires time.Time `json:"cache_expires"`
	Age          string    `json:"age"`
}

type ProcessDebugInfo struct {
	ActiveProcesses int                       `json:"active_processes"`
	ProcessingItems map[string]ProcessingItem `json:"processing_items"`
}

type ProcessingItem struct {
	ServiceName    string    `json:"service_name"`
	Type           string    `json:"type"`
	StartTime      time.Time `json:"start_time"`
	ProcessTime    string    `json:"process_time"`
	CompletedSteps []string  `json:"completed_steps"`
	PendingSteps   []string  `json:"pending_steps"`
	TotalSteps     int       `json:"total_steps"`
}

// SystemMetrics represents system-wide metrics
type SystemMetrics struct {
	TotalCacheEntries int         `json:"total_cache_entries"`
	ActiveProcesses   int         `json:"active_processes"`
	QueueStats        QueueStats  `json:"queue_stats"`
	SystemState       SystemState `json:"system_state"`
}

type QueuedJobInfo struct {
	ServiceName   string    `json:"service_name"`
	Type          string    `json:"type"`
	QueuePosition int       `json:"queue_position"`
	QueueTime     time.Time `json:"queue_time"`
	WaitTime      string    `json:"wait_time"`
	StepsToRun    []string  `json:"steps_to_run"`
}

type HandlerConfig struct {
	Name         string   `json:"name"`
	CacheSeconds int      `json:"cacheSeconds"`
	Description  string   `json:"description"`
	Dependencies []string `json:"dependencies,omitempty"`
}

type QueueStats struct {
	QueueLength    int             `json:"queue_length"`
	MaxQueueSize   int             `json:"max_queue_size"`
	ActiveChecks   int             `json:"active_checks"`
	QueuedServices []string        `json:"queued_services"`
	QueuedJobs     []QueuedJobInfo `json:"queued_jobs"`
	QueueName      string          `json:"queue_name,omitempty"`
}

type CachedResponse struct {
	Status     string                `json:"status"`
	LastUpdate time.Time             `json:"last_update"`
	Steps      map[string]StepResult `json:"steps"`
}

type StepResult struct {
	Status     string    `json:"status"`
	LastUpdate time.Time `json:"last_update"`
	Result     *Result   `json:"result,omitempty"`
}

// ServiceResponse is an interface that can be implemented by different response types
type ServiceResponse interface {
	ResponseType() string
}

// Add QueueConfig if not already present
type QueueConfig struct {
	Name          string `json:"name"`
	Type          string `json:"type"`
	MaxConcurrent int    `json:"maxConcurrent"`
	QueueSize     int    `json:"queueSize"`
}

// Ensure CachedResponse implements ServiceResponse
func (cr *CachedResponse) ResponseType() string {
	return "cached"
}

// identifyExpiredSteps checks which steps need to be reprocessed
func (pm *ProcessManager) identifyExpiredSteps(serviceType string, cache *ServiceCache) []string {
	var expiredSteps []string

	// Find the service type configuration
	svcType, exists := pm.config.System.ServiceTypes[serviceType]
	if !exists {
		return expiredSteps
	}

	cache.mu.RLock()
	defer cache.mu.RUnlock()

	// Check each handler for expiration
	for handlerID, handler := range svcType.Handlers {
		stepCache, exists := cache.Steps[handlerID]
		if !exists || time.Since(stepCache.LastUpdated) > time.Duration(handler.CacheSeconds)*time.Second {
			expiredSteps = append(expiredSteps, handlerID)
		}
	}

	return expiredSteps
}

// queueRequest adds a request to the processing queue
func (pm *ProcessManager) queueRequest(req ServiceRequest, expiredSteps []string, cacheKey string) (ServiceResponse, error) {
	check := &QueuedCheck{
		ServiceName: req.Name,
		Type:        req.Type,
		StepsToRun:  expiredSteps,
		QueueTime:   time.Now(),
	}

	if err := pm.queueMgr.EnqueueJob(check); err != nil {
		return nil, err
	}

	return &QueueStatus{
		Status:     "queued",
		Position:   check.Position,
		CacheKey:   cacheKey,
		QueueTime:  check.QueueTime,
		StepsToRun: expiredSteps,
	}, nil
}

// getCachedResults is an internal method to get cached results for a service cache
func (pm *ProcessManager) getCachedResults(cache *ServiceCache) *CachedResponse {
	cache.mu.RLock()
	defer cache.mu.RUnlock()

	response := &CachedResponse{
		Status:     "cached",
		LastUpdate: cache.LastUpdated,
		Steps:      make(map[string]StepResult),
	}

	for stepID, stepCache := range cache.Steps {
		response.Steps[stepID] = StepResult{
			Status:     "cached",
			LastUpdate: stepCache.LastUpdated,
			Result:     stepCache.Result,
		}
	}

	return response
}

type JobHistory struct {
	CompletedJobs []*JobResult
	MaxJobs       int
	mu            sync.RWMutex
}

type JobResult struct {
	ServiceName string                   `json:"service_name"`
	Type        string                   `json:"type"`
	StartTime   time.Time                `json:"start_time"`
	EndTime     time.Time                `json:"end_time"`
	Duration    string                   `json:"duration"`
	Status      string                   `json:"status"`
	Steps       map[string]*StepProgress `json:"steps"`
	TreeOutput  string                   `json:"tree_output"`
}
