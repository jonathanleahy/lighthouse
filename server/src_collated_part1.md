# Project Files Summary - Part 1

Generated on: 2024-12-24T12:26:34Z

Root Directory: ./src

---

## File: go.mod

Size: 195 bytes

Last Modified: 2024-12-15T09:35:10Z

```
module argocd

go 1.23.1

require (
	github.com/buger/jsonparser v1.1.1
	github.com/tidwall/gjson v1.18.0
	github.com/tidwall/pretty v1.2.1
)

require github.com/tidwall/match v1.1.1 // indirect

```

## File: go.sum

Size: 767 bytes

Last Modified: 2024-12-15T09:35:10Z

```
github.com/buger/jsonparser v1.1.1 h1:2PnMjfWD7wBILjqQbt530v576A/cAbQvEW9gGIpYMUs=
github.com/buger/jsonparser v1.1.1/go.mod h1:6RYKKt7H4d4+iWqouImQ9R2FZql3VbhNgx27UK13J/0=
github.com/tidwall/gjson v1.18.0 h1:FIDeeyB800efLX89e5a8Y0BNH+LOngJyGrIWxG2FKQY=
github.com/tidwall/gjson v1.18.0/go.mod h1:/wbyibRr2FHMks5tjHJ5F8dMZh3AcwJEMf5vlfC0lxk=
github.com/tidwall/match v1.1.1 h1:+Ho715JplO36QYgwN9PGYNhgZvoUSc9X2c80KVTi+GA=
github.com/tidwall/match v1.1.1/go.mod h1:eRSPERbgtNPcGhD8UCthc6PmLEQXEWd3PRB5JTxsfmM=
github.com/tidwall/pretty v1.2.0/go.mod h1:ITEVvHYasfjBbM0u2Pg8T2nJnzm8xPwvNhhsoaGGjNU=
github.com/tidwall/pretty v1.2.1 h1:qjsOFOWWQl+N3RsoF5/ssm1pHmJJwhjlSbZ51I6wMl4=
github.com/tidwall/pretty v1.2.1/go.mod h1:ITEVvHYasfjBbM0u2Pg8T2nJnzm8xPwvNhhsoaGGjNU=


```

## File: main.go

Size: 24689 bytes

Last Modified: 2024-12-19T22:02:34Z

```go
package main

import (
	"argocd/pkg/analyzer"
	"argocd/pkg/analyzerArgoCd"
	"argocd/pkg/gitProcessor"

	//"argocd/pkg/gitParser/pkg/gitProcessor"
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

	cacheTime := 30000000
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

	//// Process the data if the file does not exist or is older than cacheTime seconds
	repo, err := getRepositoryBlock(repoName)
	if err != nil {
		fmt.Println(err)
		//return
	}

	//// Call the gitProcessor code and then add those details into a new GitHub section
	options := gitProcessor.Options{
		CommitHistoryMonths:  24, // Example: 24 months of commit history
		ReleaseHistoryMonths: 24, // Example: 24 months of release history
	}
	repoModule, err := gitProcessor.NewRepositoryModule(options)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize repository module: %v", err)
	}

	repoPath := filepath.Join("projects/projects", baseRepoName, "github")
	repoDetails, err := repoModule.Extract(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to extract repository details: %v", err)
	}

	var gitDetails gitProcessor.AnalysisResult
	if err := json.Unmarshal(repoDetails, &gitDetails); err != nil {
		return nil, fmt.Errorf("failed to unmarshal repository details: %v", err)
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
		"github": map[string]interface{}{
			"commit":    gitDetails.Commit,
			"branch":    gitDetails.Branch,
			"author":    gitDetails.Author,
			"timestamp": gitDetails.Timestamp,
			"tag":       gitDetails.Tag,
		},
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
			//return fmt.Errorf("error cloning repository: %v", err)
			fmt.Printf("Error cloning Repository %s .\n", baseRepoName)
		} else {
			fmt.Printf("Repository %s cloned successfully.\n", baseRepoName)
		}
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
	forceRefresh := r.URL.Query().Get("force") == "true"

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
	forceRefresh := r.URL.Query().Get("force") == "true"

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

	// call the github and add a new section call github with the same structure as the others

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

```

## File: pismo-json-processor.go

Size: 1023 bytes

Last Modified: 2024-12-21T22:31:58Z

```go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Repository struct {
	RepositoryName string `json:"repository_name"`
	Team           string `json:"team"`
	Description    string `json:"description"`
}

type PismoData struct {
	Repositories []Repository `json:"repositories"`
	TotalCount   int          `json:"total_count"`
}

func getRepositoryBlock(repoName string) (*Repository, error) {
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

```

## Directory: pkg

## Directory: pkg/analyzer

## File: pkg/analyzer/analyzer.go

Size: 2723 bytes

Last Modified: 2024-12-15T09:35:10Z

```go
// analyzer/analyzer.go
package analyzer

import (
	"encoding/json"
	"fmt"
)

func AnalyzeDeployment(jsonData string) (string, error) {
	var data K8sData
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling input: %v", err)
	}

	// Track versions and their deployment stats
	versions := make(map[string]*struct {
		Type      string
		PodCount  int
		NodeNames map[string]bool
	})

	// First pass: Collect all versions
	for _, node := range data.Nodes {
		if node.Kind == "Pod" {
			labels := node.NetworkingInfo.Labels
			version := labels["version"]
			if version == "" {
				continue
			}

			// Initialize version tracking if needed
			if _, exists := versions[version]; !exists {
				versions[version] = &struct {
					Type      string
					PodCount  int
					NodeNames map[string]bool
				}{
					Type:      "uncategorized",
					NodeNames: make(map[string]bool),
				}
			}

			// Set type based on Argo Rollout labels
			if labels["argoproj.io/version"] == "canary" {
				versions[version].Type = "canary"
			} else if labels["argoproj.io/version"] == "stable" {
				versions[version].Type = "stable"
			}
		}
	}

	// Second pass: Count pods and nodes
	for _, node := range data.Nodes {
		if node.Kind == "Pod" {
			labels := node.NetworkingInfo.Labels
			version := labels["version"]
			if version == "" {
				continue
			}

			// Find node name from Info
			var nodeName string
			for _, info := range node.Info {
				if info.Name == "Node" {
					nodeName = info.Value
					break
				}
			}

			if versionInfo, exists := versions[version]; exists {
				versionInfo.PodCount++
				if nodeName != "" {
					versionInfo.NodeNames[nodeName] = true
				}
			}
		}
	}

	// Calculate total pods
	totalPods := 0
	for _, info := range versions {
		totalPods += info.PodCount
	}

	// Create analysis result
	analysis := DeploymentAnalysis{
		TotalPods:   totalPods,
		Deployments: make([]VersionDeployment, 0, len(versions)),
	}

	// Convert map to slice and calculate percentages
	for version, info := range versions {
		nodeNames := make([]string, 0, len(info.NodeNames))
		for nodeName := range info.NodeNames {
			nodeNames = append(nodeNames, nodeName)
		}

		deployment := VersionDeployment{
			Version:    version,
			Type:       info.Type,
			PodCount:   info.PodCount,
			NodeCount:  len(info.NodeNames),
			NodeNames:  nodeNames,
			Percentage: float64(info.PodCount) / float64(totalPods) * 100,
		}
		analysis.Deployments = append(analysis.Deployments, deployment)
	}

	// Convert to JSON
	result, err := json.MarshalIndent(analysis, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshaling result: %v", err)
	}

	return string(result), nil
}


```

## File: pkg/analyzer/types.go

Size: 1234 bytes

Last Modified: 2024-12-15T09:35:10Z

```go
// analyzer/types.go
package analyzer

type ResourceInfo struct {
	Capacity             int64  `json:"capacity"`
	RequestedByApp       int64  `json:"requestedByApp"`
	RequestedByNeighbors int64  `json:"requestedByNeighbors"`
	ResourceName         string `json:"resourceName"`
}

type Host struct {
	Name          string         `json:"name"`
	ResourcesInfo []ResourceInfo `json:"resourcesInfo"`
}

type NetworkingInfo struct {
	Labels map[string]string `json:"labels"`
}

type Pod struct {
	Name           string         `json:"name"`
	Kind           string         `json:"kind"`
	NetworkingInfo NetworkingInfo `json:"networkingInfo"`
	Info           []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"info"`
}

type K8sData struct {
	Hosts []Host `json:"hosts"`
	Nodes []Pod  `json:"nodes"`
}

type VersionDeployment struct {
	Version    string   `json:"version"`
	Type       string   `json:"type"`
	PodCount   int      `json:"podCount"`
	NodeCount  int      `json:"nodeCount"`
	NodeNames  []string `json:"nodeNames"`
	Percentage float64  `json:"percentage"`
}

type DeploymentAnalysis struct {
	TotalPods   int                 `json:"totalPods"`
	Deployments []VersionDeployment `json:"deployments"`
}


```

## Directory: pkg/analyzerArgoCd

## File: pkg/analyzerArgoCd/analyzerArgoCd.go

Size: 4292 bytes

Last Modified: 2024-12-15T09:35:10Z

```go
package analyzerArgoCd

import (
	"encoding/json"
	"fmt"
	"github.com/buger/jsonparser"
	"sort"
)

func AnalyzeArgoCd(jsonData string) (string, error) {
	// Extract the manifest field
	manifestStr, err := jsonparser.GetString([]byte(jsonData), "manifest")
	if err != nil {
		return "", fmt.Errorf("manifest field is missing or not a string: %v", err)
	}

	// Extract the spec field from the manifest
	spec, _, _, err := jsonparser.Get([]byte(manifestStr), "spec")
	if err != nil {
		return "", fmt.Errorf("spec field is missing or not a map: %v", err)
	}

	// Extract the strategy field from the spec
	strategy, _, _, err := jsonparser.Get(spec, "strategy")
	if err != nil {
		return "", fmt.Errorf("strategy field is missing or not a map: %v", err)
	}

	// Extract the canary field from the strategy
	canary, _, _, err := jsonparser.Get(strategy, "canary")
	if err != nil {
		return "", fmt.Errorf("canary field is missing or not a map: %v", err)
	}

	// Extract the status field from the manifest
	status, _, _, err := jsonparser.Get([]byte(manifestStr), "status")
	if err != nil {
		return "", fmt.Errorf("status field is missing or not a map: %v", err)
	}

	// Extract the steps field from the canary
	steps, _, _, err := jsonparser.Get(canary, "steps")
	if err != nil {
		return "", fmt.Errorf("steps field is missing or not a list: %v", err)
	}

	// Extract the currentStepIndex from the status
	currentStepIndex, err := jsonparser.GetInt(status, "currentStepIndex")
	if err != nil {
		return "", fmt.Errorf("currentStepIndex field is missing or not an integer: %v", err)
	}

	// Handle the case where currentStepIndex is -1
	if currentStepIndex == -1 {
		return "currentStepIndex is -1", nil
	}

	// Check if currentStepIndex is within the bounds of the steps array
	stepsCount := 0
	_, err = jsonparser.ArrayEach(steps, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		stepsCount++
	})
	if err != nil {
		return "", fmt.Errorf("error counting steps: %v", err)
	}

	if currentStepIndex < 0 || currentStepIndex >= int64(stepsCount) {
		//return "", fmt.Errorf("currentStepIndex %d is out of bounds", currentStepIndex)
	} else {

		// Find the latest step and the setWeight value from the previous step
		var latestStep, setWeight string
		index := 0

		// Sort steps by date-time order
		sortedSteps := make([][]byte, 0)
		_, err = jsonparser.ArrayEach(steps, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			sortedSteps = append(sortedSteps, value)
		})
		if err != nil {
			return "", fmt.Errorf("error iterating steps: %v", err)
		}

		// Sort the steps array by date-time order
		sort.Slice(sortedSteps, func(i, j int) bool {
			dateTimeI, _ := jsonparser.GetString(sortedSteps[i], "dateTime")
			dateTimeJ, _ := jsonparser.GetString(sortedSteps[j], "dateTime")
			return dateTimeI < dateTimeJ
		})

		// Iterate through the sorted steps
		for _, step := range sortedSteps {
			if index == int(currentStepIndex) {
				latestStep = string(step)
				break
			} else {
				_, dataType, _, err := jsonparser.Get(step, "setWeight")
				if err != nil {
					//return "", fmt.Errorf("error getting setWeight: %v", err)
				} else {
					switch dataType {
					case jsonparser.String:
						setWeight, err = jsonparser.GetString(step, "setWeight")
						if err != nil {
							return "", fmt.Errorf("error getting setWeight: %v", err)
						}
					case jsonparser.Number:
						setWeightFloat, err := jsonparser.GetFloat(step, "setWeight")
						if err != nil {
							return "", fmt.Errorf("error getting setWeight: %v", err)
						}
						setWeight = fmt.Sprintf("%f", setWeightFloat)
					default:
						return "", fmt.Errorf("setWeight is not a valid type: %v", dataType)
					}
				}
			}
			index++
		}

		if latestStep == "" {
			return "", fmt.Errorf("no step found for currentStepIndex: %d", currentStepIndex)
		}

		// Create the result map
		resultWithIndex := map[string]interface{}{
			"step": []string{latestStep},
			//"index":  currentStepIndex,
			"weight": setWeight,
		}

		// Convert the result to JSON
		resultJSON, err := json.Marshal(resultWithIndex)
		if err != nil {
			return "", fmt.Errorf("error marshalling result to JSON: %v", err)
		}

		return string(resultJSON), nil
	}
	return string("{}"), nil

}


```

## Directory: pkg/gitParser

## Directory: pkg/gitParser/.idea

## File: pkg/gitParser/.idea/.gitignore

Size: 176 bytes

Last Modified: 2024-12-19T19:41:49Z

```
# Default ignored files
/shelf/
/workspace.xml
# Editor-based HTTP Client requests
/httpRequests/
# Datasource local storage ignored files
/dataSources/
/dataSources.local.xml

```

## File: pkg/gitParser/.idea/git-parser.iml

Size: 322 bytes

Last Modified: 2024-12-19T19:41:49Z

```
<?xml version="1.0" encoding="UTF-8"?>
<module type="WEB_MODULE" version="4">
  <component name="Go" enabled="true" />
  <component name="NewModuleRootManager">
    <content url="file://$MODULE_DIR$" />
    <orderEntry type="inheritedJdk" />
    <orderEntry type="sourceFolder" forTests="false" />
  </component>
</module>
```

## File: pkg/gitParser/.idea/modules.xml

Size: 272 bytes

Last Modified: 2024-12-19T19:41:49Z

```xml
<?xml version="1.0" encoding="UTF-8"?>
<project version="4">
  <component name="ProjectModuleManager">
    <modules>
      <module fileurl="file://$PROJECT_DIR$/.idea/git-parser.iml" filepath="$PROJECT_DIR$/.idea/git-parser.iml" />
    </modules>
  </component>
</project>
```

## File: pkg/gitParser/.idea/vcs.xml

Size: 267 bytes

Last Modified: 2024-12-19T20:23:41Z

```xml
<?xml version="1.0" encoding="UTF-8"?>
<project version="4">
  <component name="VcsDirectoryMappings">
    <mapping directory="$PROJECT_DIR$/../../../.." vcs="Git" />
    <mapping directory="$PROJECT_DIR$/tmp/anthropic-astroids" vcs="Git" />
  </component>
</project>
```

## File: pkg/gitParser/.idea/workspace.xml

Size: 7458 bytes

Last Modified: 2024-12-19T21:15:33Z

```xml
<?xml version="1.0" encoding="UTF-8"?>
<project version="4">
  <component name="AutoImportSettings">
    <option name="autoReloadType" value="ALL" />
  </component>
  <component name="ChangeListManager">
    <list default="true" id="ccfcde17-8801-482f-b335-158d9d9d3ff9" name="Changes" comment="">
      <change afterPath="$PROJECT_DIR$/pkg/gitProcessor/gitProcessor.go1" afterDir="false" />
      <change afterPath="$PROJECT_DIR$/pkg/gitProcessor/gitResults.go1" afterDir="false" />
      <change afterPath="$PROJECT_DIR$/pkg/gitProcessor/types.go" afterDir="false" />
      <change beforePath="$PROJECT_DIR$/../../../../dashboard/package-lock.json" beforeDir="false" />
      <change beforePath="$PROJECT_DIR$/../../../../dashboard/src/app/customize-fields/page.tsx" beforeDir="false" afterPath="$PROJECT_DIR$/../../../../dashboard/src/app/customize-fields/page.tsx" afterDir="false" />
      <change beforePath="$PROJECT_DIR$/../../../../dashboard/src/app/globals.css" beforeDir="false" afterPath="$PROJECT_DIR$/../../../../dashboard/src/app/globals.css" afterDir="false" />
      <change beforePath="$PROJECT_DIR$/../../../../dashboard/src/app/microservice/[id]/page.tsx" beforeDir="false" afterPath="$PROJECT_DIR$/../../../../dashboard/src/app/microservice/[id]/page.tsx" afterDir="false" />
      <change beforePath="$PROJECT_DIR$/../../../../dashboard/src/app/page.tsx" beforeDir="false" afterPath="$PROJECT_DIR$/../../../../dashboard/src/app/page.tsx" afterDir="false" />
      <change beforePath="$PROJECT_DIR$/../../../../dashboard/src/components/dashboard-grid.tsx" beforeDir="false" afterPath="$PROJECT_DIR$/../../../../dashboard/src/components/dashboard-grid.tsx" afterDir="false" />
      <change beforePath="$PROJECT_DIR$/../../../../dashboard/src/components/dashboard.tsx" beforeDir="false" afterPath="$PROJECT_DIR$/../../../../dashboard/src/components/dashboard.tsx" afterDir="false" />
      <change beforePath="$PROJECT_DIR$/../../../../dashboard/src/components/deployments.tsx" beforeDir="false" afterPath="$PROJECT_DIR$/../../../../dashboard/src/components/deployments.tsx" afterDir="false" />
      <change beforePath="$PROJECT_DIR$/../../../projects/projects-summary/backoffice-core-bff.json" beforeDir="false" afterPath="$PROJECT_DIR$/../../../projects/projects-summary/backoffice-core-bff.json" afterDir="false" />
      <change beforePath="$PROJECT_DIR$/../../../projects/projects-summary/console-audit-bff.json" beforeDir="false" afterPath="$PROJECT_DIR$/../../../projects/projects-summary/console-audit-bff.json" afterDir="false" />
      <change beforePath="$PROJECT_DIR$/../../../server" beforeDir="false" afterPath="$PROJECT_DIR$/../../../server" afterDir="false" />
      <change beforePath="$PROJECT_DIR$/../../main.go" beforeDir="false" afterPath="$PROJECT_DIR$/../../main.go" afterDir="false" />
      <change beforePath="$PROJECT_DIR$/go.mod" beforeDir="false" afterPath="$PROJECT_DIR$/go.mod" afterDir="false" />
      <change beforePath="$PROJECT_DIR$/go.sum" beforeDir="false" afterPath="$PROJECT_DIR$/go.sum" afterDir="false" />
      <change beforePath="$PROJECT_DIR$/main.go" beforeDir="false" afterPath="$PROJECT_DIR$/main.go" afterDir="false" />
      <change beforePath="$PROJECT_DIR$/pkg/gitProcessor/gitProcessor.go" beforeDir="false" afterPath="$PROJECT_DIR$/pkg/gitProcessor/gitProcessor.go" afterDir="false" />
      <change beforePath="$PROJECT_DIR$/pkg/gitProcessor/gitResults.go" beforeDir="false" />
    </list>
    <option name="SHOW_DIALOG" value="false" />
    <option name="HIGHLIGHT_CONFLICTS" value="true" />
    <option name="HIGHLIGHT_NON_ACTIVE_CHANGELIST" value="false" />
    <option name="LAST_RESOLUTION" value="IGNORE" />
  </component>
  <component name="GOROOT" url="file:///snap/go/current" />
  <component name="Git.Settings">
    <option name="RECENT_GIT_ROOT_PATH" value="$PROJECT_DIR$/../../../.." />
  </component>
  <component name="GitHubPullRequestSearchHistory">{
  &quot;lastFilter&quot;: {
    &quot;state&quot;: &quot;OPEN&quot;,
    &quot;assignee&quot;: &quot;jonathanleahy&quot;
  }
}</component>
  <component name="GithubPullRequestsUISettings">{
  &quot;selectedUrlAndAccountId&quot;: {
    &quot;url&quot;: &quot;https://github.com/jonathanleahy/lighthouse.git&quot;,
    &quot;accountId&quot;: &quot;48ebb4ce-34fe-4984-a0aa-c40105418828&quot;
  }
}</component>
  <component name="GoLibraries">
    <option name="indexEntireGoPath" value="true" />
  </component>
  <component name="NamedScopeManager">
    <scope name="GitScopePro" pattern="" />
  </component>
  <component name="ProjectColorInfo">{
  &quot;associatedIndex&quot;: 8
}</component>
  <component name="ProjectId" id="2qRxFiQGGvUPuwca1mPESWIP4NX" />
  <component name="ProjectViewState">
    <option name="hideEmptyMiddlePackages" value="true" />
    <option name="showLibraryContents" value="true" />
  </component>
  <component name="PropertiesComponent">{
  &quot;keyToString&quot;: {
    &quot;Go Build.go build git-parser.executor&quot;: &quot;Run&quot;,
    &quot;Go Build.go build git-parser/.executor&quot;: &quot;Run&quot;,
    &quot;RunOnceActivity.ShowReadmeOnStart&quot;: &quot;true&quot;,
    &quot;RunOnceActivity.git.unshallow&quot;: &quot;true&quot;,
    &quot;RunOnceActivity.go.formatter.settings.were.checked&quot;: &quot;true&quot;,
    &quot;RunOnceActivity.go.migrated.go.modules.settings&quot;: &quot;true&quot;,
    &quot;RunOnceActivity.go.modules.go.list.on.any.changes.was.set&quot;: &quot;true&quot;,
    &quot;git-widget-placeholder&quot;: &quot;main&quot;,
    &quot;go.import.settings.migrated&quot;: &quot;true&quot;,
    &quot;go.sdk.automatically.set&quot;: &quot;true&quot;,
    &quot;last_opened_file_path&quot;: &quot;/home/jon/lighthouse/server/src/pkg/git-parser/tmp&quot;,
    &quot;node.js.detected.package.eslint&quot;: &quot;true&quot;,
    &quot;node.js.selected.package.eslint&quot;: &quot;(autodetect)&quot;,
    &quot;nodejs_package_manager_path&quot;: &quot;npm&quot;
  }
}</component>
  <component name="RecentsManager">
    <key name="CopyFile.RECENT_KEYS">
      <recent name="$PROJECT_DIR$/tmp" />
    </key>
  </component>
  <component name="RunManager">
    <configuration name="go build git-parser/" type="GoApplicationRunConfiguration" factoryName="Go Application" nameIsGenerated="true">
      <module name="git-parser" />
      <working_directory value="$PROJECT_DIR$" />
      <parameters value="-path tmp/anthropic-astroids " />
      <kind value="DIRECTORY" />
      <package value="git-parser" />
      <directory value="$PROJECT_DIR$" />
      <filePath value="$PROJECT_DIR$" />
      <method v="2" />
    </configuration>
  </component>
  <component name="SharedIndexes">
    <attachedChunks>
      <set>
        <option value="bundled-gosdk-d297c17c1fbd-b87a2f8923ed-org.jetbrains.plugins.go.sharedIndexes.bundled-GO-243.21565.208" />
        <option value="bundled-js-predefined-d6986cc7102b-e768b9ed790e-JavaScript-GO-243.21565.208" />
      </set>
    </attachedChunks>
  </component>
  <component name="SpellCheckerSettings" RuntimeDictionaries="0" Folders="0" CustomDictionaries="0" DefaultDictionary="application-level" UseSingleDictionary="true" transferred="true" />
  <component name="TypeScriptGeneratedFilesManager">
    <option name="version" value="3" />
  </component>
  <component name="VgoProject">
    <integration-enabled>false</integration-enabled>
    <settings-migrated>true</settings-migrated>
  </component>
</project>
```

## File: pkg/gitParser/go.mod

Size: 29 bytes

Last Modified: 2024-12-19T20:12:32Z

```
module git-parser

go 1.23.3

```

## File: pkg/gitParser/go.sum

Size: 0 bytes

Last Modified: 2024-12-19T20:12:32Z

```

```

## File: pkg/gitParser/main.go

Size: 6561 bytes

Last Modified: 2024-12-19T21:35:34Z

```go
package main

import (
	gitProcessor2 "argocd/pkg/gitProcessor"
	"encoding/json"
	"flag"
	"fmt"
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
	repoModule, err := gitProcessor2.NewRepositoryModule(gitProcessor2.Options{
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
	var result gitProcessor2.AnalysisResult
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

func printSummary(result *gitProcessor2.AnalysisResult) {
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

```

## Directory: pkg/gitParser/pkg

## Directory: pkg/gitParser/schemas

## File: pkg/gitParser/schemas/repository.json

Size: 1196 bytes

Last Modified: 2024-12-15T09:35:10Z

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "properties": {
    "git": {
      "type": "object",
      "properties": {
        "lastCommit": {
          "type": "object",
          "properties": {
            "hash": {"type": "string"},
            "author": {"type": "string"},
            "date": {"type": "string"},
            "message": {"type": "string"}
          }
        },
        "branch": {"type": "string"},
        "tags": {"type": "array", "items": {"type": "string"}},
        "remoteUrl": {"type": "string"}
      }
    },
    "build": {
      "type": "object",
      "properties": {
        "docker": {
          "type": "object",
          "properties": {
            "present": {"type": "boolean"},
            "baseImage": {"type": "string"},
            "ports": {"type": "array", "items": {"type": "integer"}},
            "commands": {"type": "array", "items": {"type": "string"}}
          }
        },
        "makefile": {
          "type": "object",
          "properties": {
            "present": {"type": "boolean"},
            "targets": {"type": "array", "items": {"type": "string"}}
          }
        }
      }
    }
  }
}
```

## Directory: pkg/gitProcessor

## File: pkg/gitProcessor/gitProcessor.go

Size: 11099 bytes

Last Modified: 2024-12-19T22:01:51Z

```go
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
	// Save the current working directory
	originalDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current working directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			fmt.Printf("Warning: failed to change back to original directory: %v\n", err)
		}
	}()

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

```

## File: pkg/gitProcessor/types.go

Size: 2185 bytes

Last Modified: 2024-12-19T22:04:45Z

```go
package gitProcessor

import "time"

type Options struct {
	CommitHistoryMonths  int
	ReleaseHistoryMonths int
}

type Commit struct {
	Hash    string    `json:"hash"`
	Author  string    `json:"author"`
	Date    time.Time `json:"date"`
	Message string    `json:"message"`
}

type Release struct {
	Tag            string    `json:"tag"`
	Name           string    `json:"name"`
	Date           time.Time `json:"date"`
	Author         string    `json:"author"`
	IsLatestStable bool      `json:"isLatestStable"`
}

type Tag struct {
	Name   string    `json:"name"`
	Date   time.Time `json:"date"`
	Author string    `json:"author"`
}

type Repository struct {
	URL            string    `json:"url"`
	Branch         string    `json:"branch"`
	LastCommit     Commit    `json:"lastCommit"`
	Tags           []Tag     `json:"tags"`
	CommitHistory  []Commit  `json:"commitHistory,omitempty"`
	ReleaseHistory []Release `json:"releaseHistory,omitempty"`
}

type DockerConfig struct {
	Enabled bool     `json:"enabled"`
	Ports   []string `json:"ports,omitempty"`
}

type BuildInfo struct {
	Docker   DockerConfig `json:"docker"`
	Commands []string     `json:"commands,omitempty"`
}

type DependencyInfo struct {
	Language  string            `json:"language"`
	Version   string            `json:"version"`
	Libraries map[string]string `json:"libraries,omitempty"`
}

type DocumentationInfo struct {
	Available bool   `json:"available"`
	API       bool   `json:"api"`
	Summary   string `json:"summary,omitempty"`
}

type Metadata struct {
	AnalyzedAt time.Time `json:"analyzedAt"`
	RepoPath   string    `json:"repoPath"`
	Status     string    `json:"status"`
}

type AnalysisResult struct {
	Metadata      Metadata          `json:"metadata"`
	Repository    Repository        `json:"repository"`
	Build         BuildInfo         `json:"build"`
	Dependencies  DependencyInfo    `json:"dependencies"`
	Documentation DocumentationInfo `json:"documentation"`
	Commit        string            `json:"commit"`
	Branch        string            `json:"branch"`
	Author        string            `json:"author"`
	Timestamp     string            `json:"timestamp"`
	Tag           string            `json:"tag"`
}

```

## Directory: pkg/regions

## Directory: pkg/regions/crm-core-bff

## Directory: pkg/regions/crm-core-bff/scripts

## Directory: pkg/regions/crm-core-bff/scripts/terraform

## Directory: pkg/regions/crm-core-bff/scripts/terraform/aus-prod

## Directory: pkg/regions/crm-core-bff/scripts/terraform/aus-prod/ap-southeast-3

## File: pkg/regions/crm-core-bff/scripts/terraform/aus-prod/ap-southeast-3/_variables.tf

Size: 237 bytes

Last Modified: 2024-12-15T09:35:10Z

```terraform
# region configuration
variable "region" {
  type = string
  default = "ap-southeast-3"
  description = "aws region"
}

# aws account configuration
variable "account" {
  type = string
  default = "aus-prod"
  description = "account"
}


```

## File: pkg/regions/crm-core-bff/scripts/terraform/aus-prod/ap-southeast-3/providers.tf

Size: 168 bytes

Last Modified: 2024-12-15T09:35:10Z

```terraform
provider "aws" {
  region  = var.region
  profile = var.environment
  default_tags {
    tags = {
      Squad   = "psm-crm"
      Service = "crm-core-bff"
    }
  }
}


```

## Directory: pkg/regions/crm-core-bff/scripts/terraform/ind-prod

## Directory: pkg/regions/crm-core-bff/scripts/terraform/ind-prod/ap-south-22

## File: pkg/regions/crm-core-bff/scripts/terraform/ind-prod/ap-south-22/_variables.tf

Size: 240 bytes

Last Modified: 2024-12-15T09:35:10Z

```terraform
# region configuration
variable "region" {
  type    = string
  default = "ap-south-22"
  description = "aws region"
}

# aws account configuration
variable "account" {
  type    = string
  default = "ind-prod"
  description = "account"
}


```

## File: pkg/regions/crm-core-bff/scripts/terraform/ind-prod/ap-south-22/providers.tf

Size: 168 bytes

Last Modified: 2024-12-15T09:35:10Z

```terraform
provider "aws" {
  region  = var.region
  profile = var.environment
  default_tags {
    tags = {
      Squad   = "psm-crm"
      Service = "crm-core-bff"
    }
  }
}


```

## Directory: pkg/regions/crm-core-bff/scripts/terraform/ind-prod/ap-south-3

## File: pkg/regions/crm-core-bff/scripts/terraform/ind-prod/ap-south-3/_variables.tf

Size: 201 bytes

Last Modified: 2024-12-15T09:35:10Z

```terraform
# region configuration
variable "region" {
  default = "ap-south-3"
  description = "aws region"
}

# aws account configuration
variable "account" {
  default = "ind-prod"
  description = "account"
}


```

## File: pkg/regions/crm-core-bff/scripts/terraform/ind-prod/ap-south-3/providers.tf

Size: 168 bytes

Last Modified: 2024-12-15T09:35:10Z

```terraform
provider "aws" {
  region  = var.region
  profile = var.environment
  default_tags {
    tags = {
      Squad   = "psm-crm"
      Service = "crm-core-bff"
    }
  }
}


```

## File: pkg/regions/crm-core-bff/scripts/terraform/locals.tf

Size: 87 bytes

Last Modified: 2024-12-15T09:35:10Z

```terraform
locals {
  project_name         = "crm-core-bff"
  namespace            = "psm-crm"
}


```

## File: pkg/regions/crm-core-bff/scripts/terraform/providers.tf

Size: 117 bytes

Last Modified: 2024-12-15T09:35:10Z

```terraform
provider "aws" {
  default_tags {
    tags = {
      Squad   = "psm-crm"
      Service = "crm-core-bff"
    }
  }
}


```

## File: pkg/regions/go1.mod1

Size: 15 bytes

Last Modified: 2024-12-15T09:35:10Z

```
module regions

```

## File: pkg/regions/main.go

Size: 2489 bytes

Last Modified: 2024-12-15T09:35:10Z

```go
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

```

## Directory: pkg/terraformConfig

## File: pkg/terraformConfig/parser.go

Size: 2532 bytes

Last Modified: 2024-12-15T09:35:10Z

```go
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
```

## Directory: pkg/tf-repo-scanner

## File: pkg/tf-repo-scanner/README.md

Size: 3012 bytes

Last Modified: 2024-12-15T09:35:10Z

```markdown
# Terraform Repository Scanner

This tool scans Terraform files to extract information about repository configurations and outputs the results in either JSON or CSV format. It can process both local directories and GitHub repositories.

## Features

- Scans Terraform files recursively in a directory
- Extracts repository names, team information, and descriptions
- Supports both local directories and GitHub repositories
- Outputs results in JSON or CSV format
- Automatically handles repository cloning and updates
- Uses directory names as team identifiers

## Prerequisites

- Go 1.16 or higher
- Git (for GitHub repository functionality)

## Getting Started Quickly

1. Build the program:
```bash
go build -o terraform-repo-scanner
```

2. Run it on a local directory:
```bash
./terraform-repo-scanner -path ./your-terraform-files
```

Or with a GitHub repository:
```bash
./terraform-repo-scanner -repo https://github.com/org/repo
```

## Installation

```bash
git clone https://github.com/yourusername/terraform-repo-scanner
cd terraform-repo-scanner
go build
```

## Usage

### Basic Command Structure

```bash
./terraform-repo-scanner [flags]
```

### Available Flags

- `-path`: Path to the root directory containing Terraform files (default: ".")
- `-format`: Output format, either "json" or "csv" (default: "json")
- `-repo`: GitHub repository URL (optional)
- `-tmp`: Temporary directory for cloning repositories (default: system temp directory)

### Examples

1. Process local directory:
```bash
./terraform-repo-scanner -path ./terraform-files -format json
```

2. Process GitHub repository:
```bash
./terraform-repo-scanner -repo https://github.com/org/repo -format json
```

3. Process GitHub repository with custom temp directory:
```bash
./terraform-repo-scanner -repo https://github.com/org/repo -tmp ./my-temp -format csv
```

### Example Output

JSON format:
```json
{
  "repositories": [
    {
      "repository_name": "example-repo",
      "team": "team-name",
      "description": "Repository description"
    }
  ],
  "total_count": 1
}
```

CSV format:
```
repository_name,team,description
------------------------------------
example-repo,team-name,Repository description
```

## How It Works

1. The tool walks through the specified directory (or cloned repository) recursively
2. For each `.tf` file found:
    - Extracts module blocks containing repository configurations
    - Uses the parent directory name as the team name
    - Parses repository name and description from the module block
3. Outputs the collected information in the specified format

## Notes

- When using the `-repo` flag, the tool will:
    - Clone the repository if it's not already present in the temp directory
    - Pull latest changes if the repository already exists
    - Process the files as it would for a local directory
- The tool uses the parent directory name of each Terraform file as the team name
- JSON output includes a total count of repositories found

## License

[Your chosen license]
```

## Directory: pkg/tf-repo-scanner/github-repo

## Directory: pkg/tf-repo-scanner/github-repo/psm-accounting

## File: pkg/tf-repo-scanner/github-repo/psm-accounting/repos.tf

Size: 643 bytes

Last Modified: 2024-12-15T09:35:10Z

```terraform
module "module-11" {
  source                 = "github.com/fake-org/tfmod-gh-repo.git?ref=v1.0.4"
  repository_name        = "module-11"
  description            = "Description for module-1"
  required_status_checks = ["module-11/check-1", "module-1/check-2"]
  writers                = ["team-1"]
}

module "module-22" {
  source                 = "github.com/fake-org/tfmod-gh-repo.git?ref=v1.0.4"
  repository_name        = "module-22"
  description            = "Description for module-22"
  required_status_checks = ["module-22/check-1"]
  readers                = ["team-2"]
  writers                = ["team-1", "team-3", "team-4"]
}


```

## Directory: pkg/tf-repo-scanner/github-repo/psm-antifocus

## File: pkg/tf-repo-scanner/github-repo/psm-antifocus/repos.tf

Size: 890 bytes

Last Modified: 2024-12-15T09:35:10Z

```terraform
module "module-1" {
  source                 = "github.com/fake-org/tfmod-gh-repo.git?ref=v1.0.4"
  repository_name        = "module-1"
  description            = "Description for module-1"
  required_status_checks = ["module-1/check-1", "module-1/check-2"]
  writers                = ["team-1"]
}

module "module-2" {
  source                 = "github.com/fake-org/tfmod-gh-repo.git?ref=v1.0.4"
  repository_name        = "module-2"
  description            = "Description for module-2"
  required_status_checks = ["module-2/check-1"]
  readers                = ["team-2"]
  writers                = ["team-1", "team-3", "team-4"]
}

module "module-3" {
  source          = "github.com/fake-org/tfmod-gh-repo.git?ref=v1.0.4"
  repository_name = "module-3"
  description     = "Description for module-3"
  readers         = ["team-2"]
  writers         = ["team-1", "team-3", "team-4"]
}


```

## File: pkg/tf-repo-scanner/go.mod

Size: 23 bytes

Last Modified: 2024-12-15T09:35:10Z

```
module github-projects

```

## File: pkg/tf-repo-scanner/main.go

Size: 4860 bytes

Last Modified: 2024-12-15T09:35:10Z

```go
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

```

## File: server (Skipped - Size exceeds 2MB)

