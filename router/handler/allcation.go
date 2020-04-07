package handler

import (
	"Elven9/Lab-Backend/utils"
	"bytes"
	"encoding/json"
	"log"
	"os/exec"

	"github.com/gin-gonic/gin"
)

type podData struct {
	JobName      string `json:"job_name"`
	ReplicaIndex string `json:"replca_index"`
	ReplicaType  string `json:"replica_type"`
	CPUUsage     string `json:"cpuUsage"`
	MemUsage     string `json:"memUsage"`
}

type nodesData struct {
	NodeName string    `json:"node_name"`
	Pods     []podData `json:"pods"`
}

// GetAllocation ,就是 Allocation 的 Handler
func GetAllocation(ctx *gin.Context) {
	// Get Node Information from execution of commandline
	var outBuf bytes.Buffer
	cmd := exec.Command("kubectl", "get", "pods", "-o", "json")
	cmd.Stdout = &outBuf
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	var pods struct {
		Items []struct {
			Metadata struct {
				Labels struct {
					JobName      string `json:"tf-job-name"`
					ReplicaIndex string `json:"tf-replica-index"`
					ReplicaType  string `json:"tf-replica-type"`
				} `json:"labels"`
				Name string `json:"name"`
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

	// Generate Response Payload
	var res []nodesData
	nodesName, _ := utils.GetNodesName()
	for _, name := range nodesName {
		res = append(res, nodesData{
			NodeName: name,
			Pods:     nil,
		})
	}

	// Populate Data
	for _, item := range pods.Items {
		for idx, node := range res {
			if item.Spec.NodeName == node.NodeName {
				res[idx].Pods = append(res[idx].Pods, podData{
					JobName:      item.Metadata.Labels.JobName,
					ReplicaIndex: item.Metadata.Labels.ReplicaIndex,
					ReplicaType:  item.Metadata.Labels.ReplicaType,
					CPUUsage:     item.Spec.Containers[0].Resources.Requests.CPU,
					MemUsage:     item.Spec.Containers[0].Resources.Requests.Memory,
				})
			}
		}
	}

	ctx.JSON(200, res)
}
