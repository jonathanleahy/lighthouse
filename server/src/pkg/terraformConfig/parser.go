// pkg/terraformConfig/parser.go

package terraformConfig

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type TerraformConfig struct {
	Path           string `json:"path"`
	RegionDefault  string `json:"region_default"`
	AccountDefault string `json:"account_default"`
	Namespace      string `json:"namespace"`
}

// ParseConfigs reads and parses Terraform configurations from the given root directory
func ParseConfigs(rootDir string) ([]TerraformConfig, error) {
	var configs []TerraformConfig

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
		config := TerraformConfig{
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

		if config.RegionDefault != "" && config.AccountDefault != "" {
			configs = append(configs, config)
		}

		return nil
	})

	return configs, err
}

// ToJSON converts the configurations to a JSON string
func ToJSON(configs []TerraformConfig) (string, error) {
	jsonData, err := json.MarshalIndent(configs, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error converting to JSON: %v", err)
	}
	return string(jsonData), nil
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