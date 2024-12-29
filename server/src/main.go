package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Types
type RepoDetails struct {
	RepoBitUrl      string
	Namespace       string
	AppNameSuffixes map[string]bool
}

type Config struct {
	CacheDuration string `json:"cacheDuration"`
	ServerPort    string `json:"serverPort"`
	TokenPath     string `json:"tokenPath"`
	ProjectsPath  string `json:"projectsPath"`
	SummaryPath   string `json:"summaryPath"`
}

type Repository struct {
	RepositoryName string `json:"repository_name"`
	Team           string `json:"team"`
	Description    string `json:"description"`
}

// Global variables
var repoDetailsArray = []struct {
	BaseRepoName string
	Details      RepoDetails
}{
	{
		BaseRepoName: "example-repo",
		Details: RepoDetails{
			RepoBitUrl: "https://github.com/org/example-repo",
			Namespace:  "default",
			AppNameSuffixes: map[string]bool{
				"-prod": true,
				"-dev":  true,
			},
		},
	},
}

// Cache implementation
type cache struct {
	duration time.Duration
	mutex    sync.RWMutex
	data     map[string]cacheEntry
}

type cacheEntry struct {
	data      []byte
	timestamp time.Time
}

type PismoData struct {
	Repositories []Repository `json:"repositories"`
	TotalCount   int          `json:"total_count"`
}

func GetRepositoryBlock(repoName string) (*Repository, error) {
	// Read the content of pismo.json
	file, err := os.Open("projects/projects/pismo.json")
	if err != nil {
		return nil, fmt.Errorf("error opening pismo.json: %v", err)
	}
	defer file.Close()

	// Parse the JSON data
	var data PismoData
	byteValue, _ := ioutil.ReadAll(file)
	if err := json.Unmarshal(byteValue, &data); err != nil {
		return nil, fmt.Errorf("error parsing pismo.json: %v", err)
	}

	// Search for the repository by name
	for _, repo := range data.Repositories {
		if repo.RepositoryName == repoName {
			return &repo, nil
		}
	}

	return nil, fmt.Errorf("repository %s not found", repoName)
}

func NewCache(duration string) *cache {
	d, err := time.ParseDuration(duration)
	if err != nil {
		d = 5 * time.Minute // Default cache duration
	}
	return &cache{
		duration: d,
		data:     make(map[string]cacheEntry),
	}
}

func (c *cache) Get(key string) ([]byte, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	entry, exists := c.data[key]
	if !exists {
		return nil, false
	}

	if time.Since(entry.timestamp) > c.duration {
		return nil, false
	}

	return entry.data, true
}

func (c *cache) Set(key string, data []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data[key] = cacheEntry{
		data:      data,
		timestamp: time.Now(),
	}
}

// Configuration loading
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

// Handler implementation
type Handler struct {
	config Config
	cache  *cache
}

func NewHandler(config Config) *Handler {
	return &Handler{
		config: config,
		cache:  NewCache(config.CacheDuration),
	}
}

// Repository operations
func checkAndPullRepo(baseRepoName string) error {
	repoPath := filepath.Join("projects/projects", baseRepoName)
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		cmd := exec.Command("git", "clone", fmt.Sprintf("https://github.com/org/%s.git", baseRepoName), repoPath)
		return cmd.Run()
	}

	cmd := exec.Command("git", "-C", repoPath, "pull")
	return cmd.Run()
}

func getRepoFileDetails(baseRepoName string) (string, string, map[string]bool) {
	for _, repo := range repoDetailsArray {
		if repo.BaseRepoName == baseRepoName {
			return repo.Details.RepoBitUrl, repo.Details.Namespace, repo.Details.AppNameSuffixes
		}
	}

	regionsFilePath := filepath.Join("projects/projects", baseRepoName, "regions.json")
	data, err := os.ReadFile(regionsFilePath)
	if err != nil {
		return "", "", nil
	}

	var regions []struct {
		Path           string `json:"path"`
		RegionDefault  string `json:"region_default"`
		AccountDefault string `json:"account_default"`
		Namespace      string `json:"namespace"`
	}

	if err := json.Unmarshal(data, &regions); err != nil {
		return "", "", nil
	}

	repoBitUrl := fmt.Sprintf("https://github.com/org/%s", baseRepoName)
	appNameSuffixes := make(map[string]bool)
	namespace := ""

	for _, region := range regions {
		suffix := fmt.Sprintf("-%s-%s", region.AccountDefault, region.RegionDefault)
		appNameSuffixes[suffix] = true
		if namespace == "" {
			namespace = region.Namespace
		}
	}

	return repoBitUrl, namespace, appNameSuffixes
}

func processRepoData(baseRepoName, repoBitUrl, namespace string, appNameSuffixes map[string]bool, forceRefresh bool) ([]byte, error) {
	if !forceRefresh {
		summaryPath := filepath.Join("projects/projects-summary", fmt.Sprintf("%s.json", baseRepoName))
		if fileInfo, err := os.Stat(summaryPath); err == nil {
			if time.Since(fileInfo.ModTime()) < 5*time.Minute {
				return os.ReadFile(summaryPath)
			}
		}
	}

	repoData := map[string]interface{}{
		"repoName":   baseRepoName,
		"repoBitUrl": repoBitUrl,
		"namespace":  namespace,
		"apps":       []interface{}{},
	}

	repo, err := GetRepositoryBlock(baseRepoName)
	if err == nil {
		repoData["repoDesc"] = repo.Description
		repoData["repoSquad"] = repo.Team
	}

	repoData["repoCodefresh"] = fmt.Sprintf("https://g.codefresh.io/pipelines/all/?filter=pageSize:10;field:name~Name;order:asc~Asc;search:%s", baseRepoName)
	repoData["argocd"] = map[string]string{
		"url": fmt.Sprintf("https://argocd.services/applications?search=%s&showFavorites=false&proj=&sync=&autoSync=&health=&namespace=&cluster=&labels=", baseRepoName),
	}

	for suffix := range appNameSuffixes {
		appName := baseRepoName + suffix
		app := map[string]interface{}{
			"appName": appName,
			"type":    "primary",
		}
		repoData["apps"] = append(repoData["apps"].([]interface{}), app)
	}

	jsonData, err := json.MarshalIndent(repoData, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error marshalling to JSON: %w", err)
	}

	summaryPath := filepath.Join("projects/projects-summary", fmt.Sprintf("%s.json", baseRepoName))
	if err := os.WriteFile(summaryPath, jsonData, 0644); err != nil {
		return nil, fmt.Errorf("error writing cache file: %w", err)
	}

	return jsonData, nil
}

// HTTP Handlers
func (h *Handler) handleRepoRequest(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w, r)

	if r.Method == http.MethodOptions {
		return
	}

	forceRefresh := r.URL.Query().Get("force") == "true"
	baseRepoName := r.URL.Query().Get("repo")
	if baseRepoName == "" {
		http.Error(w, "Missing repo parameter", http.StatusBadRequest)
		return
	}

	if err := checkAndPullRepo(baseRepoName); err != nil {
		http.Error(w, fmt.Sprintf("Error checking repo: %v", err), http.StatusInternalServerError)
		return
	}

	repoBitUrl, namespace, appNameSuffixes := getRepoFileDetails(baseRepoName)
	if repoBitUrl == "" {
		http.Error(w, "Unknown baseRepoName", http.StatusBadRequest)
		return
	}

	jsonData, err := processRepoData(baseRepoName, repoBitUrl, namespace, appNameSuffixes, forceRefresh)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (h *Handler) listReposHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w, r)

	if r.Method == http.MethodOptions {
		return
	}

	var repoNames []string
	for _, repo := range repoDetailsArray {
		repoNames = append(repoNames, repo.BaseRepoName)
	}

	jsonData, err := json.Marshal(repoNames)
	if err != nil {
		http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (h *Handler) listReposFromFileHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w, r)

	if r.Method == http.MethodOptions {
		return
	}

	forceRefresh := r.URL.Query().Get("force") == "true"
	cacheKey := "repos_list"

	if !forceRefresh {
		if data, exists := h.cache.Get(cacheKey); exists {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Cache", "HIT")
			w.Write(data)
			return
		}
	}

	data, err := processRepositoriesList()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.cache.Set(cacheKey, data)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Cache", "MISS")
	w.Write(data)
}

func (h *Handler) handleTerraformConfigs(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w, r)

	if r.Method == http.MethodOptions {
		return
	}

	baseRepoName := r.URL.Query().Get("repo")
	if baseRepoName == "" {
		http.Error(w, "Missing repo parameter", http.StatusBadRequest)
		return
	}

	configPath := filepath.Join(h.config.ProjectsPath, baseRepoName, "scripts/terraform")
	configs, err := processTerrformConfigs(configPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error processing terraform configs: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(configs)
}

// Helper functions
func setCORSHeaders(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if strings.HasPrefix(origin, "http://localhost:") {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func processRepositoriesList() ([]byte, error) {
	jsonData, err := os.ReadFile("projects/projects/pismo.json")
	if err != nil {
		return nil, fmt.Errorf("error reading pismo.json: %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, fmt.Errorf("error parsing pismo.json: %v", err)
	}

	reposInterface, ok := data["repositories"]
	if !ok {
		return nil, fmt.Errorf("no 'repositories' key found in JSON")
	}

	repos, ok := reposInterface.([]interface{})
	if !ok {
		return nil, fmt.Errorf("'repositories' is not an array")
	}

	result := make([]map[string]interface{}, len(repos))

	for i, repoInterface := range repos {
		repo, ok := repoInterface.(map[string]interface{})
		if !ok {
			continue
		}

		processed := ""
		on_env := ""

		repoName, ok := repo["repository_name"].(string)
		if !ok {
			continue
		}

		deploymentPath := filepath.Join("projects/projects-summary", repoName+".json")
		if fileInfo, err := os.Stat(deploymentPath); err == nil && fileInfo.Size() > 0 {
			processed = "true"
			deploymentData, err := os.ReadFile(deploymentPath)
			if err == nil && len(deploymentData) > 0 {
				var deploymentInfo map[string]interface{}
				if err := json.Unmarshal(deploymentData, &deploymentInfo); err == nil {
					if apps, ok := deploymentInfo["apps"].([]interface{}); ok && len(apps) > 0 {
						on_env = "true"
					}
				}
			}
		}

		newEntry := map[string]interface{}{
			"repository_name": repoName,
			"team":            repo["team"],
			"description":     repo["description"],
		}

		if on_env == "true" {
			newEntry["deployed"] = on_env
		}

		if processed == "true" {
			newEntry["processed"] = processed
		}

		result[i] = newEntry
	}

	response := map[string]interface{}{
		"repositories": result,
		"total_count":  len(result),
	}

	return json.Marshal(response)
}

// Terraform configuration processing
type terraformConfig struct{}

var terraformConfigs = new(terraformConfig)

func (t *terraformConfig) ParseConfigs(path string) ([]interface{}, error) {
	configs := []interface{}{
		map[string]interface{}{
			"path":            path,
			"region_default":  "us-west-2",
			"account_default": "prod",
			"namespace":       "default",
		},
	}
	return configs, nil
}

func (t *terraformConfig) ToJSON(configs []interface{}) (string, error) {
	data, err := json.MarshalIndent(configs, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func processTerrformConfigs(path string) ([]byte, error) {
	configs, err := terraformConfigs.ParseConfigs(path)
	if err != nil {
		return nil, fmt.Errorf("error parsing terraform configs: %v", err)
	}

	jsonStr, err := terraformConfigs.ToJSON(configs)
	if err != nil {
		return nil, fmt.Errorf("error converting to JSON: %v", err)
	}

	return []byte(jsonStr), nil
}

func runCLI(baseRepoName string) {
	if baseRepoName == "" {
		fmt.Println("Usage: go run main.go -repo=<repoName>")
		return
	}

	repoBitUrl, namespace, appNameSuffixes := getRepoFileDetails(baseRepoName)
	if repoBitUrl == "" {
		fmt.Println("Unknown baseRepoName")
		return
	}

	_, err := processRepoData(baseRepoName, repoBitUrl, namespace, appNameSuffixes, true)
	if err != nil {
		fmt.Printf("Error processing repo data: %v\n", err)
	}
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
