package util

import (
	"os"
	"strings"

	"io/ioutil"
	"net/http"
	"runtime"
)

type DetectedMetadata struct {
	Provider  string
	Location  string
	NodeID    string
	OSName    string
	OSVersion string
}

// DetectNodeMetadata tries to detect cloud or baremetal metadata
func DetectNodeMetadata() DetectedMetadata {
	meta := DetectedMetadata{
		Provider:  "baremetal",
		Location:  "unknown",
		NodeID:    "",
		OSName:    runtime.GOOS,
		OSVersion: "unknown",
	}

	// Try /etc/os-release for OS info (Linux only)
	if content, err := os.ReadFile("/etc/os-release"); err == nil {
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "NAME=") {
				meta.OSName = strings.Trim(line[5:], `"`)
			}
			if strings.HasPrefix(line, "VERSION_ID=") {
				meta.OSVersion = strings.Trim(line[11:], `"`)
			}
		}
	}

	// AWS EC2 metadata check
	if isReachable("http://169.254.169.254") {
		meta.Provider = "aws"
		meta.Location = tryGet("http://169.254.169.254/latest/meta-data/placement/region")
		meta.NodeID = tryGet("http://169.254.169.254/latest/meta-data/instance-id")
	}

	// GCP metadata
	if isReachable("http://metadata.google.internal") {
		meta.Provider = "gcp"
		meta.Location = tryGet("http://metadata.google.internal/computeMetadata/v1/instance/zone", "Metadata-Flavor: Google")
		meta.NodeID = tryGet("http://metadata.google.internal/computeMetadata/v1/instance/id", "Metadata-Flavor: Google")
	}

	return meta
}

func tryGet(url string, headers ...string) string {
	req, _ := http.NewRequest("GET", url, nil)
	if len(headers) == 2 {
		req.Header.Add(headers[0], headers[1])
	}
	client := http.Client{Timeout: 500 * 1e6} // 500ms
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	return string(data)
}

func isReachable(url string) bool {
	client := http.Client{Timeout: 500 * 1e6}
	_, err := client.Get(url)
	return err == nil
}
