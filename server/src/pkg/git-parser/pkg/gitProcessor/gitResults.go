// pkg/gitProcessor/result.go
package gitProcessor

import "time"

// AnalysisResult represents the final processed output
type AnalysisResult struct {
	Metadata struct {
		AnalyzedAt time.Time `json:"analyzedAt"`
		RepoPath   string    `json:"repoPath"`
		Status     string    `json:"status"`
	} `json:"metadata"`
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
}
