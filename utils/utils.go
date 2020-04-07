package utils

import (
	"bytes"
	"encoding/json"
	"os/exec"
)

// GetNodesName ,Get Nodes Name in System
func GetNodesName() ([]string, error) {
	// Get Node Information from execution of commandline
	var outBuf bytes.Buffer
	cmd := exec.Command("kubectl", "get", "node", "-o", "json")
	cmd.Stdout = &outBuf
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	// Start To Parse Information
	var nodeInfos struct {
		Items []struct {
			Metadata struct {
				Name string `json:"name"`
			} `json:"metadata"`
		} `json:"items"`
	}

	// Extract to Json
	json.Unmarshal(outBuf.Bytes(), &nodeInfos)

	var result []string
	for _, item := range nodeInfos.Items {
		result = append(result, item.Metadata.Name)
	}

	return result, nil
}
