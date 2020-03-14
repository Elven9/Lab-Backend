package handler

import (
	"bytes"
	"encoding/json"
	"log"
	"os/exec"

	"github.com/gin-gonic/gin"
)

type target struct {
	Type int    `json:"type"`
	ID   string `json:"id"`
}

type podData struct {
	NodeID string `json:"nodeId"`
	PodID  string `json:"podId"`
}

type userPayload struct {
	Targets []target `json:"identifier"`
}

type responsePayload struct {
	Target target    `json:"identifier"`
	Data   []podData `json:"data"`
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
					JobName string `json:"tf-job-name"`
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

	// Parse Reauest Body
	var requestBody userPayload
	ctx.BindJSON(&requestBody)

	// Generate Response Payload
	var res []responsePayload

	for _, target := range requestBody.Targets {
		// Construct Data Payload
		var data []podData

		for _, pod := range pods.Items {
			if target.Type == 0 {
				// Extract Node Data
				if target.ID == pod.Spec.NodeName {
					data = append(data, podData{
						NodeID: pod.Spec.NodeName,
						PodID:  pod.Metadata.Name,
					})
				}
			} else if target.Type == 1 {
				// Extract Specific Job
				if target.ID == pod.Metadata.Labels.JobName {
					data = append(data, podData{
						NodeID: pod.Spec.NodeName,
						PodID:  pod.Metadata.Name,
					})
				}
			}
		}

		res = append(res, responsePayload{
			Target: target,
			Data:   data,
		})
	}

	ctx.JSON(200, res)
}
