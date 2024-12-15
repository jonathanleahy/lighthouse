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

