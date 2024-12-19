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
