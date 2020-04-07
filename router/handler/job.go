package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"sort"
	"strconv"
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
	Conditions     []singleJobCondition
	CompletionTime string
	StartTime      string
	CreateTime     string
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

	if jc.Conditions == nil {
		return 3
	}

	if len(jc.Conditions) == 0 {
		return -1
	}

	if len(jc.Conditions) != 1 {
		sort.Sort(jc)
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

func (jc *jobConditions) GetExecutionTime(format string) (string, error) {

	state := jc.GetState()
	if state != 0 && state != 2 {
		return "", fmt.Errorf("job is still running or waiting")
	}

	// Prepare Exe Time
	startTimeP, err := time.Parse(time.RFC3339, jc.StartTime)
	if err != nil {
		return "", fmt.Errorf("error happend during parsing start time: %v", err)
	}

	endTimeP, err := time.Parse(time.RFC3339, jc.CompletionTime)
	if err != nil {
		return "", fmt.Errorf("error happend during parsing end time: %v", err)
	}

	execTime := endTimeP.Sub(startTimeP).String()

	completeDuration, err := time.ParseDuration(execTime)
	if err != nil {
		return "", fmt.Errorf("error happend during creating complete duration: %v", err)
	}

	if format == "Minute" {
		return execTime, nil
	} else if format == "Second" {
		return fmt.Sprintf("%f", completeDuration.Seconds()), nil
	} else {
		// Unsupported Format
		return "", fmt.Errorf("unsupported format encountered")
	}
}

func (jc *jobConditions) GetWaitingTime() (float64, error) {
	state := jc.GetState()
	if state == 3 {
		return 0, fmt.Errorf("job is still waiting")
	}

	var startTime time.Time
	for _, condition := range jc.Conditions {
		if condition.Reason == "TFJobRunning" {
			startTime, _ = time.Parse(time.RFC3339, condition.LastUpdateTime)
			break
		}
	}

	submitTime, _ := time.Parse(time.RFC3339, jc.CreateTime)
	waitDuration, _ := time.ParseDuration(startTime.Sub(submitTime).String())

	return waitDuration.Seconds(), nil

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

		// Create Conditions Object
		var condition jobConditions

		condition.Populate(item.Status.Conditions)
		condition.StartTime = item.Status.StartTime
		condition.CompletionTime = item.Status.CompletionTime
		condition.CreateTime = item.Metadata.CreationTimestamp

		exeTime, _ := condition.GetExecutionTime("Minute")

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

// SystemwideStatus ,GetSystemwideStatus() response payload format
type SystemwideStatus struct {
	WaitingJobNum            int     `json:"waiting_job_num"`
	RunningJobNum            int     `json:"running_job_num"`
	FinishJobNum             int     `json:"finish_job_num"`
	FailedJobNum             int     `json:"failed_job_num"`
	JobAverageWaitingTime    float64 `json:"job_average_waiting_time"`
	JobAverageCompletionTime float64 `json:"job_average_completion_time"`
}

// GetSystemwideStatus ,Get System-Wide jobs Status
func GetSystemwideStatus(ctx *gin.Context) {
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

	// Create Result Payload
	var result SystemwideStatus

	// Create Conditions Object
	for _, item := range jobInfos.Items {
		var condition jobConditions
		condition.Populate(item.Status.Conditions)
		condition.StartTime = item.Status.StartTime
		condition.CompletionTime = item.Status.CompletionTime
		condition.CreateTime = item.Metadata.CreationTimestamp

		duration, _ := condition.GetExecutionTime("Second")

		execTime, _ := strconv.ParseFloat(duration, 64)

		switch condition.GetState() {
		case 0:
			// Complete Job
			// Calculate New Average Complete Time
			result.JobAverageCompletionTime = (result.JobAverageCompletionTime*float64(result.FinishJobNum) + execTime) / float64(result.FinishJobNum+1)

			result.FinishJobNum++
		case 1:
			// Running
			result.RunningJobNum++
		case 2:
			// Failed
			result.FailedJobNum++
		case 3:
			// Waiting
			result.WaitingJobNum++
			continue
		}

		// Calculate Average Waiting Job Time
		waitTime, _ := condition.GetWaitingTime()
		result.JobAverageWaitingTime = (result.JobAverageWaitingTime*float64(result.FinishJobNum+result.FailedJobNum+result.RunningJobNum-1) + waitTime) / float64(result.FinishJobNum+result.FailedJobNum+result.RunningJobNum)
	}

	ctx.JSON(200, result)
}
