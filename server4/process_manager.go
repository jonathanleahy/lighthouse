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
	queueMgr    *QueueManager
	cache       map[string]*ServiceCache
	progress    map[string]*ServiceProgress
	state       *SystemState
	controlChan chan ControlCommand
	history     *JobHistory
	mu          sync.RWMutex
	stateMu     sync.RWMutex
}

func NewProcessManager(config TreeConfig) *ProcessManager {
	return &ProcessManager{
		config:   config,
		queueMgr: NewQueueManager(config.System.Queues),
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

// HandleRequest processes new service check requests
func (pm *ProcessManager) HandleRequest(req ServiceRequest) (ServiceResponse, error) {
	cacheKey := pm.getCacheKey(req.Name, req.Type)

	// Check if already processing
	if isProcessing := pm.queueMgr.IsProcessing(cacheKey); isProcessing {
		pos := pm.getQueuePosition(cacheKey)
		return &QueueStatus{
			Status:    "processing",
			Position:  pos,
			CacheKey:  cacheKey,
			StartTime: time.Now(),
		}, nil
	}

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
func (pm *ProcessManager) processNextInQueue() {
	check := pm.queueMgr.ProcessNextJob()
	if check == nil {
		return
	}

	// Process the check
	go func() {
		defer func() {
			queueName := pm.queueMgr.DetermineQueueName(check)
			pm.queueMgr.MarkJobCompleted(queueName, check.ServiceName)
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

	serviceType, exists := pm.config.System.ServiceTypes[progress.Type]
	if !exists {
		fmt.Printf("\n[%s] Error: service type %s not found\n",
			time.Now().Format("15:04:05"), progress.Type)
		return
	}

	// Process each step that needs to be run
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

		// Check handler configuration
		handlerConfig, exists := serviceType.Handlers[stepID]
		if !exists {
			fmt.Printf("\n[%s] No handler found for step: %s\n",
				time.Now().Format("15:04:05"), stepID)
			continue
		}

		// Check dependencies
		if !pm.areDependenciesMet(progress, handlerConfig.Dependencies) {
			fmt.Printf("\n[%s] Dependencies not met for step: %s\n",
				time.Now().Format("15:04:05"), stepID)
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

// Check if all dependencies for a step have been completed
func (pm *ProcessManager) areDependenciesMet(progress *ServiceProgress, dependencies []string) bool {
	if len(dependencies) == 0 {
		return true
	}

	progress.mu.RLock()
	defer progress.mu.RUnlock()

	for _, depStep := range dependencies {
		if step, exists := progress.Steps[depStep]; !exists || step.Status != "completed" {
			return false
		}
	}
	return true
}

func (pm *ProcessManager) generateTreeOutput(progress *ServiceProgress) string {
	if progress == nil {
		return ""
	}

	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("Service: %s (%s)\n", progress.ServiceName, progress.Type))
	builder.WriteString(fmt.Sprintf("├── Status: %s\n", progress.Status))
	builder.WriteString(fmt.Sprintf("├── Duration: %s\n", progress.LastUpdated.Sub(progress.StartTime)))
	builder.WriteString("└── Steps:\n")

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

// GetQueueStats returns current queue statistics
func (pm *ProcessManager) GetQueueStats() QueueStats {
	allStats := pm.queueMgr.GetQueueStats()

	// Aggregate stats from all queues
	totalStats := QueueStats{
		QueuedJobs:     make([]QueuedJobInfo, 0),
		QueuedServices: make([]string, 0),
	}

	for _, stats := range allStats {
		totalStats.QueueLength += stats.QueueLength
		totalStats.ActiveChecks += stats.ActiveChecks
		totalStats.QueuedServices = append(totalStats.QueuedServices, stats.QueuedServices...)
		totalStats.QueuedJobs = append(totalStats.QueuedJobs, stats.QueuedJobs...)
	}

	return totalStats
}

// GetSystemMetrics retrieves comprehensive system metrics
func (pm *ProcessManager) GetSystemMetrics() SystemMetrics {
	queueStats := pm.GetQueueStats()

	pm.stateMu.RLock()
	systemState := *pm.state
	pm.stateMu.RUnlock()

	pm.mu.RLock()
	totalCacheEntries := len(pm.cache)
	activeProcesses := len(pm.progress)
	pm.mu.RUnlock()

	return SystemMetrics{
		TotalCacheEntries: totalCacheEntries,
		ActiveProcesses:   activeProcesses,
		QueueStats:        queueStats,
		SystemState:       systemState,
	}
}

// Utility methods
// Continued from previous...

// Utility methods
func (pm *ProcessManager) getCacheKey(serviceName, processType string) string {
	return fmt.Sprintf("%s-%s", serviceName, processType)
}

func (pm *ProcessManager) getQueuePosition(cacheKey string) int {
	return pm.queueMgr.GetPosition(cacheKey)
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
	queueStats := pm.queueMgr.GetQueueStats()
	var totalLength, maxSize int

	for _, stats := range queueStats {
		totalLength += stats.QueueLength
		if stats.MaxQueueSize > maxSize {
			maxSize = stats.MaxQueueSize
		}

		for _, service := range stats.QueuedServices {
			parts := strings.Split(service, " (")
			if len(parts) == 2 {
				serviceName := parts[0]
				serviceType := strings.TrimRight(parts[1], ")")

				debug.QueueStatus.QueuedItems = append(debug.QueueStatus.QueuedItems, QueuedItem{
					ServiceName: serviceName,
					Type:        serviceType,
					Position:    len(debug.QueueStatus.QueuedItems) + 1,
					QueueTime:   time.Now(),
					WaitTime:    "N/A",
					StepsToRun:  []string{},
				})
			}
		}
	}

	debug.QueueStatus.QueueLength = totalLength
	debug.QueueStatus.MaxQueueSize = maxSize

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

// StartCleanupTask starts a periodic cleanup of old progress entries
func (pm *ProcessManager) StartCleanupTask(interval, maxAge time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		for range ticker.C {
			pm.cleanupOldProgress(maxAge)
		}
	}()
}

// cleanupOldProgress removes completed progress entries older than maxAge
func (pm *ProcessManager) cleanupOldProgress(maxAge time.Duration) {
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
