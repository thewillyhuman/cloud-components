package types

type ControlPlaneMetadata struct {
	Name   string `json:"name"`
	Region string `json:"region"`
}

type NodeInfo struct {
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`

	Provider string `json:"provider"` // aws, gcp, baremetal, etc.
	Location string `json:"location"` // us-east-1, rack42-dc1, etc.

	NodeID string `json:"node_id"` // cloud instance ID, asset tag, etc.

	OSName    string `json:"os_name"`    // Ubuntu, Debian, etc.
	OSVersion string `json:"os_version"` // 22.04, etc.

	Labels map[string]string `json:"labels"`
}
