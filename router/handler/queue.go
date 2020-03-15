package handler

import (
	"github.com/gin-gonic/gin"
)

type queueGeneralStat struct {
	QueueType      string  `json:"queue_type"`
	AvgJobWaitTime int     `json:"avg_job_wait_time"`
	AvgQueueLen    float32 `json:"avg_queue_len"`
}

type queueStatPayload struct {
	Interactive []int              `json:"interactive"`
	Train       []int              `json:"train"`
	Service     []int              `json:"service"`
	TimeInfo    []queueGeneralStat `json:"time_info"`
}

type individualQueueInfo struct {
	SubmissionTime string `json:"submission_time"`
	GroupID        int    `json:"group_id"`
	JobID          int    `json:"job_id"`
	UserID         int    `json:"user_id"`
	JobType        string `json:"job_type"`
	YamlLink       string `json:"yaml_link"`
}

// GetQueueStatics ,Get Queue Statics
func GetQueueStatics(ctx *gin.Context) {
	// TODO: Real System Infomation Source
	response := queueStatPayload{
		Interactive: []int{15, 2, 45, 12, 31, 1},
		Train:       []int{1, 34, 2, 21, 45, 6},
		Service:     []int{6, 2, 12, 23, 5, 9},
		TimeInfo: []queueGeneralStat{
			queueGeneralStat{
				QueueType:      "Interactive",
				AvgJobWaitTime: 25,
				AvgQueueLen:    1.2,
			},
			queueGeneralStat{
				QueueType:      "Train",
				AvgJobWaitTime: 190,
				AvgQueueLen:    10.7,
			},
			queueGeneralStat{
				QueueType:      "Service",
				AvgJobWaitTime: 76,
				AvgQueueLen:    5.6,
			},
		},
	}

	ctx.JSON(200, response)
}

// GetInteractiveInfo ,Interactive queue statics
func GetInteractiveInfo(ctx *gin.Context) {
	// TODO: Real System Infomation Source
	response := []individualQueueInfo{
		individualQueueInfo{
			SubmissionTime: "2020/1/1",
			GroupID:        12,
			UserID:         1,
			JobID:          3,
			JobType:        "Interactive",
			YamlLink:       "https://fake-download-link.com.tw",
		},
		individualQueueInfo{
			SubmissionTime: "2020/1/1",
			GroupID:        12,
			UserID:         1,
			JobID:          3,
			JobType:        "Interactive",
			YamlLink:       "https://fake-download-link.com.tw",
		},
		individualQueueInfo{
			SubmissionTime: "2020/1/1",
			GroupID:        12,
			UserID:         1,
			JobID:          3,
			JobType:        "Interactive",
			YamlLink:       "https://fake-download-link.com.tw",
		},
		individualQueueInfo{
			SubmissionTime: "2020/1/1",
			GroupID:        12,
			UserID:         1,
			JobID:          3,
			JobType:        "Interactive",
			YamlLink:       "https://fake-download-link.com.tw",
		},
	}

	ctx.JSON(200, response)
}

// GetTrainInfo ,Train queue statics
func GetTrainInfo(ctx *gin.Context) {
	// TODO: Real System Infomation Source
	response := []individualQueueInfo{
		individualQueueInfo{
			SubmissionTime: "2020/1/1",
			GroupID:        12,
			UserID:         1,
			JobID:          3,
			JobType:        "Train",
			YamlLink:       "https://fake-download-link.com.tw",
		},
		individualQueueInfo{
			SubmissionTime: "2020/1/1",
			GroupID:        12,
			UserID:         1,
			JobID:          3,
			JobType:        "Train",
			YamlLink:       "https://fake-download-link.com.tw",
		},
		individualQueueInfo{
			SubmissionTime: "2020/1/1",
			GroupID:        12,
			UserID:         1,
			JobID:          3,
			JobType:        "Train",
			YamlLink:       "https://fake-download-link.com.tw",
		},
		individualQueueInfo{
			SubmissionTime: "2020/1/1",
			GroupID:        12,
			UserID:         1,
			JobID:          3,
			JobType:        "Train",
			YamlLink:       "https://fake-download-link.com.tw",
		},
	}

	ctx.JSON(200, response)
}

// GetServiceInfo ,Service queue statics
func GetServiceInfo(ctx *gin.Context) {
	// TODO: Real System Infomation Source
	response := []individualQueueInfo{
		individualQueueInfo{
			SubmissionTime: "2020/1/1",
			GroupID:        12,
			UserID:         1,
			JobID:          3,
			JobType:        "Service",
			YamlLink:       "https://fake-download-link.com.tw",
		},
		individualQueueInfo{
			SubmissionTime: "2020/1/1",
			GroupID:        12,
			UserID:         1,
			JobID:          3,
			JobType:        "Service",
			YamlLink:       "https://fake-download-link.com.tw",
		},
		individualQueueInfo{
			SubmissionTime: "2020/1/1",
			GroupID:        12,
			UserID:         1,
			JobID:          3,
			JobType:        "Service",
			YamlLink:       "https://fake-download-link.com.tw",
		},
		individualQueueInfo{
			SubmissionTime: "2020/1/1",
			GroupID:        12,
			UserID:         1,
			JobID:          3,
			JobType:        "Service",
			YamlLink:       "https://fake-download-link.com.tw",
		},
	}

	ctx.JSON(200, response)
}
