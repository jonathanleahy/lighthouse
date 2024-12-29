package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type Command struct {
	fs      *flag.FlagSet
	project string
	output  string
	debug   bool
}

func main() {
	listCmd := &Command{
		fs: flag.NewFlagSet("list", flag.ExitOnError),
	}

	listCmd.fs.StringVar(&listCmd.project, "project", "", "Project name")
	listCmd.fs.StringVar(&listCmd.output, "output", "", "Output file path")
	listCmd.fs.BoolVar(&listCmd.debug, "debug", false, "Enable debug output")

	if len(os.Args) < 2 {
		fmt.Println("Usage: codefresh [list] [flags]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "list":
		listCmd.fs.Parse(os.Args[2:])
		err := listCmd.Run()
		if err != nil {
			writeErrorResponse(err.Error())
			os.Exit(1)
		}
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func (c *Command) Run() error {
	if _, err := exec.LookPath("codefresh"); err != nil {
		return fmt.Errorf("Codefresh CLI is not installed")
	}

	pipelineJson, err := c.getPipelineInfo()
	if err != nil {
		return fmt.Errorf("Failed to retrieve pipeline info: %v", err)
	}

	if c.debug {
		fmt.Printf("Pipeline JSON:\n%s\n", pipelineJson)
	}

	pipelineResult := gjson.Parse(pipelineJson)

	fmt.Println("\nPipeline Execution Status:")
	fmt.Println("==========================")

	var pipelineArray []gjson.Result
	if pipelineResult.IsArray() {
		pipelineArray = pipelineResult.Array()
	} else {
		pipelineArray = []gjson.Result{pipelineResult}
	}

	for _, p := range pipelineArray {
		metadata := p.Get("metadata")
		if !metadata.Exists() {
			continue
		}

		name := metadata.Get("name").String()
		fmt.Printf("\n📋 Pipeline: %s\n", name)

		buildsJson, err := c.getPipelineBuilds(name)
		if err != nil {
			fmt.Printf("   Started: Unable to fetch\n")
			fmt.Printf("   Duration: Unable to fetch\n")
			fmt.Printf("   Status:  ❌ Error: %v\n", err)
			continue
		}

		if c.debug && buildsJson != "" {
			fmt.Printf("Build JSON for %s:\n%s\n", name, buildsJson)
		}

		if buildsJson == "" {
			fmt.Printf("   Started: No data\n")
			fmt.Printf("   Duration: No data\n")
			fmt.Printf("   Status:  ⚪ No builds found\n")
			continue
		}

		buildResult := gjson.Parse(buildsJson)

		// Extract values with empty string defaults
		startedStr := buildResult.Get("started").String()
		duration := buildResult.Get("buildTime").String()
		status := buildResult.Get("status").String()

		// Handle empty or missing values
		if startedStr == "" {
			startedStr = "No start time available"
		}
		if duration == "" {
			duration = "No duration available"
		}
		if status == "" {
			status = "unknown"
		}

		fmt.Printf("   Started: %s\n", startedStr)
		fmt.Printf("   Duration: %s\n", duration)

		// Print status with emoji
		statusEmoji := "❓"
		switch strings.ToLower(status) {
		case "success":
			statusEmoji = "✅"
		case "failure", "error":
			statusEmoji = "❌"
		case "running":
			statusEmoji = "🔄"
		case "pending":
			statusEmoji = "⏳"
		case "terminated":
			statusEmoji = "🛑"
		}
		fmt.Printf("   Status:  %s %s\n", statusEmoji, status)
	}

	return nil
}

func (c *Command) getPipelineInfo() (string, error) {
	args := []string{"get", "pipelines", "--project", c.project, "--output", "json"}

	if c.debug {
		fmt.Printf("Executing: codefresh %s\n", strings.Join(args, " "))
	}

	cmd := exec.Command("codefresh", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		if stderr.Len() > 0 {
			return "", fmt.Errorf("command failed: %s", stderr.String())
		}
		return "", err
	}

	if stdout.Len() > 0 {
		return stdout.String(), nil
	}

	return "", fmt.Errorf("no output from command")
}

func (c *Command) getPipelineBuilds(pipelineName string) (string, error) {
	args := []string{"get", "builds", "--pipeline-name", pipelineName, "--limit", "1", "--output", "json"}

	if c.debug {
		fmt.Printf("Executing: codefresh %s\n", strings.Join(args, " "))
	}

	cmd := exec.Command("codefresh", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if stderr.Len() > 0 {
			return "", fmt.Errorf("command failed: %s", stderr.String())
		}
		return "", err
	}

	if stdout.Len() > 0 {
		return stdout.String(), nil
	}

	return "", nil
}

func writeErrorResponse(message string) {
	errorJson, _ := sjson.Set("{}", "status", "error")
	errorJson, _ = sjson.Set(errorJson, "message", message)
	fmt.Println(errorJson)
}
