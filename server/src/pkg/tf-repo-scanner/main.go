package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type Repository struct {
	Name        string `json:"repository_name"`
	Team        string `json:"team"`
	Description string `json:"description"`
}

func cloneOrPullRepo(repoURL, tmpDir string) (string, error) {
	// Extract repository name from URL
	parts := strings.Split(repoURL, "/")
	repoName := strings.TrimSuffix(parts[len(parts)-1], ".git")
	repoPath := filepath.Join(tmpDir, repoName)

	// Check if repository already exists
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		// Clone the repository
		cmd := exec.Command("git", "clone", repoURL, repoPath)
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("failed to clone repository: %v", err)
		}
		fmt.Printf("Cloned repository to %s\n", repoPath)
	} else {
		// Pull latest changes
		cmd := exec.Command("git", "-C", repoPath, "pull")
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("failed to pull repository: %v", err)
		}
		fmt.Printf("Updated existing repository at %s\n", repoPath)
	}

	return repoPath, nil
}

func main() {
	// Define command-line flags
	rootDir := flag.String("path", ".", "Path to the root directory containing Terraform files")
	format := flag.String("format", "json", "Output format: json or csv")
	repoURL := flag.String("repo", "", "GitHub repository URL (optional)")
	tmpDir := flag.String("tmp", os.TempDir(), "Temporary directory for cloning repositories")
	flag.Parse()

	// Handle GitHub repository URL if provided
	processingPath := *rootDir
	if *repoURL != "" {
		var err error
		processingPath, err = cloneOrPullRepo(*repoURL, *tmpDir)
		if err != nil {
			fmt.Printf("Error handling repository: %v\n", err)
			os.Exit(1)
		}
	}

	// Verify the directory exists
	if _, err := os.Stat(processingPath); os.IsNotExist(err) {
		fmt.Printf("Error: Directory '%s' does not exist\n", processingPath)
		flag.Usage()
		os.Exit(1)
	}

	repositories := []Repository{}

	// Walk through all directories starting from the provided root
	err := filepath.Walk(processingPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only process .tf files
		if !info.IsDir() && strings.HasSuffix(path, ".tf") {
			// Get the parent directory name as the team
			team := filepath.Base(filepath.Dir(path))

			// Read and parse the file
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			// Extract repositories from the file
			fileRepos := parseRepositories(string(content), team)
			repositories = append(repositories, fileRepos...)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking through directories: %v\n", err)
		return
	}

	// Output results based on format flag
	switch strings.ToLower(*format) {
	case "json":
		outputJSON(repositories)
	case "csv":
		outputCSV(repositories)
	default:
		fmt.Printf("Error: Unsupported format '%s'. Use 'json' or 'csv'\n", *format)
		flag.Usage()
		os.Exit(1)
	}
}

func outputJSON(repositories []Repository) {
	// Create a wrapper struct for the JSON output
	output := struct {
		Repositories []Repository `json:"repositories"`
		Count        int          `json:"total_count"`
	}{
		Repositories: repositories,
		Count:        len(repositories),
	}

	// Convert to JSON with indentation for readability
	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Printf("Error converting to JSON: %v\n", err)
		return
	}

	fmt.Println(string(jsonData))
}

func outputCSV(repositories []Repository) {
	fmt.Println("repository_name,team,description")
	fmt.Println("------------------------------------")
	for _, repo := range repositories {
		fmt.Printf("%s,%s,%s\n", repo.Name, repo.Team, repo.Description)
	}
}

func parseRepositories(content, team string) []Repository {
	repositories := []Repository{}

	// Regular expressions to extract information
	moduleRegex := regexp.MustCompile(`module\s+"([^"]+)"\s*{([^}]+)}`)
	repoNameRegex := regexp.MustCompile(`repository_name\s*=\s*"([^"]+)"`)
	descRegex := regexp.MustCompile(`description\s*=\s*"([^"]+)"`)

	// Find all module blocks
	moduleMatches := moduleRegex.FindAllStringSubmatch(content, -1)

	for _, moduleMatch := range moduleMatches {
		if len(moduleMatch) < 2 {
			continue
		}

		moduleBlock := moduleMatch[2]

		// Extract repository name
		repoNameMatch := repoNameRegex.FindStringSubmatch(moduleBlock)
		if len(repoNameMatch) < 2 {
			continue
		}
		repoName := repoNameMatch[1]

		// Extract description
		description := ""
		descMatch := descRegex.FindStringSubmatch(moduleBlock)
		if len(descMatch) >= 2 {
			description = descMatch[1]
		}

		repositories = append(repositories, Repository{
			Name:        repoName,
			Team:        team,
			Description: description,
		})
	}

	return repositories
}
