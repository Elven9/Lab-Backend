package jobs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// Job , Job Object
type Job struct {
	JobConditions
	Name           string
	MaxInstance    int
	MinInstance    int
	Workers        []WorkerNodePair
	DispersionRate float64
}

// Create , Construct a Job
func (j *Job) Create(info jobInformation) error {
	// Mount Job Information
	j.Name = fmt.Sprintf("%s/%s", info.Metadata.Namespace, info.Metadata.Name)
	j.MaxInstance = info.Spec.MaxInstances
	j.MinInstance = info.Spec.MinInstances

	// Create Condition
	j.Populate(info.Status.Conditions)
	j.StartTime = info.Status.StartTime
	j.CompletionTime = info.Status.CompletionTime
	j.CreateTime = info.Metadata.CreationTimestamp

	// Currently Won't have any Error
	return nil
}

// GetPsCount , Get The Amount of Parameter server this job currently have.
func (j *Job) GetPsCount() int {
	result := 0

	for _, pod := range j.Workers {
		if strings.Split(pod.Worker, "-")[0] == "ps" {
			result++
		}
	}

	return result
}

// GetWorkerCount ,Get the amount of worker this job currently have.
func (j *Job) GetWorkerCount() int {
	result := 0

	for _, pod := range j.Workers {
		if strings.Split(pod.Worker, "-")[0] == "worker" {
			result++
		}
	}

	return result
}

// RecordActivePod , Record Pod Currently Run in System
func (j *Job) RecordActivePod() error {
	var outBuf bytes.Buffer
	cmd := exec.Command("kubectl", "get", "pod", "-o", "json", "-l", fmt.Sprintf("tf-job-name=%s", strings.Split(j.Name, "/")[1]))
	cmd.Stdout = &outBuf
	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("RecordActivePod Failed: Kubectl Command Execution Failed: %v", err)
	}

	// Unmartial Data
	var pods jobPods
	err = json.Unmarshal(outBuf.Bytes(), &pods)
	if err != nil {
		return fmt.Errorf("RecordActivePod Failed: Parse Pod Data Failed: %v", err)
	}

	// Calculate How Many Nodes that Have Worker of This TFJOB
	nodeNameMap := make(map[string]int)

	// Construct Node-Pod Pair
	for _, pod := range pods.Items {
		nodeNameMap[pod.Spec.NodeName]++

		j.Workers = append(j.Workers, WorkerNodePair{
			Node:   pod.Spec.NodeName,
			Worker: fmt.Sprintf("%s-%s", pod.Metadata.Labels.Type, pod.Metadata.Labels.Index),
		})
	}

	j.DispersionRate = float64(len(nodeNameMap)) / float64(j.MinInstance)

	return nil
}
