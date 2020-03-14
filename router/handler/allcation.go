package handler

import (
	"bytes"
	"encoding/json"
	"log"
	"os/exec"

	"github.com/gin-gonic/gin"
)

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

	ctx.JSON(200, pods)
}
