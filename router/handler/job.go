package handler

import (
	"Elven9/Lab-Backend/service/jobs"
	"Elven9/Lab-Backend/utils"
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strconv"

	"github.com/gin-gonic/gin"
)

type jobInfo struct {
	JobName        string `json:"job_name"`
	SubmissionTime string `json:"submission_time"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
	ExeTime        string `json:"exe_time"`
	WaitTime       string `json:"wait_time"`
	State          int    `json:"state"`
}

type getJobsPayload struct {
	jobInfo
	PsCount     int `json:"ps_count"`
	WorkerCount int `json:"worker_count"`
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
	var result []getJobsPayload

	for _, job := range jobs {

		exeTime, _ := job.GetExecutionTime("Minute")
		waitTime, _ := job.GetWaitingTime("Minute")

		job.RecordActivePod()

		var p getJobsPayload
		p.JobName = job.Name
		p.WaitTime = waitTime
		p.State = job.GetState()
		p.ExeTime = exeTime
		p.PsCount = job.GetPsCount()
		p.WorkerCount = job.GetWorkerCount()
		p.SubmissionTime = job.CreateTime
		p.StartTime = job.StartTime
		p.EndTime = job.CompletionTime

		result = append(result, p)
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
		_waitTime, _ := job.GetWaitingTime("Seconds")
		waitTime, _ := strconv.ParseFloat(_waitTime, 64)
		result.JobAverageWaitingTime = (result.JobAverageWaitingTime*float64(result.FinishJobNum+result.FailedJobNum+result.RunningJobNum-1) + waitTime) / float64(result.FinishJobNum+result.FailedJobNum+result.RunningJobNum)
	}

	ctx.JSON(200, result)
}

type workerNodePair struct {
	Node   string `json:"node"`
	Worker string `json:"worker"`
}

type getJobPayload struct {
	Info             jobInfo                `json:"info"`
	DispersionRate   float64                `json:"dispersion_rate"`
	WorkerNodePair   []jobs.WorkerNodePair  `json:"worker_node_pair"`
	CPUUsageOvertime []jobs.PrometheusValue `json:"cpu_usage_overtime"`
}

// GetJob , Get Single Job Information
func GetJob(ctx *gin.Context) {
	// Extract Query
	jobName, isExist := ctx.GetQuery("name")
	if !isExist {
		utils.PushError(400, utils.CustomError{
			Msg: "Query Name Not Found",
		}, nil, ctx)
		return
	}

	// Get Node Information from execution of commandline
	var outBuf bytes.Buffer
	cmd := exec.Command("kubectl", "get", "tfjob", "-o", "json", "--field-selector", fmt.Sprintf("metadata.name=%s", jobName))
	cmd.Stdout = &outBuf
	err := cmd.Run()
	if err != nil {
		utils.PushError(500, utils.CustomError{
			Msg:     "Kubectl Command Execution Failed",
			Command: fmt.Sprintf("kubectl get tfjob -o json --field-selector metadata.name=%s", jobName),
		}, nil, ctx)
		return
	}

	jobs, _ := jobs.CreateJobs(outBuf)

	// Error Handling
	if len(jobs) == 0 {
		utils.PushError(404, utils.CustomError{
			Msg: "Job Specified by UID Does Not Exist",
		}, nil, ctx)
		return
	} else if len(jobs) > 1 {
		utils.PushError(500, utils.CustomError{
			Msg: "Multiple Job Specified by Name Found",
		}, nil, ctx)
		return
	}

	// Compute Time
	exeTime, _ := jobs[0].GetExecutionTime("Minute")
	waitTime, _ := jobs[0].GetWaitingTime("Minute")

	// Compute Dispersion Rate and Get Workers Pod
	err = jobs[0].RecordActivePod()
	if err != nil {
		utils.PushError(500, utils.CustomError{
			Msg: "Server Internal Error",
		}, err, ctx)
		return
	}

	// cpuOvertime, err := jobs[0].GetCPUUsageOvertime()
	// if err != nil {
	// 	utils.PushError(500, utils.CustomError{
	// 		Msg: "Server Internal Error",
	// 	}, err, ctx)
	// 	return
	// }

	// Construct Response Payload
	result := getJobPayload{
		Info: jobInfo{
			JobName:        jobs[0].Name,
			SubmissionTime: jobs[0].CreateTime,
			StartTime:      jobs[0].StartTime,
			EndTime:        jobs[0].CompletionTime,
			ExeTime:        exeTime,
			WaitTime:       waitTime,
			State:          jobs[0].GetState(),
		},
		DispersionRate:   jobs[0].DispersionRate,
		WorkerNodePair:   jobs[0].Workers,
		CPUUsageOvertime: cpuOvertime,
	}

	ctx.JSON(200, result)
}
