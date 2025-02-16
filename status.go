package jobrunner

import (
	"time"

	"github.com/robfig/cron/v3"
)

type StatusData struct {
	Id        cron.EntryID
	JobRunner *Job
	Next      time.Time
	Prev      time.Time
}

// Entries returns a list of all cron job entries
func Entries() []cron.Entry {
	return MainCron.Entries()
}

// StatusPage returns a list of all cron job entries in a human-readable format
func StatusPage() []StatusData {
	entries := MainCron.Entries()
	statuses := make([]StatusData, len(entries))

	for i, entry := range entries {
		statuses[i] = StatusData{
			Id:        entry.ID,
			JobRunner: AddJob(entry.Job),
			Next:      entry.Next,
			Prev:      entry.Prev,
		}
	}
	return statuses
}

// StatusJson returns a list of all cron job entries in a JSON format
func StatusJson() map[string]interface{} {
	return map[string]interface{}{
		"jobrunner": StatusPage(),
	}
}

// AddJob returns a Job from a cron.Job
func AddJob(job cron.Job) *Job {
	return job.(*Job)
}
