package jobrunner

import (
	"time"

	"github.com/robfig/cron/v3"
)

// StatusData holds the observable state of a single scheduled job entry.
// Prev reflects the last execution start time across all execution paths
// (scheduled and manual), sourced from the Job's own atomic lastRun field.
type StatusData struct {
	Id        cron.EntryID `json:"id"`
	JobRunner *Job         `json:"jobRunner"`
	Next      time.Time    `json:"next"`
	Prev      time.Time    `json:"prev"`
	RunCount  int64        `json:"runCount"`
}

// Entries returns a snapshot of all registered cron entries.
func Entries() []cron.Entry {
	return MainCron.Entries()
}

// StatusPage returns the current state of all registered jobs.
// Prev is taken from the Job's own lastRun field — the single source of truth
// for execution time regardless of whether the run was scheduled or manual.
func StatusPage() []StatusData {
	entries := MainCron.Entries()
	statuses := make([]StatusData, len(entries))

	for i, entry := range entries {
		j := AddJob(entry.Job)
		statuses[i] = StatusData{
			Id:        entry.ID,
			JobRunner: j,
			Next:      entry.Next,
			Prev:      j.LastRun(),
			RunCount:  j.RunCount(),
		}
	}
	return statuses
}

// StatusJson returns the job status as a JSON-serialisable map.
func StatusJson() map[string]interface{} {
	return map[string]interface{}{
		"jobrunner": StatusPage(),
	}
}

// AddJob returns the *Job wrapper from a cron.Job interface value.
func AddJob(job cron.Job) *Job {
	return job.(*Job)
}
