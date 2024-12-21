package main

import (
	"fmt"
	"math/rand"
	"time"
)

type HandlerFunc func(ctx *Context) *Result

// HandlerRegistry contains all available handlers
var HandlerRegistry = map[string]HandlerFunc{
	// Check workflow handlers
	"dnsHandler": func(ctx *Context) *Result {
		return simulateWork("DNS Check", ctx, 1, 5)
	},
	"githubHandler": func(ctx *Context) *Result {
		return simulateWork("GitHub Check", ctx, 1, 5)
	},
	"performanceHandler": func(ctx *Context) *Result {
		return simulateWork("Performance Analysis", ctx, 2, 5)
	},
	"aiHandler": func(ctx *Context) *Result {
		return simulateWork("AI Analysis", ctx, 3, 5)
	},

	// Report workflow handlers
	"metricCollector": func(ctx *Context) *Result {
		return simulateWork("Metric Collection", ctx, 1, 4)
	},
	"logAnalyzer": func(ctx *Context) *Result {
		return simulateWork("Log Analysis", ctx, 2, 5)
	},
	"trendAnalyzer": func(ctx *Context) *Result {
		return simulateWork("Trend Analysis", ctx, 2, 4)
	},
	"reportGenerator": func(ctx *Context) *Result {
		return simulateWork("Report Generation", ctx, 3, 5)
	},
}

// simulateWork simulates processing with random delay and console output
func simulateWork(handlerName string, ctx *Context, minSec, maxSec int) *Result {
	start := time.Now()

	fmt.Printf("\n[%s] Starting %s for service '%s' (%s workflow)\n",
		start.Format("15:04:05"),
		handlerName,
		ctx.ServiceName,
		ctx.ProcessType)

	// Random delay between minSec and maxSec seconds
	delay := time.Duration(minSec+rand.Intn(maxSec-minSec+1)) * time.Second
	time.Sleep(delay)

	// Simulate occasional failures (10% chance)
	status := "completed"
	if rand.Float32() < 0.1 {
		status = "failed"
	}

	end := time.Now()

	fmt.Printf("[%s] %s %s for service '%s' after %v\n",
		end.Format("15:04:05"),
		status,
		handlerName,
		ctx.ServiceName,
		delay)

	return &Result{
		Status:    status,
		Message:   fmt.Sprintf("%s %s in %v", handlerName, status, delay),
		StartTime: start,
		EndTime:   end,
	}
}
