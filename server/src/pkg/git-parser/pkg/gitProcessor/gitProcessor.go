// pkg/gitProcessor/gitProcessor.go
package gitProcessor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// RepoData represents the raw data extracted from the repository
type RepoData struct {
	Git struct {
		LastCommit struct {
			Hash    string `json:"hash"`
			Author  string `json:"author"`
			Date    string `json:"date"`
			Message string `json:"message"`
		} `json:"lastCommit"`
		Branch    string   `json:"branch"`
		Tags      []string `json:"tags"`
		RemoteURL string   `json:"remoteUrl"`
	} `json:"git"`

	Build struct {
		Docker struct {
			Present   bool     `json:"present"`
			BaseImage string   `json:"baseImage,omitempty"`
			Ports     []int    `json:"ports,omitempty"`
			Commands  []string `json:"commands,omitempty"`
		} `json:"docker"`
		Makefile struct {
			Present bool     `json:"present"`
			Targets []string `json:"targets,omitempty"`
		} `json:"makefile"`
	} `json:"build"`

	Dependencies struct {
		Go struct {
			Version  string   `json:"version"`
			Module   string   `json:"module"`
			Direct   []string `json:"direct"`
			Indirect []string `json:"indirect"`
		} `json:"go"`
	} `json:"dependencies"`

	Documentation struct {
		ReadmePresent bool   `json:"readmePresent"`
		ApiDocs       bool   `json:"apiDocs"`
		Summary       string `json:"summary,omitempty"`
	} `json:"documentation"`
}

type RepositoryModule struct {
	name string
}

func NewRepositoryModule() (*RepositoryModule, error) {
	return &RepositoryModule{
		name: "repository",
	}, nil
}

func (r *RepositoryModule) Extract(path string) (json.RawMessage, error) {
	var repoData RepoData

	// Extract Git information
	if err := r.extractGitInfo(path, &repoData); err != nil {
		return nil, fmt.Errorf("git info extraction failed: %w", err)
	}

	// Extract build configuration
	if err := r.extractBuildInfo(path, &repoData); err != nil {
		return nil, fmt.Errorf("build info extraction failed: %w", err)
	}

	// Extract dependencies
	if err := r.extractDependencies(path, &repoData); err != nil {
		return nil, fmt.Errorf("dependencies extraction failed: %w", err)
	}

	// Extract documentation
	if err := r.extractDocumentation(path, &repoData); err != nil {
		return nil, fmt.Errorf("documentation extraction failed: %w", err)
	}

	return json.Marshal(repoData)
}

func (r *RepositoryModule) extractGitInfo(path string, data *RepoData) error {
	// Get last commit
	cmd := exec.Command("git", "-C", path, "log", "-1", "--pretty=format:%H|%an|%aI|%s")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get git log: %w", err)
	}

	parts := strings.Split(string(output), "|")
	if len(parts) == 4 {
		data.Git.LastCommit.Hash = parts[0]
		data.Git.LastCommit.Author = parts[1]
		data.Git.LastCommit.Date = parts[2]
		data.Git.LastCommit.Message = parts[3]
	}

	// Get current branch
	cmd = exec.Command("git", "-C", path, "rev-parse", "--abbrev-ref", "HEAD")
	if output, err := cmd.Output(); err == nil {
		data.Git.Branch = strings.TrimSpace(string(output))
	}

	// Get tags
	cmd = exec.Command("git", "-C", path, "tag", "--sort=-creatordate")
	if output, err := cmd.Output(); err == nil && len(output) > 0 {
		data.Git.Tags = strings.Split(strings.TrimSpace(string(output)), "\n")
	}

	// Get remote URL
	cmd = exec.Command("git", "-C", path, "remote", "get-url", "origin")
	if output, err := cmd.Output(); err == nil {
		data.Git.RemoteURL = strings.TrimSpace(string(output))
	}

	return nil
}

func (r *RepositoryModule) extractBuildInfo(path string, data *RepoData) error {
	// Check Dockerfile
	dockerfilePath := filepath.Join(path, "Dockerfile")
	if content, err := ioutil.ReadFile(dockerfilePath); err == nil {
		data.Build.Docker.Present = true

		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)

			// Extract base image
			if strings.HasPrefix(line, "FROM ") {
				data.Build.Docker.BaseImage = strings.TrimPrefix(line, "FROM ")
			}

			// Extract exposed ports
			if strings.HasPrefix(line, "EXPOSE ") {
				ports := strings.TrimPrefix(line, "EXPOSE ")
				for _, port := range strings.Fields(ports) {
					var portNum int
					if _, err := fmt.Sscanf(port, "%d", &portNum); err == nil {
						data.Build.Docker.Ports = append(data.Build.Docker.Ports, portNum)
					}
				}
			}

			// Extract CMD and ENTRYPOINT
			if strings.HasPrefix(line, "CMD ") || strings.HasPrefix(line, "ENTRYPOINT ") {
				data.Build.Docker.Commands = append(data.Build.Docker.Commands, line)
			}
		}
	}

	// Check Makefile
	makefilePath := filepath.Join(path, "Makefile")
	if content, err := ioutil.ReadFile(makefilePath); err == nil {
		data.Build.Makefile.Present = true

		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			if strings.Contains(line, ":") && !strings.HasPrefix(line, "\t") {
				target := strings.TrimSpace(strings.Split(line, ":")[0])
				if target != "" && !strings.Contains(target, " ") {
					data.Build.Makefile.Targets = append(data.Build.Makefile.Targets, target)
				}
			}
		}
	}

	return nil
}

func (r *RepositoryModule) extractDependencies(path string, data *RepoData) error {
	// Check for go1.mod1
	goModPath := filepath.Join(path, "go1.mod1")
	if content, err := ioutil.ReadFile(goModPath); err == nil {
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "module ") {
				data.Dependencies.Go.Module = strings.TrimPrefix(line, "module ")
			} else if strings.HasPrefix(line, "go ") {
				data.Dependencies.Go.Version = strings.TrimPrefix(line, "go ")
			}
		}

		// Get dependencies using go list
		cmd := exec.Command("go", "list", "-m", "all")
		cmd.Dir = path
		if output, err := cmd.Output(); err == nil {
			deps := strings.Split(string(output), "\n")
			for _, dep := range deps {
				dep = strings.TrimSpace(dep)
				if dep != "" && !strings.HasPrefix(dep, data.Dependencies.Go.Module) {
					data.Dependencies.Go.Direct = append(data.Dependencies.Go.Direct, dep)
				}
			}
		}
	}

	return nil
}

func (r *RepositoryModule) extractDocumentation(path string, data *RepoData) error {
	// Check README files
	readmePaths := []string{
		filepath.Join(path, "README.md"),
		filepath.Join(path, "README"),
		filepath.Join(path, "readme.md"),
	}

	for _, readmePath := range readmePaths {
		if content, err := ioutil.ReadFile(readmePath); err == nil {
			data.Documentation.ReadmePresent = true
			// Extract first paragraph as summary
			paragraphs := strings.Split(string(content), "\n\n")
			if len(paragraphs) > 0 {
				data.Documentation.Summary = strings.TrimSpace(paragraphs[0])
			}
			break
		}
	}

	// Check for API documentation
	apiPaths := []string{
		filepath.Join(path, "api"),
		filepath.Join(path, "docs", "api"),
		filepath.Join(path, "swagger.yaml"),
		filepath.Join(path, "swagger.json"),
		filepath.Join(path, "openapi.yaml"),
		filepath.Join(path, "openapi.json"),
	}

	for _, apiPath := range apiPaths {
		if _, err := os.Stat(apiPath); err == nil {
			data.Documentation.ApiDocs = true
			break
		}
	}

	return nil
}

func (r *RepositoryModule) Validate(data json.RawMessage) error {
	var repoData RepoData
	if err := json.Unmarshal(data, &repoData); err != nil {
		return fmt.Errorf("invalid repository data format: %w", err)
	}

	// Basic validation
	if repoData.Git.LastCommit.Hash == "" {
		return fmt.Errorf("repository must have at least one commit")
	}

	if repoData.Git.Branch == "" {
		return fmt.Errorf("repository must have a current branch")
	}

	return nil
}

func (r *RepositoryModule) Transform(data json.RawMessage) (json.RawMessage, error) {
	var repoData RepoData
	if err := json.Unmarshal(data, &repoData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal repo data: %w", err)
	}

	// Transform to standardized format matching the AnalysisResult struct in main.go
	transformed := struct {
		Repository struct {
			URL        string `json:"url"`
			Branch     string `json:"branch"`
			LastCommit struct {
				Hash    string `json:"hash"`
				Author  string `json:"author"`
				Date    string `json:"date"`
				Message string `json:"message"`
			} `json:"lastCommit"`
			Tags []string `json:"tags"`
		} `json:"repository"`
		Build struct {
			Docker struct {
				Enabled bool  `json:"enabled"`
				Ports   []int `json:"ports,omitempty"`
			} `json:"docker"`
			Commands []string `json:"commands,omitempty"`
		} `json:"build"`
		Dependencies struct {
			Language  string   `json:"language"`
			Version   string   `json:"version"`
			Libraries []string `json:"libraries"`
		} `json:"dependencies"`
		Documentation struct {
			Available bool   `json:"available"`
			API       bool   `json:"api"`
			Summary   string `json:"summary,omitempty"`
		} `json:"documentation"`
	}{}

	// Fill transformed structure
	transformed.Repository.URL = repoData.Git.RemoteURL
	transformed.Repository.Branch = repoData.Git.Branch
	transformed.Repository.LastCommit = repoData.Git.LastCommit
	transformed.Repository.Tags = repoData.Git.Tags

	transformed.Build.Docker.Enabled = repoData.Build.Docker.Present
	transformed.Build.Docker.Ports = repoData.Build.Docker.Ports
	if repoData.Build.Makefile.Present {
		transformed.Build.Commands = repoData.Build.Makefile.Targets
	}

	if repoData.Dependencies.Go.Version != "" {
		transformed.Dependencies.Language = "go"
		transformed.Dependencies.Version = repoData.Dependencies.Go.Version
		transformed.Dependencies.Libraries = repoData.Dependencies.Go.Direct
	}

	transformed.Documentation.Available = repoData.Documentation.ReadmePresent
	transformed.Documentation.API = repoData.Documentation.ApiDocs
	transformed.Documentation.Summary = repoData.Documentation.Summary

	return json.Marshal(transformed)
}
