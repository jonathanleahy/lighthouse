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

