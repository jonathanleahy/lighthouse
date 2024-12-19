package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"git-parser/pkg/gitProcessor"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	defaultHistoryMonths = 24 // 2 years default for both commit and release history
)

// Command line flags for history
type HistoryFlags struct {
	commitMonths  int
	releaseMonths int
}

func main() {
	// Setup logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Parse command line flags
	repoPath := flag.String("path", ".", "Path to the repository to analyze")
	outputPath := flag.String("output", "", "Path to write JSON output (optional)")
	prettyPrint := flag.Bool("pretty", true, "Pretty print the output")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")

	// Add history flags with 2-year defaults
	commitHistory := flag.Int("commit-history", defaultHistoryMonths, "Number of months of commit history to include (default 24 months)")
	releaseHistory := flag.Int("release-history", defaultHistoryMonths, "Number of months of release history to include (default 24 months)")

	flag.Parse()

	historyFlags := HistoryFlags{
		commitMonths:  *commitHistory,
		releaseMonths: *releaseHistory,
	}

	// Create logger based on verbose flag
	logger := log.New(os.Stdout, "", log.LstdFlags)
	if !*verbose {
		logger.SetOutput(os.Stderr)
	}

	// Validate repository path
	absPath, err := filepath.Abs(*repoPath)
	if err != nil {
		logger.Fatalf("Invalid repository path: %v", err)
	}

	// Verify directory exists and is a git repository
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		logger.Fatalf("Repository path does not exist: %s", absPath)
	}
	if _, err := os.Stat(filepath.Join(absPath, ".git")); os.IsNotExist(err) {
		logger.Fatalf("Not a git repository: %s", absPath)
	}

	// Initialize repository module with history options
	logger.Printf("Initializing repository analyzer for: %s", absPath)
	repoModule, err := gitProcessor.NewRepositoryModule(gitProcessor.Options{
		CommitHistoryMonths:  historyFlags.commitMonths,
		ReleaseHistoryMonths: historyFlags.releaseMonths,
	})
	if err != nil {
		logger.Fatalf("Failed to initialize repository module: %v", err)
	}

	// Extract data
	logger.Print("Extracting repository data...")
	rawData, err := repoModule.Extract(absPath)
	if err != nil {
		logger.Fatalf("Failed to extract repository data: %v", err)
	}

	// Validate data
	logger.Print("Validating extracted data...")
	if err := repoModule.Validate(rawData); err != nil {
		logger.Fatalf("Data validation failed: %v", err)
	}

	// Transform data
	logger.Print("Transforming data to standard format...")
	transformedData, err := repoModule.Transform(rawData)
	if err != nil {
		logger.Fatalf("Data transformation failed: %v", err)
	}

	// Create final result
	var result gitProcessor.AnalysisResult
	if err := json.Unmarshal(transformedData, &result); err != nil {
		logger.Fatalf("Failed to parse transformed data: %v", err)
	}

	// Add metadata
	result.Metadata.AnalyzedAt = time.Now()
	result.Metadata.RepoPath = absPath
	result.Metadata.Status = "success"

	// Marshal the final result
	var output []byte
	if *prettyPrint {
		output, err = json.MarshalIndent(result, "", "  ")
	} else {
		output, err = json.Marshal(result)
	}
	if err != nil {
		logger.Fatalf("Failed to marshal result: %v", err)
	}

	// Handle output
	if *outputPath != "" {
		// Ensure output directory exists
		outputDir := filepath.Dir(*outputPath)
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			logger.Fatalf("Failed to create output directory: %v", err)
		}

		// Write to file
		if err := os.WriteFile(*outputPath, output, 0644); err != nil {
			logger.Fatalf("Failed to write output file: %v", err)
		}
		logger.Printf("Analysis results written to: %s", *outputPath)
	} else {
		// Print to stdout
		fmt.Println(string(output))
	}

	// Print summary if verbose
	if *verbose {
		printSummary(&result)
	}

	logger.Print("Analysis completed successfully")
}

func printSummary(result *gitProcessor.AnalysisResult) {
	timeFormat := "2006-01-02 15:04:05"

	fmt.Println("\nRepository Analysis Summary")
	fmt.Println("==========================")

	// Repository Info
	if result.Repository.URL != "" {
		fmt.Printf("Repository URL: %s\n", result.Repository.URL)
	}
	fmt.Printf("Current Branch: %s\n", result.Repository.Branch)
	fmt.Printf("Latest Commit: %s\n", result.Repository.LastCommit.Hash)
	fmt.Printf("Commit Author: %s\n", result.Repository.LastCommit.Author)
	fmt.Printf("Commit Date: %s\n", result.Repository.LastCommit.Date)
	fmt.Printf("Commit Message: %s\n", result.Repository.LastCommit.Message)

	if len(result.Repository.Tags) > 0 {
		fmt.Println("\nTags:")
		for _, tag := range result.Repository.Tags {
			fmt.Printf("- %s | %s | by %s\n",
				tag.Date.Format(timeFormat),
				tag.Name,
				tag.Author)
		}
	}

	// Commit History
	if len(result.Repository.CommitHistory) > 0 {
		fmt.Println("\nRecent Commits:")
		for _, commit := range result.Repository.CommitHistory {
			fmt.Printf("- %s | %s | %s | %s\n",
				commit.Date.Format(timeFormat),
				commit.Hash[:8],
				commit.Author,
				commit.Message)
		}
	}

	// Release History
	if len(result.Repository.ReleaseHistory) > 0 {
		fmt.Println("\nRecent Releases:")
		for _, release := range result.Repository.ReleaseHistory {
			fmt.Printf("- %s | %s | %s\n",
				release.Date.Format(timeFormat),
				release.Tag,
				release.Name)
		}
	}

	// Build Info
	fmt.Println("\nBuild Configuration:")
	fmt.Printf("- Docker Enabled: %v\n", result.Build.Docker.Enabled)
	if len(result.Build.Docker.Ports) > 0 {
		fmt.Printf("- Exposed Ports: %v\n", result.Build.Docker.Ports)
	}
	if len(result.Build.Commands) > 0 {
		fmt.Printf("- Build Commands: %v\n", result.Build.Commands)
	}

	// Dependencies
	if result.Dependencies.Language != "" {
		fmt.Println("\nDependencies:")
		fmt.Printf("- Language: %s %s\n", result.Dependencies.Language, result.Dependencies.Version)
		fmt.Printf("- Number of Dependencies: %d\n", len(result.Dependencies.Libraries))
	}

	// Documentation
	fmt.Println("\nDocumentation:")
	fmt.Printf("- README Available: %v\n", result.Documentation.Available)
	fmt.Printf("- API Documentation: %v\n", result.Documentation.API)
	if result.Documentation.Summary != "" {
		fmt.Printf("\nRepository Summary:\n%s\n", result.Documentation.Summary)
	}

	// Metadata
	fmt.Println("\nAnalysis Metadata:")
	fmt.Printf("- Analyzed At: %s\n", result.Metadata.AnalyzedAt.Format(time.RFC3339))
	fmt.Printf("- Repository Path: %s\n", result.Metadata.RepoPath)
	fmt.Printf("- Analysis Status: %s\n", result.Metadata.Status)
}
