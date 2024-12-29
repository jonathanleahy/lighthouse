package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

// ProcessManager handles all service processing
type ProcessManager struct {
	config      TreeConfig
	queue       *WorkQueue
	cache       map[string]*ServiceCache
	progress    map[string]*ServiceProgress
	state       *SystemState
	controlChan chan ControlCommand
	history     *JobHistory
	mu          sync.RWMutex
	stateMu     sync.RWMutex
}

// WorkQueue manages the processing queue
type WorkQueue struct {
	items      []*QueuedCheck
	processing map[string]bool
	maxSize    int
	mu         sync.RWMutex
}

// JobHistory tracks completed jobs
type JobHistory struct {
	CompletedJobs []*JobResult
	MaxJobs       int
	mu            sync.RWMutex
}

// JobResult stores the complete result of a job
type JobResult struct {
	ServiceName string
	Type        string
	StartTime   time.Time
	EndTime     time.Time
	Duration    string
	Status      string
	Steps       map[string]*StepProgress
	TreeOutput  string
}

func NewProcessManager(config TreeConfig) *ProcessManager {
	return &ProcessManager{
		config:   config,
		queue:    NewWorkQueue(100),
		cache:    make(map[string]*ServiceCache),
		progress: make(map[string]*ServiceProgress),
		state: &SystemState{
			Status:      "running",
			LastUpdated: time.Now(),
			StepMode:    false,
		},
		controlChan: make(chan ControlCommand, 10),
		history: &JobHistory{
			CompletedJobs: make([]*JobResult, 0),
			MaxJobs:       1000,
		},
	}
}

func NewWorkQueue(maxSize int) *WorkQueue {
	return &WorkQueue{
		items:      make([]*QueuedCheck, 0),
		processing: make(map[string]bool),
		maxSize:    maxSize,
	}
}

// HandleRequest processes new service check requests
func (pm *ProcessManager) HandleRequest(req ServiceRequest) (ServiceResponse, error) {
	cacheKey := pm.getCacheKey(req.Name, req.Type)

	// Check if already processing
	pm.queue.mu.RLock()
	if pm.queue.processing[cacheKey] {
		pos := pm.getQueuePosition(cacheKey)
		pm.queue.mu.RUnlock()
		return &QueueStatus{
			Status:    "processing",
			Position:  pos,
			CacheKey:  cacheKey,
			StartTime: time.Now(),
		}, nil
	}
	pm.queue.mu.RUnlock()

	// Initialize or get cache entry
	pm.mu.Lock()
	serviceCache, exists := pm.cache[cacheKey]
	if !exists {
		serviceCache = &ServiceCache{
			ServiceName: req.Name,
			Type:        req.Type,
			Steps:       make(map[string]StepCache),
			LastUpdated: time.Now(),
		}
		pm.cache[cacheKey] = serviceCache
	}
	pm.mu.Unlock()

	// Check which steps need processing
	expiredSteps := pm.identifyExpiredSteps(req.Type, serviceCache)

	// If no steps need processing, return cached results
	if len(expiredSteps) == 0 {
		return pm.getCachedResults(serviceCache), nil
	}

	// Add to queue if steps need processing
	return pm.queueRequest(req, expiredSteps, cacheKey)
}

// processNextInQueue processes the next item in the queue
// processNextInQueue processes the next item in the queue
func (pm *ProcessManager) processNextInQueue() {
	pm.queue.mu.Lock()

	// Log total queue items before processing
	fmt.Printf("\n[%s] Queue Processing - Total Items: %d\n",
		time.Now().Format("15:04:05"),
		len(pm.queue.items))

	// Check if queue is empty or processing limit is reached
	if len(pm.queue.items) == 0 {
		pm.queue.mu.Unlock()
		return
	}

	// Determine max concurrent jobs from configuration
	var maxConcurrent int
	for _, queueConfig := range pm.config.System.Queues {
		if queueConfig.Name == "service_checks" {
			maxConcurrent = queueConfig.MaxConcurrent
			break
		}
	}

	// Log current processing state
	fmt.Printf("Current Processing Status:\n")
	fmt.Printf("Max Concurrent Jobs: %d\n", maxConcurrent)
	fmt.Printf("Active Checks: %d\n", len(pm.queue.processing))

	// If processing limit is reached, do not process more jobs
	if len(pm.queue.processing) >= maxConcurrent {
		fmt.Printf("Processing limit reached. Skipping job processing.\n")
		pm.queue.mu.Unlock()
		return
	}

	// Get next check from queue
	check := pm.queue.items[0]
	pm.queue.items = pm.queue.items[1:]

	// Mark as processing
	cacheKey := pm.getCacheKey(check.ServiceName, check.Type)
	pm.queue.processing[cacheKey] = true

	// Update positions for remaining items
	for i, item := range pm.queue.items {
		item.Position = i + 1
	}
	pm.queue.mu.Unlock()

	// Log job being processed
	fmt.Printf("[%s] Processing Job: %s (%s)\n",
		time.Now().Format("15:04:05"),
		check.ServiceName,
		check.Type)

	// Process the check
	go func() {
		defer func() {
			pm.queue.mu.Lock()
			delete(pm.queue.processing, cacheKey)
			pm.queue.mu.Unlock()
		}()

		pm.processCheck(check)
	}()
}
func (pm *ProcessManager) processCheck(check *QueuedCheck) {
	if check == nil {
		fmt.Printf("\n[%s] Error: nil check received\n", time.Now().Format("15:04:05"))
		return
	}

	progress := pm.getOrCreateProgress(ServiceRequest{
		Name: check.ServiceName,
		Type: check.Type,
	})

	if progress == nil {
		fmt.Printf("\n[%s] Error: could not create progress tracker for %s\n",
			time.Now().Format("15:04:05"), check.ServiceName)
		return
	}

	progress.mu.Lock()
	progress.Status = "processing"
	progress.StartTime = time.Now()
	progress.mu.Unlock()

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("\n[%s] Recovered from panic in processCheck: %v\n",
				time.Now().Format("15:04:05"), r)

			progress.mu.Lock()
			progress.Status = "failed"
			progress.LastUpdated = time.Now()
			progress.mu.Unlock()
		}
	}()

	// Identify steps to run with dependencies
	serviceType, exists := pm.config.System.ServiceTypes[progress.Type]
	if !exists {
		fmt.Printf("\n[%s] Error: service type %s not found\n",
			time.Now().Format("15:04:05"), progress.Type)
		return
	}

	// Determine which steps can be run based on dependencies
	for _, stepID := range check.StepsToRun {
		pm.stateMu.RLock()
		currentState := pm.state.Status
		stepMode := pm.state.StepMode
		pm.stateMu.RUnlock()

		if currentState == "paused" && !stepMode {
			fmt.Printf("\n[%s] Processing paused before step: %s\n",
				time.Now().Format("15:04:05"), stepID)
			return
		}

		// Check handler configuration and dependencies
		handlerConfig, exists := serviceType.Handlers[stepID]
		if !exists {
			fmt.Printf("\n[%s] No handler found for step: %s\n",
				time.Now().Format("15:04:05"), stepID)
			progress.mu.Lock()
			if progress.Steps[stepID] != nil {
				progress.Steps[stepID].Status = "failed"
				progress.Steps[stepID].Result = &Result{
					Status:  "error",
					Message: fmt.Sprintf("No handler found for step %s", stepID),
				}
			}
			progress.mu.Unlock()
			continue
		}

		// Check if all dependencies are completed successfully
		dependenciesMet := true
		for _, depStep := range handlerConfig.Dependencies {
			progress.mu.RLock()
			depStepProgress, exists := progress.Steps[depStep]
			progress.mu.RUnlock()

			if !exists || depStepProgress.Status != "completed" {
				dependenciesMet = false
				fmt.Printf("\n[%s] Dependencies not met for step: %s\n",
					time.Now().Format("15:04:05"), stepID)
				break
			}
		}

		// Skip if dependencies are not met
		if !dependenciesMet {
			continue
		}

		// Update current step in system state
		pm.stateMu.Lock()
		pm.state.CurrentStep = stepID
		pm.state.LastUpdated = time.Now()
		pm.stateMu.Unlock()

		// Execute the handler
		ctx := &Context{
			ServiceName: progress.ServiceName,
			ProcessType: progress.Type,
			StepID:      stepID,
		}

		result := HandlerRegistry[handlerConfig.Name](ctx)

		// Update progress
		progress.mu.Lock()
		if progress.Steps[stepID] != nil {
			progress.Steps[stepID].Status = "completed"
			progress.Steps[stepID].Result = result
			progress.Steps[stepID].EndTime = time.Now()
			progress.Steps[stepID].LastUpdated = time.Now()
		}
		progress.CompletedSteps++
		progress.mu.Unlock()

		// Handle step mode
		if stepMode {
			select {
			case cmd := <-pm.controlChan:
				if cmd != CommandStep {
					return
				}
			case <-time.After(time.Second * 30):
				fmt.Printf("\n[%s] Step timeout for: %s\n",
					time.Now().Format("15:04:05"), stepID)
				return
			}
		}
	}

	// Update final status
	progress.mu.Lock()
	progress.Status = "completed"
	progress.LastUpdated = time.Now()
	progress.mu.Unlock()

	// Generate tree output and save history
	treeOutput := pm.generateTreeOutput(progress)
	if treeOutput != "" {
		result := &JobResult{
			ServiceName: progress.ServiceName,
			Type:        progress.Type,
			StartTime:   progress.StartTime,
			EndTime:     progress.LastUpdated,
			Duration:    progress.LastUpdated.Sub(progress.StartTime).String(),
			Status:      progress.Status,
			Steps:       progress.Steps,
			TreeOutput:  treeOutput,
		}

		// Save to history
		pm.history.mu.Lock()
		pm.history.CompletedJobs = append(pm.history.CompletedJobs, result)
		if len(pm.history.CompletedJobs) > pm.history.MaxJobs {
			pm.history.CompletedJobs = pm.history.CompletedJobs[1:]
		}
		pm.history.mu.Unlock()

		// Print tree output
		fmt.Printf("\n%s\n", treeOutput)
	}

	// Update cache
	cacheKey := pm.getCacheKey(progress.ServiceName, progress.Type)
	pm.mu.Lock()
	if serviceCache, exists := pm.cache[cacheKey]; exists {
		serviceCache.mu.Lock()
		for stepID, step := range progress.Steps {
			if step.Result != nil {
				serviceCache.Steps[stepID] = StepCache{
					Result:      step.Result,
					LastUpdated: time.Now(),
				}
			}
		}
		serviceCache.LastUpdated = time.Now()
		serviceCache.mu.Unlock()
	}
	pm.mu.Unlock()
}

// generateTreeOutput creates a tree visualization of the job
func (pm *ProcessManager) generateTreeOutput(progress *ServiceProgress) string {
	if progress == nil {
		return ""
	}

	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("Service: %s (%s)\n", progress.ServiceName, progress.Type))
	builder.WriteString(fmt.Sprintf("├── Status: %s\n", progress.Status))
	builder.WriteString(fmt.Sprintf("├── Duration: %s\n", progress.LastUpdated.Sub(progress.StartTime)))
	builder.WriteString("└── Steps:\n")

	// Sort steps for consistent output
	stepIDs := make([]string, 0, len(progress.Steps))
	for stepID := range progress.Steps {
		stepIDs = append(stepIDs, stepID)
	}
	sort.Strings(stepIDs)

	for i, stepID := range stepIDs {
		step := progress.Steps[stepID]
		prefix := "    ├──"
		if i == len(stepIDs)-1 {
			prefix = "    └──"
		}

		builder.WriteString(fmt.Sprintf("%s %s\n", prefix, stepID))
		childPrefix := "    │   "
		if i == len(stepIDs)-1 {
			childPrefix = "        "
		}

		builder.WriteString(fmt.Sprintf("%s├── Status: %s\n", childPrefix, step.Status))
		if step.StartTime.IsZero() {
			builder.WriteString(fmt.Sprintf("%s└── Duration: Not started\n", childPrefix))
		} else if step.EndTime.IsZero() {
			builder.WriteString(fmt.Sprintf("%s└── Duration: In progress\n", childPrefix))
		} else {
			builder.WriteString(fmt.Sprintf("%s└── Duration: %s\n", childPrefix, step.EndTime.Sub(step.StartTime)))
		}
	}

	return builder.String()
}

// HandleControlCommand processes system control commands
func (pm *ProcessManager) HandleControlCommand(cmd ControlCommand) (*SystemState, error) {
	pm.stateMu.Lock()
	defer pm.stateMu.Unlock()

	switch cmd {
	case CommandPause:
		if pm.state.Status == "paused" {
			return pm.state, fmt.Errorf("system already paused")
		}
		pm.state.Status = "paused"
		pm.state.PausedAt = time.Now()
		pm.state.Message = "System paused"

	case CommandResume:
		if pm.state.Status == "running" {
			return pm.state, fmt.Errorf("system already running")
		}
		pm.state.Status = "running"
		pm.state.StepMode = false
		pm.state.Message = "System resumed"

	case CommandStep:
		if pm.state.Status != "paused" {
			return pm.state, fmt.Errorf("system must be paused to step")
		}
		pm.state.StepMode = true
		pm.state.Status = "stepping"
		pm.state.Message = "Executing one step"

		go func() {
			pm.controlChan <- CommandStep
		}()

	case CommandReset:
		pm.state.Status = "running"
		pm.state.StepMode = false
		pm.state.CurrentStep = ""
		pm.state.PendingSteps = nil
		pm.state.Message = "System reset"
	}

	pm.state.LastUpdated = time.Now()
	return pm.state, nil
}

// InvalidateCache handles cache invalidation requests
func (pm *ProcessManager) InvalidateCache(req InvalidationRequest) error {
	if req.ServiceName == "*" {
		pm.mu.Lock()
		for cacheKey, cache := range pm.cache {
			if req.Type == "*" || strings.HasSuffix(cacheKey, "-"+req.Type) {
				cache.mu.Lock()
				if len(req.Handlers) > 0 {
					for _, handlerID := range req.Handlers {
						delete(cache.Steps, handlerID)
					}
				} else {
					cache.Steps = make(map[string]StepCache)
				}
				if req.ResetTimes {
					cache.LastUpdated = time.Time{}
				}
				cache.mu.Unlock()
			}
		}
		pm.mu.Unlock()
		return nil
	}

	cacheKey := pm.getCacheKey(req.ServiceName, req.Type)
	pm.mu.Lock()
	cache, exists := pm.cache[cacheKey]
	pm.mu.Unlock()

	if !exists {
		return fmt.Errorf("no cache found for service %s (%s)", req.ServiceName, req.Type)
	}

	cache.mu.Lock()
	defer cache.mu.Unlock()

	if len(req.Handlers) > 0 {
		for _, handlerID := range req.Handlers {
			delete(cache.Steps, handlerID)
		}
	} else {
		cache.Steps = make(map[string]StepCache)
	}

	if req.ResetTimes {
		cache.LastUpdated = time.Time{}
	}

	return nil
}

// GetSystemDebugInfo returns comprehensive system debug information
func (pm *ProcessManager) GetSystemDebugInfo() SystemDebugInfo {
	debug := SystemDebugInfo{
		Timestamp: time.Now(),
		QueueStatus: QueueDebugInfo{
			QueuedItems: make([]QueuedItem, 0),
		},
		CacheStatus: CacheDebugInfo{
			Entries: make(map[string]CacheEntry),
		},
		ProcessStatus: ProcessDebugInfo{
			ProcessingItems: make(map[string]ProcessingItem),
		},
	}

	// Queue Status
	pm.queue.mu.RLock()
	debug.QueueStatus.QueueLength = len(pm.queue.items)
	debug.QueueStatus.MaxQueueSize = pm.queue.maxSize

	for _, item := range pm.queue.items {
		debug.QueueStatus.QueuedItems = append(debug.QueueStatus.QueuedItems, QueuedItem{
			ServiceName: item.ServiceName,
			Type:        item.Type,
			Position:    item.Position,
			QueueTime:   item.QueueTime,
			WaitTime:    time.Since(item.QueueTime).String(),
			StepsToRun:  item.StepsToRun,
		})
	}
	pm.queue.mu.RUnlock()

	// Cache Status
	pm.mu.RLock()
	debug.CacheStatus.TotalEntries = len(pm.cache)
	for key, cache := range pm.cache {
		cache.mu.RLock()
		entry := CacheEntry{
			ServiceName:  cache.ServiceName,
			Type:         cache.Type,
			LastUpdated:  cache.LastUpdated,
			StepStatuses: make(map[string]StepStatus),
		}

		if serviceType, exists := pm.config.System.ServiceTypes[cache.Type]; exists {
			for stepID, stepCache := range cache.Steps {
				if handlerConfig, exists := serviceType.Handlers[stepID]; exists {
					cacheExpiry := stepCache.LastUpdated.Add(time.Duration(handlerConfig.CacheSeconds) * time.Second)
					status := "valid"
					if time.Now().After(cacheExpiry) {
						status = "expired"
					}

					entry.StepStatuses[stepID] = StepStatus{
						Status:       status,
						LastUpdated:  stepCache.LastUpdated,
						CacheExpires: cacheExpiry,
						Age:          time.Since(stepCache.LastUpdated).String(),
					}
				}
			}
		}
		cache.mu.RUnlock()
		debug.CacheStatus.Entries[key] = entry
	}
	pm.mu.RUnlock()

	// Process Status
	pm.mu.RLock()
	debug.ProcessStatus.ActiveProcesses = len(pm.progress)
	for key, progress := range pm.progress {
		progress.mu.RLock()
		if progress.Status == "processing" || progress.Status == "initializing" {
			var completedSteps, pendingSteps []string
			for stepID, step := range progress.Steps {
				if step.Status == "completed" {
					completedSteps = append(completedSteps, stepID)
				} else if step.Status == "pending" || step.Status == "processing" {
					pendingSteps = append(pendingSteps, stepID)
				}
			}

			debug.ProcessStatus.ProcessingItems[key] = ProcessingItem{
				ServiceName:    progress.ServiceName,
				Type:           progress.Type,
				StartTime:      progress.StartTime,
				ProcessTime:    time.Since(progress.StartTime).String(),
				CompletedSteps: completedSteps,
				PendingSteps:   pendingSteps,
				TotalSteps:     progress.TotalSteps,
			}
		}
		progress.mu.RUnlock()
	}
	pm.mu.RUnlock()

	return debug
}

// GetQueueStatus returns current queue statistics
func (pm *ProcessManager) GetQueueStats() QueueStats {
	pm.queue.mu.RLock()
	defer pm.queue.mu.RUnlock()

	// Log detailed queue information
	fmt.Printf("\n[%s] Queue Statistics:\n", time.Now().Format("15:04:05"))
	fmt.Printf("Total Queue Items: %d\n", len(pm.queue.items))
	fmt.Printf("Active Processing: %d\n", len(pm.queue.processing))

	stats := QueueStats{
		QueueLength:    len(pm.queue.items),
		MaxQueueSize:   pm.queue.maxSize,
		ActiveChecks:   len(pm.queue.processing),
		QueuedServices: make([]string, 0),
		QueuedJobs:     make([]QueuedJobInfo, 0),
	}

	// Log each queued item
	for i, item := range pm.queue.items {
		queuedService := fmt.Sprintf("%s (%s)", item.ServiceName, item.Type)
		stats.QueuedServices = append(stats.QueuedServices, queuedService)

		queuedJob := QueuedJobInfo{
			ServiceName:   item.ServiceName,
			Type:          item.Type,
			QueuePosition: i + 1,
			QueueTime:     item.QueueTime,
			WaitTime:      time.Since(item.QueueTime).String(),
			StepsToRun:    item.StepsToRun,
		}
		stats.QueuedJobs = append(stats.QueuedJobs, queuedJob)

		fmt.Printf("Queued Job %d: %s (%s), Position: %d, Queued At: %s\n",
			i+1,
			item.ServiceName,
			item.Type,
			i+1,
			item.QueueTime.Format(time.RFC3339))
	}

	return stats
}

// Utility methods
func (pm *ProcessManager) getCacheKey(serviceName, processType string) string {
	return fmt.Sprintf("%s-%s", serviceName, processType)
}

func (pm *ProcessManager) getQueuePosition(cacheKey string) int {
	pm.queue.mu.RLock()
	defer pm.queue.mu.RUnlock()

	for i, check := range pm.queue.items {
		if pm.getCacheKey(check.ServiceName, check.Type) == cacheKey {
			return i + 1
		}
	}
	return 0
}

func (pm *ProcessManager) getHandlerConfig(serviceType, handlerID string) *HandlerConfig {
	if svcType, exists := pm.config.System.ServiceTypes[serviceType]; exists {
		if handler, exists := svcType.Handlers[handlerID]; exists {
			return &handler
		}
	}
	return nil
}

func (pm *ProcessManager) getOrCreateProgress(req ServiceRequest) *ServiceProgress {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	cacheKey := pm.getCacheKey(req.Name, req.Type)
	if progress, exists := pm.progress[cacheKey]; exists {
		return progress
	}

	progress := &ServiceProgress{
		ServiceName: req.Name,
		Type:        req.Type,
		Status:      "initializing",
		Steps:       make(map[string]*StepProgress),
		StartTime:   time.Now(),
		LastUpdated: time.Now(),
	}

	// Initialize steps based on service type
	if serviceType, exists := pm.config.System.ServiceTypes[req.Type]; exists {
		progress.TotalSteps = len(serviceType.Handlers)
		for stepID := range serviceType.Handlers {
			progress.Steps[stepID] = &StepProgress{
				Status: "pending",
			}
		}
	}

	pm.progress[cacheKey] = progress
	return progress
}

// StartCleanupTask starts a periodic cleanup of old progress entries
func (pm *ProcessManager) StartCleanupTask(interval, maxAge time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		for range ticker.C {
			pm.CleanupOldProgress(maxAge)
		}
	}()
}

// CleanupOldProgress removes completed progress entries older than the specified duration
func (pm *ProcessManager) CleanupOldProgress(maxAge time.Duration) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	now := time.Now()
	for key, progress := range pm.progress {
		progress.mu.RLock()
		if progress.Status == "completed" && now.Sub(progress.LastUpdated) > maxAge {
			delete(pm.progress, key)
		}
		progress.mu.RUnlock()
	}
}
