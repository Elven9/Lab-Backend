package jobs

import (
	"fmt"
	"sort"
	"time"
)

// JobConditions , Main Object Implement Sorting and Extraction of Job Condition Information
type JobConditions struct {
	Conditions     []SingleJobCondition
	CompletionTime string
	StartTime      string
	CreateTime     string
}

// Interface Implementation
func (jc *JobConditions) Populate(list []SingleJobCondition) {
	for _, item := range list {
		jc.Conditions = append(jc.Conditions, item)
	}
}

func (jc *JobConditions) Len() int { return len(jc.Conditions) }
func (jc *JobConditions) Swap(i, j int) {
	jc.Conditions[i], jc.Conditions[j] = jc.Conditions[j], jc.Conditions[i]
}
func (jc *JobConditions) Less(i, j int) bool {
	iTime, _ := time.Parse(time.RFC3339, jc.Conditions[i].LastUpdateTime)
	jTime, _ := time.Parse(time.RFC3339, jc.Conditions[j].LastUpdateTime)

	return !iTime.Before(jTime)
}

func (jc *JobConditions) GetState() int {
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

func (jc *JobConditions) GetExecutionTime(format string) (string, error) {

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

func (jc *JobConditions) GetWaitingTime(format string) (string, error) {
	state := jc.GetState()
	if state == 3 {
		return "", fmt.Errorf("job is still waiting")
	}

	var startTime time.Time
	for _, condition := range jc.Conditions {
		if condition.Reason == "TFJobRunning" {
			startTime, _ = time.Parse(time.RFC3339, condition.LastUpdateTime)
			break
		}
	}

	submitTime, _ := time.Parse(time.RFC3339, jc.CreateTime)
	waitTime := startTime.Sub(submitTime).String()
	waitDuration, _ := time.ParseDuration(startTime.Sub(submitTime).String())

	if format == "Minute" {
		return waitTime, nil
	} else if format == "Second" {
		return fmt.Sprintf("%f", waitDuration.Seconds()), nil
	} else {
		// Unsupported Format
		return "", fmt.Errorf("unsupported format encountered")
	}
}
