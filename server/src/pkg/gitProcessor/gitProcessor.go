package gitProcessor

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
	"unicode"
)

type RepositoryModule struct {
	options Options
}

func NewRepositoryModule(opts Options) (*RepositoryModule, error) {
	return &RepositoryModule{
		options: opts,
	}, nil
}

func (m *RepositoryModule) Extract(repoPath string) ([]byte, error) {
	// Change to repository directory
	if err := os.Chdir(repoPath); err != nil {
		return nil, fmt.Errorf("failed to change to repository directory: %v", err)
	}

	// Get repository info
	repo := Repository{}

	// Get remote URL
	if url, err := m.getRemoteURL(); err == nil {
		repo.URL = url
	}

	// Get current branch
	if branch, err := m.getCurrentBranch(); err == nil {
		repo.Branch = branch
	}

	// Get latest commit
	if commit, err := m.getLatestCommit(); err == nil {
		repo.LastCommit = commit
	}

	// Get tags
	if tags, err := m.getTags(); err == nil {
		repo.Tags = tags
	}

	// Get commit history
	if commits, err := m.getCommitHistory(m.options.CommitHistoryMonths); err == nil {
		repo.CommitHistory = commits
	}

	// Get release history
	if releases, err := m.getReleaseHistory(m.options.ReleaseHistoryMonths); err == nil {
		repo.ReleaseHistory = releases
	}

	// Create initial analysis result
	result := AnalysisResult{
		Repository: repo,
		Build: BuildInfo{
			Docker: DockerConfig{
				Enabled: false,
			},
		},
		Dependencies: DependencyInfo{
			Libraries: make(map[string]string),
		},
		Documentation: DocumentationInfo{},
	}

	// Check for Dockerfile
	if _, err := os.Stat("Dockerfile"); err == nil {
		result.Build.Docker.Enabled = true
		if ports, err := m.parseDockerPorts("Dockerfile"); err == nil {
			result.Build.Docker.Ports = ports
		}
	}

	// Check for dependencies
	if err := m.detectDependencies(&result); err != nil {
		fmt.Printf("Warning: Failed to detect dependencies: %v\n", err)
	}

	// Check for documentation
	if err := m.detectDocumentation(&result.Documentation); err != nil {
		fmt.Printf("Warning: Failed to detect documentation: %v\n", err)
	}

	return json.Marshal(result)
}

func (m *RepositoryModule) Validate(data []byte) error {
	var result AnalysisResult
	return json.Unmarshal(data, &result)
}

func (m *RepositoryModule) Transform(data []byte) ([]byte, error) {
	return data, nil
}

// Git Operations

func (m *RepositoryModule) getRemoteURL() (string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func (m *RepositoryModule) getCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func (m *RepositoryModule) getLatestCommit() (Commit, error) {
	format := "--format=%H%n%an%n%aI%n%s"
	cmd := exec.Command("git", "log", "-1", format)
	output, err := cmd.Output()
	if err != nil {
		return Commit{}, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) < 4 {
		return Commit{}, fmt.Errorf("invalid commit format")
	}

	date, err := time.Parse(time.RFC3339, lines[2])
	if err != nil {
		return Commit{}, err
	}

	return Commit{
		Hash:    lines[0],
		Author:  lines[1],
		Date:    date,
		Message: lines[3],
	}, nil
}

func (m *RepositoryModule) getTags() ([]Tag, error) {
	cmd := exec.Command("git", "for-each-ref",
		"--sort=-creatordate",
		"--format=%(refname:short)%09%(taggerdate:iso8601)%09%(taggername)",
		"refs/tags")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var tags []Tag
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return nil, nil
	}

	for _, line := range lines {
		parts := strings.Split(line, "\t")
		if len(parts) < 3 {
			// Try getting commit info for lightweight tags
			tagName := strings.TrimSpace(parts[0])
			tagCmd := exec.Command("git", "show", "-s", "--format=%aI%n%an", tagName)
			tagInfo, err := tagCmd.Output()
			if err != nil {
				continue
			}
			tagLines := strings.Split(strings.TrimSpace(string(tagInfo)), "\n")
			if len(tagLines) >= 2 {
				date, err := time.Parse(time.RFC3339, tagLines[0])
				if err != nil {
					continue
				}
				tags = append(tags, Tag{
					Name:   tagName,
					Date:   date,
					Author: tagLines[1],
				})
			}
			continue
		}

		date, err := time.Parse("2006-01-02 15:04:05 -0700", strings.TrimSpace(parts[1]))
		if err != nil {
			continue
		}

		tags = append(tags, Tag{
			Name:   strings.TrimSpace(parts[0]),
			Date:   date,
			Author: strings.TrimSpace(parts[2]),
		})
	}

	return tags, nil
}

func (m *RepositoryModule) getCommitHistory(months int) ([]Commit, error) {
	since := time.Now().AddDate(0, -months, 0).Format("2006-01-02")
	format := "--format=%H%n%an%n%aI%n%s%n--COMMIT--"
	cmd := exec.Command("git", "log", fmt.Sprintf("--since=%s", since), format)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	commits := []Commit{}
	commitStrings := strings.Split(strings.TrimSpace(string(output)), "--COMMIT--")

	for _, commitStr := range commitStrings {
		if strings.TrimSpace(commitStr) == "" {
			continue
		}

		lines := strings.Split(strings.TrimSpace(commitStr), "\n")
		if len(lines) < 4 {
			continue
		}

		date, err := time.Parse(time.RFC3339, lines[2])
		if err != nil {
			continue
		}

		commits = append(commits, Commit{
			Hash:    lines[0],
			Author:  lines[1],
			Date:    date,
			Message: lines[3],
		})
	}

	return commits, nil
}

func (m *RepositoryModule) getReleaseHistory(months int) ([]Release, error) {
	since := time.Now().AddDate(0, -months, 0).Format("2006-01-02")
	cmd := exec.Command("git", "for-each-ref",
		"--sort=-creatordate",
		"--format=%(refname:short)%09%(creatordate:iso8601)%09%(subject)%09%(taggername)%(authorname)",
		"refs/tags",
		fmt.Sprintf("--since=%s", since))

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	releases := []Release{}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")

	// Get latest stable version
	latestStable := ""
	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, "\t")
		if len(parts) < 2 {
			continue
		}

		tagName := strings.TrimSpace(parts[0])
		if isStableVersion(tagName) {
			latestStable = tagName
			break
		}
	}

	// Process all releases
	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, "\t")
		if len(parts) < 2 {
			continue
		}

		tagName := strings.TrimSpace(parts[0])

		// Get detailed tag info using git show
		tagCmd := exec.Command("git", "show", "--format=%aI%n%an", tagName)
		tagInfo, err := tagCmd.Output()
		if err == nil {
			tagLines := strings.Split(strings.TrimSpace(string(tagInfo)), "\n")
			if len(tagLines) >= 2 {
				date, err := time.Parse(time.RFC3339, tagLines[0])
				if err == nil {
					name := ""
					if len(parts) > 2 {
						name = strings.TrimSpace(parts[2])
					}

					releases = append(releases, Release{
						Tag:            tagName,
						Date:           date,
						Name:           name,
						Author:         tagLines[1],
						IsLatestStable: tagName == latestStable,
					})
					continue
				}
			}
		}

		// Fallback to for-each-ref date if git show fails
		date, err := time.Parse("2006-01-02 15:04:05 -0700", strings.TrimSpace(parts[1]))
		if err != nil {
			continue
		}

		name := ""
		if len(parts) > 2 {
			name = strings.TrimSpace(parts[2])
		}

		author := ""
		if len(parts) > 3 {
			author = strings.TrimSpace(parts[3])
		}

		releases = append(releases, Release{
			Tag:            tagName,
			Date:           date,
			Name:           name,
			Author:         author,
			IsLatestStable: tagName == latestStable,
		})
	}

	return releases, nil
}

// isStableVersion checks if a version tag represents a stable release
// This assumes semantic versioning or similar versioning schemes
func isStableVersion(tag string) bool {
	// Remove 'v' prefix if present
	tag = strings.TrimPrefix(tag, "v")

	// Check if the tag starts with a number (indicating a version)
	if len(tag) == 0 || !strings.Contains(tag, ".") || !unicode.IsDigit(rune(tag[0])) {
		return false
	}

	// Check for pre-release indicators
	preReleaseIndicators := []string{
		"-alpha", "-beta", "-rc", "-dev", "-snapshot",
		".alpha", ".beta", ".rc", ".dev", ".snapshot",
		"alpha", "beta", "rc", "dev", "snapshot",
	}

	tagLower := strings.ToLower(tag)
	for _, indicator := range preReleaseIndicators {
		if strings.Contains(tagLower, indicator) {
			return false
		}
	}

	return true
}

// Additional Detection Functions

func (m *RepositoryModule) parseDockerPorts(dockerfilePath string) ([]string, error) {
	data, err := os.ReadFile(dockerfilePath)
	if err != nil {
		return nil, err
	}

	var ports []string
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(strings.ToUpper(strings.TrimSpace(line)), "EXPOSE") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				ports = append(ports, parts[1:]...)
			}
		}
	}
	return ports, nil
}

func (m *RepositoryModule) detectDependencies(result *AnalysisResult) error {
	// Check for go.mod
	if data, err := os.ReadFile("go.mod"); err == nil {
		result.Dependencies.Language = "Go"
		lines := strings.Split(string(data), "\n")
		if len(lines) > 0 {
			parts := strings.Fields(lines[0])
			if len(parts) > 1 {
				result.Dependencies.Version = parts[1]
			}
		}
	}

	// Check for package.json
	if data, err := os.ReadFile("package.json"); err == nil {
		var pkg struct {
			Dependencies    map[string]string `json:"dependencies"`
			DevDependencies map[string]string `json:"devDependencies"`
		}
		if json.Unmarshal(data, &pkg) == nil {
			result.Dependencies.Language = "JavaScript/Node.js"
			result.Dependencies.Libraries = pkg.Dependencies
			// Merge dev dependencies
			for k, v := range pkg.DevDependencies {
				result.Dependencies.Libraries[k+" (dev)"] = v
			}
		}
	}

	return nil
}

func (m *RepositoryModule) detectDocumentation(result *DocumentationInfo) error {
	// Check for README files
	readmePatterns := []string{"README.md", "README.txt", "README"}
	for _, pattern := range readmePatterns {
		if _, err := os.Stat(pattern); err == nil {
			result.Available = true
			if data, err := os.ReadFile(pattern); err == nil {
				result.Summary = string(data)
			}
			break
		}
	}

	// Check for API documentation
	apiDocPatterns := []string{
		"api.md", "API.md", "docs/api.md", "docs/API.md",
		"swagger.json", "swagger.yaml", "openapi.json", "openapi.yaml",
	}
	for _, pattern := range apiDocPatterns {
		if _, err := os.Stat(pattern); err == nil {
			result.API = true
			break
		}
	}

	return nil
}
