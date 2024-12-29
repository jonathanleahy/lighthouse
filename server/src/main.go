package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

// Move types to separate package
type RepoDetails struct {
	RepoBitUrl      string
	Namespace       string
	AppNameSuffixes map[string]bool
}

// Configuration struct to hold all app config
type Config struct {
	CacheDuration string `json:"cacheDuration"`
	ServerPort    string `json:"serverPort"`
	TokenPath     string `json:"tokenPath"`
	ProjectsPath  string `json:"projectsPath"`
	SummaryPath   string `json:"summaryPath"`
}

func loadConfig(path string) (Config, error) {
	var config Config
	file, err := os.ReadFile(path)
	if err != nil {
		return config, fmt.Errorf("error reading config file: %w", err)
	}

	if err := json.Unmarshal(file, &config); err != nil {
		return config, fmt.Errorf("error parsing config file: %w", err)
	}

	return config, nil
}

// Main HTTP handler struct
type Handler struct {
	config Config
	cache  *cache // Changed to lowercase to match the type name
}

func NewHandler(config Config) *Handler {
	return &Handler{
		config: config,
		cache:  NewCache(config.CacheDuration),
	}
}

// Main function updated to use config file
func main() {
	baseRepoName := flag.String("repo", "", "The base repository name")
	webserver := flag.Bool("webserver", false, "Run as a webserver")
	configPath := flag.String("config", "config.json", "Path to config file")
	flag.Parse()

	config, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if *webserver {
		runWebServer(config)
		return
	}

	runCLI(*baseRepoName)
}

func runWebServer(config Config) {
	handler := NewHandler(config)

	router := http.NewServeMux()
	router.HandleFunc("/", handler.handleRepoRequest)
	router.HandleFunc("/repos", handler.listReposHandler)
	router.HandleFunc("/list-repos", handler.listReposFromFileHandler)
	router.HandleFunc("/terraform-configs", handler.handleTerraformConfigs)

	log.Printf("Starting web server on %s", config.ServerPort)
	if err := http.ListenAndServe(config.ServerPort, router); err != nil {
		log.Fatalf("Error starting web server: %v", err)
	}
}
