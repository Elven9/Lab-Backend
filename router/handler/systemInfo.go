package handler

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
)

type systemInfo struct {
	NodeType      string `json:"node_type"`
	NodeCount     int    `json:"node_count"`
	CPUType       string `json:"cpu_type"`
	CPUCapacity   string `json:"cpu_capacity"`
	CPUCount      string `json:"cpu_count"`
	MemorySize    string `json:"memory_size"`
	GPUType       string `json:"gqu_type"`
	GPUCapacity   string `json:"gpu_capacity"`
	GPUCount      string `json:"gpu_count"`
	GPUMemorySize string `json:"gpu_memory_size"`
}

type hardcodedInfo struct {
	NodeName      string `json:"nodeName"`
	CPUType       string `json:"cpuType"`
	GPUType       string `json:"gpuType"`
	GPUCapacity   string `json:"gpuCapacity"`
	GPUCount      string `json:"gpuCount"`
	GPUMemorySize string `json:"gpuMemorySize"`
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
			} `json:"status"`
		} `json:"items"`
	}

	// Extract to Json
	json.Unmarshal(outBuf.Bytes(), &nodeInfos)

	// Construct final payload
	var realHardwareInfo []hardcodedInfo
	hardwareInfoFile, _ := os.Open("./hardwareInfo.json")
	decoder := json.NewDecoder(hardwareInfoFile)
	decoder.Decode(&realHardwareInfo)

	var result []systemInfo
	for _, item := range nodeInfos.Items {
		// real hardware info
		var info hardcodedInfo
		for _, i := range realHardwareInfo {
			if item.Metadata.Name == i.NodeName {
				info = i
			}
		}
		result = append(result, systemInfo{
			NodeType:      item.Metadata.Name,
			NodeCount:     1,
			CPUType:       info.CPUType,
			CPUCapacity:   item.Status.Capacity.Pods,
			CPUCount:      item.Status.Capacity.CPU,
			MemorySize:    item.Status.Capacity.Memory,
			GPUType:       info.GPUType,
			GPUCapacity:   info.GPUCapacity,
			GPUCount:      info.GPUCount,
			GPUMemorySize: info.GPUMemorySize,
		})
	}

	ctx.JSON(200, result)
}
