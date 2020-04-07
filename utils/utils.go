package utils

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"strconv"
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

// ResourceResponse ,
type ResourceResponse struct {
	Name   string `json:"name"`
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

// GetNodesResource ,Get Nodes' Resource in System
func GetNodesResource() ([]ResourceResponse, error) {
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
			Status struct {
				Allocatable struct {
					CPU    string `json:"cpu"`
					Memory string `json:"memory"`
				} `json:"allocatable"`
			} `json:"status"`
		} `json:"items"`
	}

	// Extract to Json
	json.Unmarshal(outBuf.Bytes(), &nodeInfos)

	var result []ResourceResponse
	for _, item := range nodeInfos.Items {
		result = append(result, ResourceResponse{
			Name:   item.Metadata.Name,
			CPU:    item.Status.Allocatable.CPU,
			Memory: item.Status.Allocatable.Memory,
		})
	}

	return result, nil
}

// MemoryConverter ,
func MemoryConverter(str string) int64 {
	runeOfStr := []rune(str)
	unit := string(runeOfStr[len(str)-2 : len(str)])
	result, _ := strconv.ParseInt(string(runeOfStr[0:len(str)-2]), 10, 64)
	if unit == "ki" {
		return result
	} else if unit == "Gi" {
		return result * 1024 * 1024
	} else {
		return result
	}
}
