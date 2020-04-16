package handler

import (
	"Elven9/Lab-Backend/service/jobs"
	"Elven9/Lab-Backend/utils"
	"bytes"
	"encoding/json"
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
		waitTime, _ := job.GetWaitingTime("Minute")

		result = append(result, jobInfo{
			JobName:        job.Name,
			JobID:          job.UID,
			SubmissionTime: job.CreateTime,
			StartTime:      job.StartTime,
			EndTime:        job.CompletionTime,
			ExeTime:        exeTime,
			WaitTime:       waitTime,
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
	Info           jobInfo          `json:"info"`
	DispersionRate float64          `json:"dispersion_rate"`
	WorkerNodePair []workerNodePair `json:"worker_node_pair"`
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

	// Compute Dispersion Rate
	outBuf.Reset()
	cmd = exec.Command("kubectl", "get", "pod", "-o", "json", "-l", fmt.Sprintf("tf-job-name=%s", jobName))
	cmd.Stdout = &outBuf
	err = cmd.Run()
	if err != nil {
		utils.PushError(500, utils.CustomError{
			Msg:     "Kubectl Command Execution Failed",
			Command: fmt.Sprintf("kubectl get pod -o json -l tf-job-name=%s", jobName),
		}, err, ctx)
		return
	}

	var jobPods struct {
		Items []struct {
			Metadata struct {
				Labels struct {
					Type  string `json:"tf-replica-type"`
					Index string `json:"tf-replica-index"`
				} `json:"labels"`
			} `json:"metadata"`
			Spec struct {
				NodeName string `json:"nodeName"`
			} `json:"spec"`
		} `json:"items"`
	}

	err = json.Unmarshal(outBuf.Bytes(), &jobPods)
	if err != nil {
		utils.PushError(500, utils.CustomError{
			Msg: "Error Happened During Parse Pod Information",
		}, err, ctx)
		return
	}

	// Calculate How Many Nodes that Have Worker of This TFJOB
	nodeNameMap := make(map[string]int)

	// Construct Node-Pod Pair
	var pairList []workerNodePair

	for _, pod := range jobPods.Items {
		nodeNameMap[pod.Spec.NodeName]++

		pairList = append(pairList, workerNodePair{
			Node:   pod.Spec.NodeName,
			Worker: fmt.Sprintf("%s-%s", pod.Metadata.Labels.Type, pod.Metadata.Labels.Index),
		})
	}

	// Construct Response Payload
	result := getJobPayload{
		Info: jobInfo{
			JobName:        jobs[0].Name,
			JobID:          jobs[0].UID,
			SubmissionTime: jobs[0].CreateTime,
			StartTime:      jobs[0].StartTime,
			EndTime:        jobs[0].CompletionTime,
			ExeTime:        exeTime,
			WaitTime:       waitTime,
			State:          jobs[0].GetState(),
		},
		DispersionRate: float64(len(nodeNameMap)) / float64(jobs[0].MinInstance),
		WorkerNodePair: pairList,
	}

	ctx.JSON(200, result)
}
