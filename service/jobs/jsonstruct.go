package jobs

// extractInfos , Target Information Need to Extract
type extractInfos struct {
	Items []jobInformation `json:"items"`
}

// Pods Extraction Json Struct
type jobPods struct {
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

// WorkerNodePair ,
type WorkerNodePair struct {
	Node   string `json:"node"`
	Worker string `json:"worker"`
}

// JobInformation , Extracted Information Stored in Job Object
type jobInformation struct {
	Metadata struct {
		Name              string `json:"name"`
		CreationTimestamp string `json:"creationTimestamp"`
		Namespace         string `json:"namespace"`
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
