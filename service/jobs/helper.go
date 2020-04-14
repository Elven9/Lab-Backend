package jobs

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// CreateJobs , take a buffer as input an return as a list of job
func CreateJobs(buf bytes.Buffer) ([]Job, error) {
	var jobInfos extractInfos

	// Extract Information
	err := json.Unmarshal(buf.Bytes(), &jobInfos)
	if err != nil {
		return nil, fmt.Errorf("parse info failed at Helper.CreateJobs: %v", err)
	}

	var finalJobs []Job

	for _, item := range jobInfos.Items {
		var job Job

		// Create Job and Push to final list of result
		err := job.Create(item)
		if err != nil {
			return nil, fmt.Errorf("Helper.CreateJobs failed: %v", err)
		}

		finalJobs = append(finalJobs, job)
	}

	return finalJobs, nil
}
