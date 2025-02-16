package jobrunner

import "github.com/robfig/cron/v3"

// Stop ALL active jobs from running at the next scheduled time
func Stop() {
	go MainCron.Stop()
}

// Remove a job from the scheduler by its id
func Remove(id cron.EntryID) {
	MainCron.Remove(id)
}
