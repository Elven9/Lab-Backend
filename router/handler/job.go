package handler

import (
	"bytes"
	"encoding/json"
	"log"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
)

type jobInfo struct {
	JobID          string `json:"job_id"`
	JobType        string `json:"job_type"`
	SubmissionTime string `json:"submission_time"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
	ExeTime        string `json:"exe_time"`
	WaitTime       string `json:"wait_time"`
	State          int    `json:"state"`
}

func isStateSucceed(succeeded int) int {
	if succeeded == 1 {
		return 0
	}
	return 1
}

func GetJobs(ctx *gin.Context) {
	// Get Node Information from execution of commandline
	var outBuf bytes.Buffer
	cmd := exec.Command("kubectl", "get", "job", "-o", "json")
	cmd.Stdout = &outBuf
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Start To Parse Information
	var jobInfos struct {
		Items []struct {
			Metadata struct {
				Name              string `json:"name"`
				UID               string `json:"uid"`
				CreationTimestamp string `json:"creationTimestamp"`
			} `json:"metadata"`
			Status struct {
				CompletionTime string `json:"completionTime"`
				StartTime      string `json:"startTime"`
				Succeeded      int    `json:"succeeded"`
			} `json:"status"`
		} `json:"items"`
	}

	json.Unmarshal(outBuf.Bytes(), &jobInfos)

	// Construct Final Result
	var result []jobInfo

	for _, item := range jobInfos.Items {
		startTimeP, _ := time.Parse(time.RFC3339, item.Status.StartTime)
		endTimeP, _ := time.Parse(time.RFC3339, item.Status.CompletionTime)

		result = append(result, jobInfo{
			JobID:          item.Metadata.UID,
			JobType:        "TODO",
			SubmissionTime: "TODO",
			StartTime:      item.Status.StartTime,
			EndTime:        item.Status.CompletionTime,
			ExeTime:        endTimeP.Sub(startTimeP).String(),
			WaitTime:       "TODO",
			State:          isStateSucceed(item.Status.Succeeded),
		})
	}

	ctx.JSON(200, result)
}
