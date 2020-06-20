package jobs

import "time"

type JobState int32

const (
	JobStateRunning  JobState = 1
	JobStateFinished JobState = 2
	JobStateAborted  JobState = 3
)

type Job struct {
	Id                 int64
	State              JobState
	Query              string
	StartTime, EndTime *time.Time
}

type JobStats struct {
	EstimatedProgress    float32
	NumMatchedEvents     int64
	FieldOccurences      map[string]int
	FieldValueOccurences map[string]map[string]int
}
