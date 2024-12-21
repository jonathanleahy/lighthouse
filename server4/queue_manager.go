package main

import (
	"fmt"
	"strings"
	"time"
)

// NewQueueManager creates a new queue manager with the specified configuration
func NewQueueManager(configs []QueueConfig) *QueueManager {
	qm := &QueueManager{
		queues: make(map[string]*WorkQueue),
		config: make(map[string]QueueConfig),
	}

	// Initialize each queue
	for _, cfg := range configs {
		qm.queues[cfg.Name] = &WorkQueue{
			name:          cfg.Name,
			items:         make([]*QueuedCheck, 0),
			processing:    make(map[string]bool),
			maxSize:       cfg.QueueSize,
			maxConcurrent: cfg.MaxConcurrent,
		}
		qm.config[cfg.Name] = cfg
	}

	return qm
}

// IsProcessing checks if a job is currently being processed
func (qm *QueueManager) IsProcessing(cacheKey string) bool {
	qm.mu.RLock()
	defer qm.mu.RUnlock()

	for _, queue := range qm.queues {
		queue.mu.RLock()
		if queue.processing[cacheKey] {
			queue.mu.RUnlock()
			return true
		}
		queue.mu.RUnlock()
	}
	return false
}

// GetPosition returns the position of a job in its queue
func (qm *QueueManager) GetPosition(cacheKey string) int {
	qm.mu.RLock()
	defer qm.mu.RUnlock()

	for _, queue := range qm.queues {
		queue.mu.RLock()
		for i, check := range queue.items {
			if fmt.Sprintf("%s-%s", check.ServiceName, check.Type) == cacheKey {
				queue.mu.RUnlock()
				return i + 1
			}
		}
		queue.mu.RUnlock()
	}
	return 0
}

// DetermineQueueName returns the appropriate queue for a check
func (qm *QueueManager) DetermineQueueName(check *QueuedCheck) string {
	if check.Type == "check" {
		if containsAny(check.StepsToRun, "ai") {
			return "ai_analysis"
		}
		if containsAny(check.StepsToRun, "performance") {
			return "performance_analysis"
		}
	}
	return "service_checks"
}

// EnqueueJob places a job in the appropriate queue based on handler type
func (qm *QueueManager) EnqueueJob(check *QueuedCheck) error {
	targetQueue := qm.DetermineQueueName(check)

	queue, exists := qm.queues[targetQueue]
	if !exists {
		return fmt.Errorf("queue %s not found", targetQueue)
	}

	queue.mu.Lock()
	defer queue.mu.Unlock()

	// Check queue capacity
	if len(queue.items) >= queue.maxSize {
		return fmt.Errorf("queue %s is full (capacity: %d)", targetQueue, queue.maxSize)
	}

	// Add to queue
	check.QueueTime = time.Now()
	check.Position = len(queue.items) + 1
	queue.items = append(queue.items, check)

	fmt.Printf("\n[%s] Job for service '%s' added to queue '%s' at position %d\n",
		time.Now().Format("15:04:05"),
		check.ServiceName,
		targetQueue,
		check.Position)

	return nil
}

// ProcessNextJob processes the next available job from appropriate queue
func (qm *QueueManager) ProcessNextJob() *QueuedCheck {
	// Try queues in priority order
	queuePriority := []string{"ai_analysis", "performance_analysis", "service_checks"}

	for _, queueName := range queuePriority {
		queue, exists := qm.queues[queueName]
		if !exists {
			continue
		}

		queue.mu.Lock()

		// Skip if queue is empty or at max concurrent capacity
		if len(queue.items) == 0 || len(queue.processing) >= queue.maxConcurrent {
			queue.mu.Unlock()
			continue
		}

		// Get next job
		job := queue.items[0]
		queue.items = queue.items[1:]

		// Mark as processing
		queue.processing[job.ServiceName] = true

		// Update positions for remaining items
		for i, item := range queue.items {
			item.Position = i + 1
		}

		queue.mu.Unlock()

		fmt.Printf("\n[%s] Processing job from queue '%s' for service '%s'\n",
			time.Now().Format("15:04:05"),
			queueName,
			job.ServiceName)

		return job
	}

	return nil
}

// MarkJobCompleted removes a job from the processing map
func (qm *QueueManager) MarkJobCompleted(queueName, serviceName string) {
	if queue, exists := qm.queues[queueName]; exists {
		queue.mu.Lock()
		delete(queue.processing, serviceName)
		queue.mu.Unlock()

		fmt.Printf("\n[%s] Completed job for service '%s' in queue '%s'\n",
			time.Now().Format("15:04:05"),
			serviceName,
			queueName)
	}
}

// GetQueueStats returns statistics for all queues
func (qm *QueueManager) GetQueueStats() map[string]QueueStats {
	qm.mu.RLock()
	defer qm.mu.RUnlock()

	stats := make(map[string]QueueStats)

	for name, queue := range qm.queues {
		queue.mu.RLock()
		queueStats := QueueStats{
			QueueLength:    len(queue.items),
			MaxQueueSize:   queue.maxSize,
			ActiveChecks:   len(queue.processing),
			QueuedServices: make([]string, 0),
			QueuedJobs:     make([]QueuedJobInfo, 0),
			QueueName:      name,
		}

		// Add details for each queued item
		for _, item := range queue.items {
			queueStats.QueuedServices = append(
				queueStats.QueuedServices,
				fmt.Sprintf("%s (%s)", item.ServiceName, item.Type),
			)

			queueStats.QueuedJobs = append(queueStats.QueuedJobs, QueuedJobInfo{
				ServiceName:   item.ServiceName,
				Type:          item.Type,
				QueuePosition: item.Position,
				QueueTime:     item.QueueTime,
				WaitTime:      time.Since(item.QueueTime).String(),
				StepsToRun:    item.StepsToRun,
			})
		}
		queue.mu.RUnlock()

		stats[name] = queueStats
	}

	return stats
}

// GetActiveJobCount returns the total number of jobs currently being processed
func (qm *QueueManager) GetActiveJobCount() int {
	qm.mu.RLock()
	defer qm.mu.RUnlock()

	total := 0
	for _, queue := range qm.queues {
		queue.mu.RLock()
		total += len(queue.processing)
		queue.mu.RUnlock()
	}
	return total
}

func (pm *ProcessManager) InvalidateCache(req InvalidationRequest) error {
	if req.ServiceName == "*" {
		pm.mu.Lock()
		for cacheKey, cache := range pm.cache {
			if req.Type == "*" || strings.HasSuffix(cacheKey, "-"+req.Type) {
				cache.mu.Lock()
				if len(req.Handlers) > 0 {
					// Invalidate specific handlers
					for _, handlerID := range req.Handlers {
						delete(cache.Steps, handlerID)
					}
				} else {
					// Invalidate all handlers
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

	// Handle single service invalidation
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
		// Invalidate specific handlers
		for _, handlerID := range req.Handlers {
			delete(cache.Steps, handlerID)
		}
	} else {
		// Invalidate all handlers
		cache.Steps = make(map[string]StepCache)
	}

	if req.ResetTimes {
		cache.LastUpdated = time.Time{}
	}

	return nil
}
