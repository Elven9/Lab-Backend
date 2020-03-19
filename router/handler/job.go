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
	JobName        string `json:"job_name"`
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

	if len(jc.Conditions) == 0 {
		return -1
	}

	switch jc.Conditions[0].Reason {
	case "TFJobSucceeded":
		return 0
	case "TFJobRunning":
		return 1
	case "TFJobFailed":
		return 2
	}
	return -1
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

		// Prepare Exe Time
		startTimeP, _ := time.Parse(time.RFC3339, item.Status.StartTime)
		endTimeP, _ := time.Parse(time.RFC3339, item.Status.CompletionTime)

		// Create Conditions Object
		var condition jobConditions

		condition.Populate(item.Status.Conditions)

		if len(condition.Conditions) != 1 {
			sort.Sort(&condition)
		}

		exeTime := ""

		if condition.GetState() != -1 && condition.GetState() != 1 {
			// Have Already Failed or Succeeded.
			exeTime = endTimeP.Sub(startTimeP).String()
		}

		result = append(result, jobInfo{
			JobName:        item.Metadata.Name,
			JobID:          item.Metadata.UID,
			JobType:        "TODO",
			SubmissionTime: item.Metadata.CreationTimestamp,
			StartTime:      item.Status.StartTime,
			EndTime:        item.Status.CompletionTime,
			ExeTime:        exeTime,
			WaitTime:       "TODO",
			State:          condition.GetState(),
		})
	}

	ctx.JSON(200, result)
}
