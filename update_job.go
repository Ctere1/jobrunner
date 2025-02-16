package jobrunner

import (
	"errors"
	"time"
)

// RemoveJobByName removes a scheduled job by its name.
func RemoveJobByName(jobName string) error {
	if MainCron == nil {
		return errors.New("scheduler is not initialized")
	}

	// Search for the job by name
	var found bool
	for _, entry := range Entries() {
		job, ok := entry.Job.(*Job)
		if ok && job.Name == jobName {
			MainCron.Remove(entry.ID)
			found = true
		}
	}

	if !found {
		return errors.New("job with given name not found")
	}
	return nil
}

// UpdateJobScheduleByName updates an existing job's schedule by its name.
func UpdateJobScheduleByName(jobName, newSpec string) error {
	if MainCron == nil {
		return errors.New("scheduler is not initialized")
	}

	// Search for the job by name
	var found bool
	var existingJob *Job

	for _, entry := range Entries() {
		job, ok := entry.Job.(*Job)
		if ok && job.Name == jobName {
			existingJob = job
			Remove(entry.ID) // Remove old job
			found = true
		}
	}

	if !found {
		return errors.New("job with given name not found")
	}

	// Add the updated job with the new schedule
	return Schedule(newSpec, existingJob)
}

// UpdateJobIntervalByName updates an existing job's execution interval by its name.
func UpdateJobIntervalByName(jobName string, newInterval time.Duration) error {
	if MainCron == nil {
		return errors.New("scheduler is not initialized")
	}

	// Search for the job by name
	var found bool
	var existingJob *Job

	for _, entry := range Entries() {
		job, ok := entry.Job.(*Job)
		if ok && job.Name == jobName {
			existingJob = job
			Remove(entry.ID) // Remove old job
			found = true
		}
	}

	if !found {
		return errors.New("job with given name not found")
	}

	// Add the updated job with the new interval
	Every(newInterval, existingJob)

	return nil
}
