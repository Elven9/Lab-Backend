package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"

	"github.com/gin-gonic/gin"
)

type systemInfo struct {
	NodeName    string `json:"node_name"`
	NodeInfo    string `json:"node_infp"`
	CPUCapacity string `json:"cpu_capacity"`
	CPUCount    string `json:"cpu_count"`
	MemorySize  string `json:"memory_size"`
}

// GetSystemInfo ,Handler Func for System Info.
func GetSystemInfo(ctx *gin.Context) {
	// Get Node Information from execution of commandline
	var outBuf bytes.Buffer
	cmd := exec.Command("kubectl", "get", "node", "-o", "json")
	cmd.Stdout = &outBuf
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Start To Parse Information
	var nodeInfos struct {
		Items []struct {
			Metadata struct {
				Name string `json:"name"`
				UID  string `json:"uid"`
			} `json:"metadata"`
			Status struct {
				Capacity struct {
					CPU    string `json:"cpu"`
					Memory string `json:"memory"`
					Pods   string `json:"pods"`
				} `json:"capacity"`
				NodeInfo struct {
					Architecture string `json:"architecture"`
					OS           string `json:"operatingSystem"`
					OSImage      string `json:"osImage"`
				} `json:"nodeInfo"`
			} `json:"status"`
		} `json:"items"`
	}

	// Extract to Json
	json.Unmarshal(outBuf.Bytes(), &nodeInfos)

	// Construct final payload
	var result []systemInfo
	for _, item := range nodeInfos.Items {
		result = append(result, systemInfo{
			NodeName:    item.Metadata.Name,
			NodeInfo:    fmt.Sprintf("%s %s, %s", item.Status.NodeInfo.OS, item.Status.NodeInfo.Architecture, item.Status.NodeInfo.OSImage),
			CPUCapacity: fmt.Sprintf("%s pod(s)", item.Status.Capacity.Pods),
			CPUCount:    fmt.Sprintf("%s core(s)", item.Status.Capacity.CPU),
			MemorySize:  item.Status.Capacity.Memory,
		})
	}

	ctx.JSON(200, result)
}
