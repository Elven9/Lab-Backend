package jobs

// Job , Job Object
type Job struct {
	JobConditions
	Name        string
	UID         string
	MaxInstance int
	MinInstance int
}

// Create , Construct a Job
func (j *Job) Create(info jobInformation) error {
	// Mount Job Information
	j.Name = info.Metadata.Name
	j.UID = info.Metadata.UID
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
