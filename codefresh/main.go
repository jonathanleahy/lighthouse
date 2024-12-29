package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type Command struct {
	fs      *flag.FlagSet
	project string
	output  string
	debug   bool
}

type PipelineInfo struct {
	Name     string
	Started  string
	Duration string
	Tag      string
	Status   string
	Error    error
}

func getBuildVariables(buildId string, progressId string) (map[string]interface{}, error) {
	apiKey := "677090d6f12fd1f488cd8d88.2aea720444e3c14b848c1633b9c98652"
	if apiKey == "" {
		return nil, fmt.Errorf("CF_API_KEY environment variable not set")
	}

	url := fmt.Sprintf("https://g.codefresh.io/api/builds/%s", buildId)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	result := gjson.Parse(string(body))

	var variables = make(map[string]interface{})
	result.Get("exposedVariables.pipeline").ForEach(func(key, value gjson.Result) bool {
		if value.Get("key").String() == "tag" {
			variables["tag"] = value.Get("value").String()
			return false
		}
		return true
	})

	return variables, nil
}

func (c *Command) getPipelineInfo() ([]gjson.Result, error) {
	args := []string{"get", "pipelines", "--project", c.project, "--output", "json"}

	cmd := exec.Command("codefresh", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		if stderr.Len() > 0 {
			return nil, fmt.Errorf("command failed: %s", stderr.String())
		}
		return nil, err
	}

	if stdout.Len() > 0 {
		pipelineResult := gjson.Parse(stdout.String())
		var pipelineArray []gjson.Result
		if pipelineResult.IsArray() {
			pipelineArray = pipelineResult.Array()
		} else {
			pipelineArray = []gjson.Result{pipelineResult}
		}
		return pipelineArray, nil
	}

	return nil, fmt.Errorf("no output from command")
}

func (c *Command) getPipelineBuilds(pipelineName string) (string, error) {
	args := []string{
		"get",
		"builds",
		"--pipeline-name",
		pipelineName,
		"--limit",
		"1",
		"--output",
		"json",
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

func (c *Command) processPipeline(p gjson.Result, resultChan chan<- PipelineInfo, wg *sync.WaitGroup) {
	defer wg.Done()

	metadata := p.Get("metadata")
	if !metadata.Exists() {
		return
	}

	name := metadata.Get("name").String()
	pipelineInfo := PipelineInfo{
		Name: name,
	}

	buildsJson, err := c.getPipelineBuilds(name)
	if err != nil {
		pipelineInfo.Error = fmt.Errorf("failed to get builds: %v", err)
		pipelineInfo.Started = "Unable to fetch"
		pipelineInfo.Duration = "Unable to fetch"
		pipelineInfo.Tag = "Unable to fetch"
		pipelineInfo.Status = "❌ Error"
		resultChan <- pipelineInfo
		return
	}

	if buildsJson == "" {
		pipelineInfo.Started = "No data"
		pipelineInfo.Duration = "No data"
		pipelineInfo.Tag = "No data"
		pipelineInfo.Status = "⚪ No builds found"
		resultChan <- pipelineInfo
		return
	}

	buildResult := gjson.Parse(buildsJson)

	buildId := buildResult.Get("id").String()
	progressId := buildResult.Get("progress").String()

	// Get build variables
	variables, err := getBuildVariables(buildId, progressId)
	var tag string
	if err != nil {
		tag = "Unable to fetch tag"
	} else {
		if tagVal, ok := variables["tag"]; ok {
			tag = fmt.Sprintf("%v", tagVal)
		} else {
			tag = "No tag found"
		}
	}

	pipelineInfo.Started = buildResult.Get("started").String()
	if pipelineInfo.Started == "" {
		pipelineInfo.Started = "No start time available"
	}

	pipelineInfo.Duration = buildResult.Get("buildTime").String()
	if pipelineInfo.Duration == "" {
		pipelineInfo.Duration = "No duration available"
	}

	pipelineInfo.Tag = tag

	status := buildResult.Get("status").String()
	if status == "" {
		status = "unknown"
	}

	// Determine status emoji
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

	pipelineInfo.Status = fmt.Sprintf("%s %s", statusEmoji, status)

	resultChan <- pipelineInfo
}

func (c *Command) Run() error {
	if _, err := exec.LookPath("codefresh"); err != nil {
		return fmt.Errorf("Codefresh CLI is not installed")
	}

	pipelineArray, err := c.getPipelineInfo()
	if err != nil {
		return fmt.Errorf("Failed to retrieve pipeline info: %v", err)
	}

	fmt.Println("\nPipeline Execution Status:")
	fmt.Println("==========================")

	// Create a channel to receive pipeline results
	resultChan := make(chan PipelineInfo, len(pipelineArray))
	var wg sync.WaitGroup

	// Process pipelines concurrently
	for _, p := range pipelineArray {
		wg.Add(1)
		go c.processPipeline(p, resultChan, &wg)
	}

	// Close the channel when all goroutines are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect and print results
	for pinfo := range resultChan {
		fmt.Printf("\n📋 Pipeline: %s\n", pinfo.Name)
		fmt.Printf("   Started: %s\n", pinfo.Started)
		fmt.Printf("   Duration: %s\n", pinfo.Duration)
		fmt.Printf("   Tag: %s\n", pinfo.Tag)

		if pinfo.Error != nil {
			fmt.Printf("   Status:  %s\n", pinfo.Status)
		} else {
			fmt.Printf("   Status:  %s\n", pinfo.Status)
		}
	}

	return nil
}

func writeErrorResponse(message string) {
	errorJson, _ := sjson.Set("{}", "status", "error")
	errorJson, _ = sjson.Set(errorJson, "message", message)
	fmt.Println(errorJson)
}

func main() {
	listCmd := &Command{
		fs: flag.NewFlagSet("list", flag.ExitOnError),
	}

	listCmd.fs.StringVar(&listCmd.project, "project", "", "Project name")
	listCmd.fs.StringVar(&listCmd.output, "output", "", "Output file path")
	listCmd.fs.BoolVar(&listCmd.debug, "debug", false, "Enable debug output (no effect)")

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
