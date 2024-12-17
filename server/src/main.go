package main

import (
	"argocd/pkg/analyzer"
	"argocd/pkg/analyzerArgoCd"
	"argocd/pkg/regions"
	"argocd/pkg/terraformConfig"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type RepoDetails struct {
	RepoBitUrl      string
	Namespace       string
	AppNameSuffixes map[string]bool
}

var repoDetailsArray = []struct {
	BaseRepoName string
	Details      RepoDetails
}{}

type RegionDetails struct {
	Path           string `json:"path"`
	RegionDefault  string `json:"region_default"`
	AccountDefault string `json:"account_default"`
	Namespace      string `json:"namespace"`
}

type ReposListCache struct {
	Data      []byte
	Timestamp time.Time
}

var (
	reposListCache ReposListCache
	reposListMux   sync.RWMutex
	cacheDuration  = 1 * time.Minute
)

type DeploymentVersions struct {
	Stable string
	Canary string
}

func handleTerraformConfigs(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w, r)

	// Handle preflight requests
	if r.Method == http.MethodOptions {
		return
	}

	baseRepoName := r.URL.Query().Get("repo")
	if baseRepoName == "" {
		http.Error(w, "Missing repo parameter", http.StatusBadRequest)
		return
	}

	// Construct path to repo's terraform configs
	configPath := filepath.Join("projects/projects", baseRepoName, "scripts/terraform")

	configs, err := terraformConfig.ParseConfigs(configPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing terraform configs: %v", err), http.StatusInternalServerError)
		return
	}

	jsonData, err := terraformConfig.ToJSON(configs)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error converting to JSON: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonData))
}

func getRepoFileDetails(baseRepoName string) (string, string, map[string]bool) {
	// Construct the path to the regions.json file
	regionsFilePath := filepath.Join("projects/projects", baseRepoName, "regions.json")

	// Get the region details directly using the parser
	regions, err := regions.ParseRegions(baseRepoName)
	if err != nil {
		fmt.Printf("Error parsing regions: %v\n", err)
		return "", "", nil
	}

	// Convert to JSON and write to file
	jsonData, err := json.MarshalIndent(regions, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling regions to JSON: %v\n", err)
		return "", "", nil
	}

	// Write to file
	err = ioutil.WriteFile(regionsFilePath, jsonData, 0644)
	if err != nil {
		fmt.Printf("Error writing regions.json file: %v\n", err)
		return "", "", nil
	}

	// Populate the data structure
	repoBitUrl := fmt.Sprintf("https://github.com/pismo/%s", baseRepoName)
	appNameSuffixes := make(map[string]bool)
	namespace := ""

	// Use the first found namespace
	for _, region := range regions {
		suffix := fmt.Sprintf("-%s-%s", region.AccountDefault, region.RegionDefault)
		namespace = region.Namespace
		appNameSuffixes[suffix] = true
	}

	return repoBitUrl, namespace, appNameSuffixes
}

func extractImageFromJSON(jsonData string) (string, error) {
	// Define a regular expression to find the image field
	cleanedJSON := strings.ReplaceAll(jsonData, "\\", "")
	re := regexp.MustCompile(`"image":"([^"]+)"`)
	matches := re.FindStringSubmatch(cleanedJSON)
	if len(matches) < 2 {
		return "", fmt.Errorf("no image found in JSON")
	}
	return matches[1], nil
}

func getURLContent(url, token string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: received status code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Replace all occurrences of backslashes
	//cleanedBody := strings.ReplaceAll(string(body), "\\", "")
	return []byte(body), nil

	//return []byte(cleanedBody), nil
}

func writeFormattedJSONToFile(filename string, jsonData []byte) error {
	formattedData := pretty.Pretty(jsonData)

	err := ioutil.WriteFile(filename, formattedData, 0644)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	fmt.Printf("Writing: %s\n", filename)

	return nil
}

func fetchImages(baseRepoName, appNameSuffix, spr string, isPrimary bool) map[string]interface{} {

	appName := baseRepoName + appNameSuffix
	escapedAppName := url.QueryEscape(appName)
	url1 := fmt.Sprintf("https://argocd.pismo.services/api/v1/applications/%s/resource-tree?appNamespace=argocd", escapedAppName)
	url2 := fmt.Sprintf("https://argocd.pismo.services/api/v1/applications/%s/resource?name=%s&appNamespace=argocd&namespace=%s&resourceName=%s&version=v1alpha1&kind=Rollout&group=argoproj.io", escapedAppName, baseRepoName, spr, baseRepoName)

	tokenBytes, err := ioutil.ReadFile("token.txt")
	if err != nil {
		log.Fatalf("Error reading token file: %v", err)
	}
	token := strings.TrimSpace(string(tokenBytes))

	fmt.Printf("baseRepoName: %s, appNameSuffix: %s, spr: %s, isPrimary: %t\n", baseRepoName, appNameSuffix, spr, isPrimary)

	app := map[string]interface{}{"appName": appName, "type": "primary"}
	if !isPrimary {
		app["type"] = "failover"
	}
	errors := []string{}
	warnings := []string{}

	body1, err := getURLContent(url1, token)
	if err != nil {
		errors = append(errors, err.Error())
		app["error"] = errors
		return app
	}

	err = writeFormattedJSONToFile("tmp/"+appName+"-url1.json", body1)
	//if err != nil {
	//	fmt.Println("Error writing to file:", err)
	//	return map[string]interface{}{"appName": appName, "error": []string{err.Error()}}
	//}

	body2, err2 := getURLContent(url2, token)
	if err2 != nil {
		errors = append(errors, err2.Error())
		app["error"] = errors
		return app
	}

	err = writeFormattedJSONToFile("tmp/"+appName+"-url2.json", body2)
	//if err != nil {
	//	fmt.Println("Error writing to file:", err)
	//	return map[string]interface{}{"appName": appName, "error": []string{err.Error()}}
	//}

	// Process the first URL response
	jsonString1 := string(body1)
	phase := gjson.Get(jsonString1, "phase").String()
	if phase == "Error" {
		message := gjson.Get(jsonString1, "message").String()
		errors = append(errors, message)
	}

	// Process the second URL response
	jsonString2 := string(body2)
	if strings.Contains(jsonString2, "Error") {
		re := regexp.MustCompile(`"phase":"Error","message":"([^"]+)"`)
		match := re.FindStringSubmatch(jsonString2)
		if len(match) > 1 {
			errors = append(errors, "Error found in application response: "+match[1])
		} else {
			errors = append(errors, "Error found in application response, but no specific message could be extracted.")
		}
	}

	// Extract the image field from the JSON string
	image, err := extractImageFromJSON(jsonString2)
	if err != nil {
		errors = append(errors, err.Error())
	}

	var imageList []string
	printedImages := make(map[string]bool)

	// Add the extracted image to the list if it matches the baseRepoName
	if image != "" {
		if strings.Contains(image, baseRepoName) && !printedImages[image] {
			imageList = append(imageList, image)
			printedImages[image] = true
		}
	}

	if len(imageList) == 0 {
		errors = append(errors, "No images found")
	}
	app["images"] = imageList

	// Add warning if PR is found in any image name
	for _, image := range imageList {
		if strings.Contains(image, "PR-") && !strings.Contains(appNameSuffix, "integration") && !strings.Contains(appNameSuffix, "ext") {
			warnings = append(warnings, "PR found on non-integration/ext environment")
			break
		}
	}

	// Check if non-ext or non-integration images are the same
	if !strings.Contains(appNameSuffix, "integration") && !strings.Contains(appNameSuffix, "ext") {
		var referenceImage string
		for _, image := range imageList {
			if referenceImage == "" {
				referenceImage = image
			} else if image != referenceImage {
				errors = append(errors, "Non-ext or non-integration images are different")
				break
			}
		}
	}

	phase2 := gjson.Get(jsonString2, "operationState.phase").String()
	if phase2 == "Error" {
		message2 := gjson.Get(jsonString2, "operationState.message").String()
		errors = append(errors, message2)
	}

	// Add errors and warnings to the app map if they exist
	if len(errors) > 0 {
		app["error"] = errors
	}
	if len(warnings) > 0 {
		app["warning"] = warnings
	}

	// Analyze deployment and add the result to the app map
	result1, err := analyzer.AnalyzeDeployment(string(body1))
	if err != nil {
		log.Fatalf("Error analyzing deployment: %v", err)
	}

	var deployment interface{}
	err = json.Unmarshal([]byte(result1), &deployment)
	if err != nil {
		log.Fatalf("Error unmarshalling deployment: %v", err)
	}

	app["deployment"] = deployment

	result2, err := analyzerArgoCd.AnalyzeArgoCd(string(body2))
	if err != nil {
		log.Fatalf("Error analyzing deployment: %v", err)
	}

	var result2Map map[string]interface{}
	err = json.Unmarshal([]byte(result2), &result2Map)
	if err != nil {
		log.Fatalf("Error unmarshalling result2: %v", err)
	}

	health := ""
	nodesHealthStatus := gjson.Get(jsonString2, "health.status").String()
	if nodesHealthStatus == "Error" {
		health = gjson.Get(jsonString2, "health.status").String()
	}

	app["argocd"] = map[string]interface{}{
		"url":    "https://argocd.pismo.services/applications/argocd/" + appName + "?view=tree&orphaned=false&resource=",
		"status": result2Map,
		"health": health,
	}

	app["grafana"] = map[string]string{
		"url": "https://pismo.grafana.net/explore?schemaVersion=1&panes=%7B%22jqe%22%3A%7B%22datasource%22%3A%22grafanacloud-logs%22%2C%22queries%22%3A%5B%7B%22refId%22%3A%22A%22%2C%22expr%22%3A%22%7Bcontainer%3D%5C%22" + baseRepoName + "%5C%22%2C+env%3D%5C%22prod%5C%22%2C+version%3D%5C%221.24.0%5C%22%7D%22%2C%22queryType%22%3A%22range%22%2C%22datasource%22%3A%7B%22type%22%3A%22loki%22%2C%22uid%22%3A%22grafanacloud-logs%22%7D%2C%22editorMode%22%3A%22code%22%7D%5D%2C%22range%22%3A%7B%22from%22%3A%22now-7d%22%2C%22to%22%3A%22now%22%7D%7D%7D&orgId=1",
	}

	app["codefresh"] = map[string]string{
		"url": "https://pismo.grafana.net/explore?schemaVersion=1&panes=%7B%22jqe%22%3A%7B%22datasource%22%3A%22grafanacloud-logs%22%2C%22queries%22%3A%5B%7B%22refId%22%3A%22A%22%2C%22expr%22%3A%22%7Bcontainer%3D%5C%22" + baseRepoName + "%5C%22%2C+env%3D%5C%22prod%5C%22%2C+version%3D%5C%221.24.0%5C%22%7D%22%2C%22queryType%22%3A%22range%22%2C%22datasource%22%3A%7B%22type%22%3A%22loki%22%2C%22uid%22%3A%22grafanacloud-logs%22%7D%2C%22editorMode%22%3A%22code%22%7D%5D%2C%22range%22%3A%7B%22from%22%3A%22now-7d%22%2C%22to%22%3A%22now%22%7D%7D%7D&orgId=1",
	}

	return app
}

func getRepoDetails(baseRepoName string) (string, string, map[string]bool) {
	for _, repo := range repoDetailsArray {
		if repo.BaseRepoName == baseRepoName {
			return repo.Details.RepoBitUrl, repo.Details.Namespace, repo.Details.AppNameSuffixes
		}
	}
	fmt.Println("Unknown baseRepoName")
	return "", "", nil
}

func processRepoData(baseRepoName, repoBitUrl, namespace string, appNameSuffixes map[string]bool, forceRefresh bool) ([]byte, error) {
	repoName := baseRepoName

	// Ensure the projects/summary directory exists
	summaryDir := "projects/projects-summary"
	if _, err := os.Stat(summaryDir); os.IsNotExist(err) {
		err := os.MkdirAll(summaryDir, os.ModePerm)
		if err != nil {
			log.Fatalf("Error creating directory: %v", err)
		}
	}

	// Construct the file path in the projects/summary directory
	filename := filepath.Join(summaryDir, fmt.Sprintf("%s.json", repoName))

	cacheTime := 30
	if forceRefresh {
		cacheTime = 1
	}

	// Check if the file exists and is less than cacheTime seconds old
	fileInfo, err := os.Stat(filename)
	if err == nil && time.Since(fileInfo.ModTime()) < time.Duration(cacheTime)*time.Second {
		// Read the data from the file
		jsonData, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf("error reading from file: %w", err)
		}
		fmt.Println("JSON data read from", filename)
		return jsonData, nil
	}

	// Process the data if the file does not exist or is older than cacheTime seconds
	repo, err := getRepositoryBlock(repoName)
	if err != nil {
		fmt.Println(err)
		//return
	}

	repoData := map[string]interface{}{
		"repoName":      baseRepoName,
		"repoBitUrl":    repoBitUrl,
		"repoCodefresh": "https://g.codefresh.io/pipelines/all/?filter=pageSize:10;field:name~Name;order:asc~Asc;search:" + baseRepoName,
		"apps":          []map[string]interface{}{},
		"argocd": map[string]string{
			"url": "https://argocd.pismo.services/applications?search=" + baseRepoName + "&showFavorites=false&proj=&sync=&autoSync=&health=&namespace=&cluster=&labels=",
		},
		"repoDesc":  repo.Description,
		"repoSquad": repo.Team,
	}

	// Limit the number of concurrent goroutines to 5 by using a semaphore pattern with a buffered channel.
	var wg sync.WaitGroup
	var mu sync.Mutex
	sem := make(chan struct{}, 10) // Semaphore with a capacity of 1

	for appNameSuffix, isPrimary := range appNameSuffixes {
		wg.Add(1)
		sem <- struct{}{} // Acquire a slot
		go func(appNameSuffix string, isPrimary bool) {
			defer wg.Done()
			defer func() { <-sem }() // Release the slot

			app := fetchImages(baseRepoName, appNameSuffix, namespace, isPrimary)
			// a random 1 to 3 second pause
			time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
			mu.Lock()
			repoData["apps"] = append(repoData["apps"].([]map[string]interface{}), app)
			mu.Unlock()
		}(appNameSuffix, isPrimary)
	}

	wg.Wait()

	jsonData, err := json.MarshalIndent(repoData, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error marshalling to JSON: %w", err)
	}

	err = ioutil.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return nil, fmt.Errorf("error writing to file: %w", err)
	}

	fmt.Println("JSON data written to", filename)
	return jsonData, nil
}

func setCORSHeaders(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if strings.HasPrefix(origin, "http://localhost:") {
		portStr := strings.TrimPrefix(origin, "http://localhost:")
		port, err := strconv.Atoi(portStr)
		if err == nil && port >= 3000 && port <= 3100 {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
	}
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func listReposHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w, r)

	// Handle preflight requests
	if r.Method == http.MethodOptions {
		return
	}

	// Extract the list of BaseRepoNames
	var repoNames []string
	for _, repo := range repoDetailsArray {
		repoNames = append(repoNames, repo.BaseRepoName)
	}

	// Convert the list to JSON
	jsonData, err := json.Marshal(repoNames)
	if err != nil {
		http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
		return
	}

	// Set the content type and write the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// checkAndPullRepo checks if the GitHub folder exists and pulls it if it doesn't.
func checkAndPullRepo(baseRepoName string) error {
	repoPath := filepath.Join("projects/projects", baseRepoName, "github")
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		fmt.Printf("Repository %s does not exist. Cloning...\n", baseRepoName)
		cloneURL := fmt.Sprintf("https://github.com/pismo/%s.git", baseRepoName)
		cmd := exec.Command("git", "clone", cloneURL, repoPath)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("error cloning repository: %v", err)
		}
		fmt.Printf("Repository %s cloned successfully.\n", baseRepoName)
	}
	return nil
}

func handleRepoRequest(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w, r)

	// Handle preflight requests
	if r.Method == http.MethodOptions {
		return
	}

	// Check for force refresh parameter
	forceRefresh := r.URL.Query().Get("force") != ""

	fmt.Println(forceRefresh)

	baseRepoName := r.URL.Query().Get("repo")
	if baseRepoName == "" {
		http.Error(w, "Missing repo parameter", http.StatusBadRequest)
		return
	}

	if err := checkAndPullRepo(baseRepoName); err != nil {
		log.Fatalf("Error checking and pulling repository: %v", err)
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

func getEnvironmentVersions(deploymentInfo map[string]interface{}) map[string]string {
	envVersions := make(map[string]string)

	if apps, ok := deploymentInfo["apps"].([]interface{}); ok {
		for _, app := range apps {
			if appMap, ok := app.(map[string]interface{}); ok {
				appName := appMap["appName"].(string)

				// Skip if there's an error in the app data
				if _, hasError := appMap["error"]; hasError {
					continue
				}

				// Extract environment path from app name
				parts := strings.Split(appName, "-")
				if len(parts) <= 1 {
					continue
				}

				// Remove the microservice name prefix
				envPath := strings.TrimPrefix(appName, parts[0]+"-"+parts[1]+"-"+parts[2]+"-")
				if envPath == "" {
					continue
				}

				// Get deployment info if it exists
				if deployment, ok := appMap["deployment"].(map[string]interface{}); ok {
					if deployments, ok := deployment["deployments"].([]interface{}); ok && len(deployments) > 0 {
						// Get the first deployment's version
						if deploy, ok := deployments[0].(map[string]interface{}); ok {
							if version, ok := deploy["version"].(string); ok {
								envKey := "env-" + envPath
								envVersions[envKey] = version
							}
						}
					}
				}
			}
		}
	}

	return envVersions
}

func listReposFromFileHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w, r)

	// Handle preflight requests
	if r.Method == http.MethodOptions {
		return
	}

	// Check for force refresh parameter
	forceRefresh := r.URL.Query().Get("force") != ""
	
	// Check cache if not forcing refresh
	if !forceRefresh {
		reposListMux.RLock()
		if !reposListCache.Timestamp.IsZero() && time.Since(reposListCache.Timestamp) < cacheDuration {
			reposListMux.RUnlock()
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Cache", "HIT")
			w.Write(reposListCache.Data)
			return
		}
		reposListMux.RUnlock()
	}

	fmt.Println("Full refresh")
	// Read the content of pismo.json
	jsonData, err := ioutil.ReadFile("projects/projects/pismo.json")
	if err != nil {
		log.Printf("Error reading file: %v", err)
		http.Error(w, "Error reading pismo.json", http.StatusInternalServerError)
		return
	}

	// Parse the JSON into a map
	var data map[string]interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		log.Printf("JSON Unmarshal error: %v", err)
		http.Error(w, "Error parsing pismo.json", http.StatusInternalServerError)
		return
	}

	// Check if repositories exists and is an array
	reposInterface, ok := data["repositories"]
	if !ok {
		log.Printf("No 'repositories' key found in JSON")
		http.Error(w, "Invalid JSON structure", http.StatusInternalServerError)
		return
	}

	repos, ok := reposInterface.([]interface{})
	if !ok {
		log.Printf("'repositories' is not an array")
		http.Error(w, "Invalid JSON structure", http.StatusInternalServerError)
		return
	}

	// Convert repositories to the desired format
	result := make([]map[string]interface{}, len(repos))

	for i, repoInterface := range repos {
		repo, ok := repoInterface.(map[string]interface{})
		if !ok {
			log.Printf("Repository at index %d is not an object", i)
			continue
		}

		processed := ""
		on_env := ""

		repoName, ok := repo["repository_name"].(string)
		if !ok {
			log.Printf("Repository at index %d has no valid repository_name", i)
			continue
		}

		deploymentPath := filepath.Join("projects/projects-summary", repoName+".json")
		if fileInfo, err := os.Stat(deploymentPath); err == nil && fileInfo.Size() > 0 {
			processed = "true"
			deploymentData, err := ioutil.ReadFile(deploymentPath)
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

		if fileInfo, err := os.Stat(deploymentPath); err == nil && fileInfo.Size() > 0 {
			processed = "true"
			deploymentData, err := ioutil.ReadFile(deploymentPath)
			if err == nil && len(deploymentData) > 0 {
				var deploymentInfo map[string]interface{}
				if err := json.Unmarshal(deploymentData, &deploymentInfo); err == nil {
					if apps, ok := deploymentInfo["apps"].([]interface{}); ok && len(apps) > 0 {
						on_env = "true"
						// Get environment versions
						envVersions := getEnvironmentVersions(deploymentInfo)
						// Add environment versions to newEntry
						for env, version := range envVersions {
							newEntry[env] = version
						}
					}
				}
			}
		}

		result[i] = newEntry
	}

	response := map[string]interface{}{
		"repositories": result,
	}

	// Convert to JSON
	updatedJSON, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}

	// Update cache
	reposListMux.Lock()
	reposListCache = ReposListCache{
		Data:      updatedJSON,
		Timestamp: time.Now(),
	}
	reposListMux.Unlock()

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Cache", "MISS")
	w.Write(updatedJSON)
}

func main() {
	baseRepoNamePtr := flag.String("repo", "", "The base repository name")
	webserverPtr := flag.Bool("webserver", false, "Run as a webserver")
	flag.Parse()
	baseRepoName := *baseRepoNamePtr
	webserver := *webserverPtr

	if webserver {
		http.HandleFunc("/", handleRepoRequest)
		http.HandleFunc("/repos", listReposHandler)
		http.HandleFunc("/list-repos", listReposFromFileHandler)

		fmt.Println("Starting web server on :8083")
		if err := http.ListenAndServe(":8083", nil); err != nil {
			fmt.Println("Error starting web server:", err)
		}
		return
	}

	if err := checkAndPullRepo(baseRepoName); err != nil {
		log.Fatalf("Error checking and pulling repository: %v", err)
	}

	repoBitUrl, namespace, appNameSuffixes := getRepoFileDetails(baseRepoName)
	if repoBitUrl == "" {
		fmt.Println("Usage: go run main.go -repo=<repoName>")
		fmt.Println("Available baseRepoNames:")
		for _, repo := range repoDetailsArray {
			fmt.Printf("  - %s\n", repo.BaseRepoName)
		}
		return
	}

	_, err := processRepoData(baseRepoName, repoBitUrl, namespace, appNameSuffixes, true)
	if err != nil {
		fmt.Println(err)
	}
}
