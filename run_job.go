package jobrunner

import (
	"time"

	"github.com/robfig/cron/v3"
)

// Func is a wrapper for a function to be used as a job
type Func func()

func (r Func) Run() { r() }

// Schedule calls the given job on the given schedule.
func Schedule(spec string, job cron.Job) error {
	sched, err := cron.ParseStandard(spec)
	if err != nil {
		return err
	}
	MainCron.Schedule(sched, New(job))
	return nil
}

// Every runs the given job at the given interval.
func Every(duration time.Duration, job cron.Job) {
	MainCron.Schedule(cron.Every(duration), New(job))
}

// Now runs the given job once, immediately.
func Now(job cron.Job) {
	go New(job).Run()
}

// In runs the given job once, after the given duration.
func In(duration time.Duration, job cron.Job) {
	go func() {
		time.Sleep(duration)
		New(job).Run()
	}()
}
