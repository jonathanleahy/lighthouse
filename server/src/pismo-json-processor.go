package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

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
