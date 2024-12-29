package jobscheduler

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/jonathanleahy/project/jobscheduler/internal/executor"
)

// Scheduler manages the job scheduling and processing
type Scheduler struct {
	config     Config
	executor   *executor.Executor
	channels   map[string]*Channel
	stats      map[string]*ChannelStats
	mu         sync.RWMutex
	processLog *os.File
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	startTime  time.Time
}

// Channel represents a processing channel
type Channel struct {
	Name      string
	Jobs      chan JobPayload
	Workers   int
	Timeout   time.Duration
	processor *Processor
}

// NewScheduler creates and returns a new Scheduler instance
func NewScheduler(cfg Config) (*Scheduler, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %v", err)
	}

	// Create working directory if it doesn't exist
	if err := os.MkdirAll(cfg.WorkDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create working directory: %v", err)
	}

	// Open process log
	processLog, err := os.OpenFile(cfg.ProcessingLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open process log: %v", err)
	}

	// Create executor
	exec, err := executor.NewExecutor(cfg.WorkDir)
	if err != nil {
		processLog.Close()
		return nil, fmt.Errorf("failed to create executor: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	s := &Scheduler{
		config:     cfg,
		executor:   exec,
		channels:   make(map[string]*Channel),
		stats:      make(map[string]*ChannelStats),
		processLog: processLog,
		ctx:        ctx,
		cancel:     cancel,
		startTime:  time.Now(),
	}

	return s, nil
}

// SubmitJob submits a new job for processing
func (s *Scheduler) SubmitJob(job JobPayload) error {
	if err := job.Validate(); err != nil {
		return fmt.Errorf("invalid job payload: %v", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Create or get channel
	channel, err := s.getOrCreateChannel(job)
	if err != nil {
		return err
	}

	// Initialize job status
	job.Status = JobStatusPending
	job.StartTime = time.Now()

	// Submit to channel
	select {
	case channel.Jobs <- job:
		// Update statistics
		s.updateStatsForNewJob(job.Channel)
		return nil
	default:
		return fmt.Errorf("channel %s is full", job.Channel)
	}
}

// getOrCreateChannel creates a new channel if it doesn't exist
func (s *Scheduler) getOrCreateChannel(job JobPayload) (*Channel, error) {
	channel, exists := s.channels[job.Channel]
	if !exists {
		// Create new channel
		workers := job.Workers
		if workers <= 0 {
			workers = s.config.DefaultWorkers
		}

		timeout := job.Timeout
		if timeout <= 0 {
			timeout = s.config.DefaultTimeout
		}

		channel = &Channel{
			Name:    job.Channel,
			Jobs:    make(chan JobPayload, s.config.ChannelBufferSize),
			Workers: workers,
			Timeout: timeout,
		}

		// Initialize channel processor
		processor := NewProcessor(ProcessorConfig{
			Channel:       channel,
			Executor:      s.executor,
			ProcessLog:    s.processLog,
			MaxOutputSize: s.config.MaxOutputSize,
		})

		channel.processor = processor
		s.channels[job.Channel] = channel
		s.stats[job.Channel] = &ChannelStats{}

		// Start the processor
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			processor.Start(s.ctx)
		}()
	}

	return channel, nil
}

// updateStatsForNewJob updates channel statistics for a new job
func (s *Scheduler) updateStatsForNewJob(channelName string) {
	stats := s.stats[channelName]
	stats.TotalJobs++
	stats.LastJobTime = time.Now()
}

// GetJobStatus retrieves the status of a specific job
func (s *Scheduler) GetJobStatus(jobID string) (*JobPayload, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Search for job in all channels
	for _, channel := range s.channels {
		jobs := channel.processor.GetActiveJobs()
		for _, job := range jobs {
			if job.ID == jobID {
				return &job, nil
			}
		}
	}

	return nil, fmt.Errorf("job not found: %s", jobID)
}

// ListJobs returns a list of jobs filtered by channel and status
func (s *Scheduler) ListJobs(channel, status string) ([]JobPayload, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var jobs []JobPayload

	// If channel is specified, only look in that channel
	if channel != "" {
		if ch, exists := s.channels[channel]; exists {
			activeJobs := ch.processor.GetActiveJobs()
			for _, job := range activeJobs {
				if status == "" || string(job.Status) == status {
					jobs = append(jobs, job)
				}
			}
		}
		return jobs, nil
	}

	// Otherwise, look in all channels
	for _, ch := range s.channels {
		activeJobs := ch.processor.GetActiveJobs()
		for _, job := range activeJobs {
			if status == "" || string(job.Status) == status {
				jobs = append(jobs, job)
			}
		}
	}

	return jobs, nil
}

// CancelJob cancels a specific job if it's still running
func (s *Scheduler) CancelJob(jobID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Search for job in all channels
	for _, channel := range s.channels {
		jobs := channel.processor.GetActiveJobs()
		for _, job := range jobs {
			if job.ID == jobID && job.Status == JobStatusRunning {
				// Signal cancellation through context
				if s.cancel != nil {
					s.cancel()
				}
				return nil
			}
		}
	}

	return fmt.Errorf("no active job found with ID: %s", jobID)
}

// GetStatsSummary returns summarized statistics for a time range
func (s *Scheduler) GetStatsSummary(from, to string) (*StatsSummary, error) {
	// Parse time range
	fromTime, err := time.Parse(time.RFC3339, from)
	if err != nil {
		return nil, fmt.Errorf("invalid from time: %v", err)
	}

	toTime, err := time.Parse(time.RFC3339, to)
	if err != nil {
		return nil, fmt.Errorf("invalid to time: %v", err)
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	summary := &StatsSummary{
		ActiveChannels: len(s.channels),
	}

	var totalDuration time.Duration
	var completedJobCount int64

	// Calculate statistics
	for _, channel := range s.channels {
		summary.TotalJobs += s.stats[channel.Name].TotalJobs
		summary.FailedJobs += s.stats[channel.Name].FailedJobs

		jobs := channel.processor.GetActiveJobs()
		for _, job := range jobs {
			if job.StartTime.After(fromTime) && job.StartTime.Before(toTime) {
				if job.Status == JobStatusComplete {
					summary.CompletedJobs++
					completedJobCount++
					if !job.EndTime.IsZero() {
						totalDuration += job.EndTime.Sub(job.StartTime)
					}
				}
			}
		}
	}

	// Calculate average runtime
	if completedJobCount > 0 {
		summary.AverageRuntime = float64(totalDuration) / float64(completedJobCount) / float64(time.Second)
	}

	return summary, nil
}

// GetOverallStats returns system-wide statistics
func (s *Scheduler) GetOverallStats() *OverallStats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := &OverallStats{
		ActiveChannels: len(s.channels),
		LastUpdate:     time.Now(),
		Uptime:         time.Since(s.startTime).String(),
	}

	// Calculate totals across all channels
	for _, channel := range s.channels {
		activeJobs := channel.processor.GetActiveJobs()
		stats.ActiveJobs += int64(len(activeJobs))

		if channelStats, exists := s.stats[channel.Name]; exists {
			stats.CompletedJobs += channelStats.TotalJobs - channelStats.FailedJobs
			stats.FailedJobs += channelStats.FailedJobs
			stats.QueuedJobs += int64(len(channel.Jobs))
		}
	}

	return stats
}

// GetChannelStats returns statistics for all channels
func (s *Scheduler) GetChannelStats() map[string]*ChannelStats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Create a copy of stats
	statsCopy := make(map[string]*ChannelStats)
	for name, stats := range s.stats {
		statsCopy[name] = &ChannelStats{
			Workers:     s.channels[name].Workers,
			ActiveJobs:  stats.ActiveJobs,
			TotalJobs:   stats.TotalJobs,
			FailedJobs:  stats.FailedJobs,
			LastJobTime: stats.LastJobTime,
		}
	}

	return statsCopy
}

// Shutdown gracefully shuts down the scheduler
func (s *Scheduler) Shutdown() error {
	log.Println("Starting graceful shutdown...")

	// Signal shutdown
	s.cancel()

	// Create shutdown context with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
	defer cancel()

	// Wait for all processors to complete with timeout
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("All processors completed successfully")
	case <-shutdownCtx.Done():
		log.Println("Shutdown timed out, some processors may still be running")
	}

	// Cleanup executor
	s.executor.Cleanup()

	// Close process log
	if err := s.processLog.Close(); err != nil {
		return fmt.Errorf("failed to close process log: %v", err)
	}

	return nil
}
