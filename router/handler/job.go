package handler

import (
	"bytes"
	"encoding/json"
	"log"
	"os/exec"
	"sort"
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

type singleJobCondition struct {
	LastUpdateTime string `json:"lastUpdateTime"`
	Reason         string `json:"reason"`
	Type           string `json:"type"`
}

// Sorting Functionality
type jobConditions struct {
	Conditions []singleJobCondition
}

// Interface Implementation
func (jc *jobConditions) Populate(list []singleJobCondition) {
	for _, item := range list {
		jc.Conditions = append(jc.Conditions, item)
	}
}

func (jc *jobConditions) Len() int { return len(jc.Conditions) }
func (jc *jobConditions) Swap(i, j int) {
	jc.Conditions[i], jc.Conditions[j] = jc.Conditions[j], jc.Conditions[i]
}
func (jc *jobConditions) Less(i, j int) bool {
	iTime, _ := time.Parse(time.RFC3339, jc.Conditions[i].LastUpdateTime)
	jTime, _ := time.Parse(time.RFC3339, jc.Conditions[j].LastUpdateTime)

	return !iTime.Before(jTime)
}

func (jc *jobConditions) GetState() int {
	// State 0: Success
	// State 1: Running
	// State 2: Failed
	// State 3: Waiting

	switch jc.Conditions[0].Reason {
	case "TFJobSucceeded":
		return 0
	case "TFJobRunning":
		return 1
	case "TFJobFailed":
		return 2
	}
	return 2
}

// GetJobs ,取得 TF JOBS 的資料
func GetJobs(ctx *gin.Context) {
	// Get Node Information from execution of commandline
	var outBuf bytes.Buffer
	cmd := exec.Command("kubectl", "get", "tfjob", "-o", "json")
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
				Conditions     []singleJobCondition
			} `json:"status"`
		} `json:"items"`
	}

	json.Unmarshal(outBuf.Bytes(), &jobInfos)

	// Construct Final Result
	var result []jobInfo

	for _, item := range jobInfos.Items {
		startTimeP, _ := time.Parse(time.RFC3339, item.Status.StartTime)
		endTimeP, _ := time.Parse(time.RFC3339, item.Status.CompletionTime)

		// Create Conditions Object
		var condition jobConditions

		condition.Populate(item.Status.Conditions)
		sort.Sort(&condition)

		result = append(result, jobInfo{
			JobID:          item.Metadata.UID,
			JobType:        "TODO",
			SubmissionTime: item.Metadata.CreationTimestamp,
			StartTime:      item.Status.StartTime,
			EndTime:        item.Status.CompletionTime,
			ExeTime:        endTimeP.Sub(startTimeP).String(),
			WaitTime:       "TODO",
			State:          condition.GetState(),
		})
	}

	ctx.JSON(200, result)
}
