package jobs

// extractInfos , Target Information Need to Extract
type extractInfos struct {
	Items []jobInformation `json:"items"`
}

// JobInformation , Extracted Information Stored in Job Object
type jobInformation struct {
	Metadata struct {
		Name              string `json:"name"`
		UID               string `json:"uid"`
		CreationTimestamp string `json:"creationTimestamp"`
	} `json:"metadata"`
	Status struct {
		CompletionTime string `json:"completionTime"`
		StartTime      string `json:"startTime"`
		Conditions     []SingleJobCondition
	} `json:"status"`
	Spec struct {
		MaxInstances int `json:"max-instances"`
		MinInstances int `json:"min-instances"`
	} `json:"spec"`
}

// SingleJobCondition ,Job Condition Section Wrapper
type SingleJobCondition struct {
	LastUpdateTime string `json:"lastUpdateTime"`
	Reason         string `json:"reason"`
	Type           string `json:"type"`
}
