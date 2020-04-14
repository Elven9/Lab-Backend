package handler

import (
	"Elven9/Lab-Backend/service/jobs"
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strconv"

	"github.com/gin-gonic/gin"
)

type jobInfo struct {
	JobName        string `json:"job_name"`
	JobID          string `json:"job_id"`
	SubmissionTime string `json:"submission_time"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
	ExeTime        string `json:"exe_time"`
	WaitTime       string `json:"wait_time"`
	State          int    `json:"state"`
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

	// Get Jobs
	jobs, _ := jobs.CreateJobs(outBuf)

	// Construct Final Result
	var result []jobInfo

	for _, job := range jobs {

		exeTime, _ := job.GetExecutionTime("Minute")
		waitTime, _ := job.GetWaitingTime()

		result = append(result, jobInfo{
			JobName:        job.Name,
			JobID:          job.UID,
			SubmissionTime: job.CreateTime,
			StartTime:      job.StartTime,
			EndTime:        job.CompletionTime,
			ExeTime:        exeTime,
			WaitTime:       fmt.Sprintf("%f Sec.", waitTime),
			State:          job.GetState(),
		})
	}

	ctx.JSON(200, result)
}

// SystemwideStatus ,GetSystemwideStatus() response payload format
type systemwideStatus struct {
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

	// Get Jobs
	jobs, _ := jobs.CreateJobs(outBuf)

	// Create Result Payload
	var result systemwideStatus

	// Create Conditions Object
	for _, job := range jobs {

		duration, _ := job.GetExecutionTime("Second")

		execTime, _ := strconv.ParseFloat(duration, 64)

		switch job.GetState() {
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
		waitTime, _ := job.GetWaitingTime()
		result.JobAverageWaitingTime = (result.JobAverageWaitingTime*float64(result.FinishJobNum+result.FailedJobNum+result.RunningJobNum-1) + waitTime) / float64(result.FinishJobNum+result.FailedJobNum+result.RunningJobNum)
	}

	ctx.JSON(200, result)
}
