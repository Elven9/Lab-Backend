package jobs

// Job , Job Object
type Job struct {
	JobConditions
	Name string
	UID  string
}

// Create , Construct a Job
func (j *Job) Create(info jobInformation) error {
	// Mount Job Information
	j.Name = info.Metadata.Name
	j.UID = info.Metadata.UID

	// Create Condition
	j.Populate(info.Status.Conditions)
	j.StartTime = info.Status.StartTime
	j.CompletionTime = info.Status.CompletionTime
	j.CreateTime = info.Metadata.CreationTimestamp

	// Currently Won't have any Error
	return nil
}
