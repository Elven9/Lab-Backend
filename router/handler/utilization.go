package handler

import (
	"Elven9/Lab-Backend/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strconv"

	"github.com/gin-gonic/gin"
)

type utilizationData struct {
	NodeName      string `json:"node_name"`
	MaximaCPU     string `json:"maxima_cpu"`
	AggregateCPU  string `json:"aggregate_cpu"`
	MaximaMem     string `json:"maxima_mem"`
	AggregateMem  string `json:"aggregate_mem"`
	CPUPercentage int    `json:"cpu_percentage"`
	MemPercentage int    `json:"mem_percentage"`
}

// GetUtilization ,
func GetUtilization(ctx *gin.Context) {
	// Generate Response Payload
	var res []utilizationData
	nodeResources, _ := utils.GetNodesResource()
	for _, node := range nodeResources {
		var outBuf bytes.Buffer
		cmd := exec.Command("kubectl", "get", "pods", "-o", "json", "--field-selector", fmt.Sprintf("spec.nodeName=%s", node.Name))
		cmd.Stdout = &outBuf
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}

		var pods struct {
			Items []struct {
				Metadata struct {
					Labels struct {
						ReplicaType string `json:"tf-replica-type"`
					} `json:"labels"`
				} `json:"metadata"`
				Spec struct {
					Containers []struct {
						Resources struct {
							Requests struct {
								CPU    string `json:"cpu"`
								Memory string `json:"memory"`
							} `json:"requests"`
						} `json:"resources"`
					} `json:"containers"`
					NodeName string `json:"nodeName"`
				} `json:"spec"`
			} `json:"items"`
		}

		json.Unmarshal(outBuf.Bytes(), &pods)

		// Prepare Payload
		payload := utilizationData{
			NodeName:  node.Name,
			MaximaCPU: node.CPU,
			MaximaMem: node.Memory,
		}

		var curCPU int8 = 0
		var curMem int64 = 0
		for _, pod := range pods.Items {
			if pod.Metadata.Labels.ReplicaType != "worker" {
				continue
			}
			podCPU, _ := strconv.ParseInt(pod.Spec.Containers[0].Resources.Requests.CPU, 10, 8)
			curCPU += int8(podCPU)

			curMem += utils.MemoryConverter(pod.Spec.Containers[0].Resources.Requests.Memory)
		}

		payload.AggregateCPU = strconv.Itoa(int(curCPU))
		payload.AggregateMem = strconv.Itoa(int(curMem))
		maximaCPU, _ := strconv.ParseInt(payload.MaximaCPU, 10, 64)
		payload.CPUPercentage = int(float64(curCPU) / float64(maximaCPU) * 100)
		payload.MemPercentage = int(float64(curMem) / float64(utils.MemoryConverter(payload.MaximaMem)) * 100)

		res = append(res, payload)
	}

	ctx.JSON(200, res)
}
