// pkg/regions/parser.go
package regions

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type RegionDetails struct {
	Path           string `json:"path"`
	RegionDefault  string `json:"region_default"`
	AccountDefault string `json:"account_default"`
	Namespace      string `json:"namespace"`
}

// ParseRegions returns the region configuration for a given repository
func ParseRegions(baseRepoName string) ([]RegionDetails, error) {
	var configs []RegionDetails

	// Construct path to repo's terraform configs
	rootDir := filepath.Join("projects/projects", baseRepoName, "github/scripts/terraform")

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip if not a directory
		if !info.IsDir() {
			return nil
		}

		// Check if directory contains _variables.tf
		varsPath := filepath.Join(path, "_variables.tf")
		if !fileExists(varsPath) {
			return nil
		}

		// Parse the configuration
		config := RegionDetails{
			Path: path,
		}

		// Read variables file
		varsContent, err := os.ReadFile(varsPath)
		if err != nil {
			return err
		}

		// Extract values using regex
		config.RegionDefault = extractValue(string(varsContent), "variable \"region\"", "default")
		config.AccountDefault = extractValue(string(varsContent), "variable \"account\"", "default")

		// Try to get namespace from providers.tf if it exists
		providersPath := filepath.Join(path, "providers.tf")
		if fileExists(providersPath) {
			if content, err := os.ReadFile(providersPath); err == nil {
				config.Namespace = extractNamespace(string(content))
			}
		}

		// Only add if we have both region and account values
		if config.RegionDefault != "" && config.AccountDefault != "" {
			configs = append(configs, config)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking through terraform configs: %v", err)
	}

	return configs, nil
}

func extractValue(content, blockName, fieldName string) string {
	pattern := fmt.Sprintf(`%s\s*{[^}]*%s\s*=\s*"([^"]*)"`, blockName, fieldName)
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(content)

	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func extractNamespace(content string) string {
	re := regexp.MustCompile(`Squad\s*=\s*"([^"]*)"`)
	matches := re.FindStringSubmatch(content)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
