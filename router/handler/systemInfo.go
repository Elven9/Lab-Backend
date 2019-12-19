package handler

import (
	"github.com/gin-gonic/gin"
)

type systemInfo struct {
	NodeType      string `json:"node_type"`
	NodeCount     int    `json:"node_count"`
	CPUType       string `json:"cpu_type"`
	CPUCapacity   int    `json:"cpu_capacity"`
	CPUCount      int    `json:"cpu_count"`
	MemorySize    string `json:"memory_size"`
	GPUType       string `json:"gqu_type"`
	GPUCapacity   int    `json:"gpu_capacity"`
	GPUCount      int    `json:"gpu_count"`
	GPUMemorySize string `json:"gpu_memory_size"`
}

// GetSystemInfo ,Handler Func for System Info.
func GetSystemInfo(ctx *gin.Context) {
	// TODO: Real System Infomation Source
	// Currently implemented with static info, hardcoded.
	response := [4]systemInfo{
		systemInfo{
			NodeType:      "Test Node Type",
			NodeCount:     5,
			CPUType:       "i7-9700【8核/8緒】3.0GHz(↑4.7GHz)/12M/UHD630/65W",
			CPUCapacity:   12,
			CPUCount:      4,
			MemorySize:    "32 GB",
			GPUType:       "NVIDIA® Tesla® P100",
			GPUCapacity:   12,
			GPUCount:      2,
			GPUMemorySize: "64 GB",
		},
		systemInfo{
			NodeType:      "Test Node Type",
			NodeCount:     5,
			CPUType:       "i7-9700【8核/8緒】3.0GHz(↑4.7GHz)/12M/UHD630/65W",
			CPUCapacity:   12,
			CPUCount:      4,
			MemorySize:    "32 GB",
			GPUType:       "NVIDIA® Tesla® P100",
			GPUCapacity:   12,
			GPUCount:      2,
			GPUMemorySize: "64 GB",
		},
		systemInfo{
			NodeType:      "Test Node Type",
			NodeCount:     5,
			CPUType:       "i7-9700【8核/8緒】3.0GHz(↑4.7GHz)/12M/UHD630/65W",
			CPUCapacity:   12,
			CPUCount:      4,
			MemorySize:    "32 GB",
			GPUType:       "NVIDIA® Tesla® P100",
			GPUCapacity:   12,
			GPUCount:      2,
			GPUMemorySize: "64 GB",
		},
		systemInfo{
			NodeType:      "Test Node Type",
			NodeCount:     5,
			CPUType:       "i7-9700【8核/8緒】3.0GHz(↑4.7GHz)/12M/UHD630/65W",
			CPUCapacity:   12,
			CPUCount:      4,
			MemorySize:    "32 GB",
			GPUType:       "NVIDIA® Tesla® P100",
			GPUCapacity:   12,
			GPUCount:      2,
			GPUMemorySize: "64 GB",
		},
	}

	ctx.JSON(200, response)
}
