package jobrunner

import (
	"reflect"
	"sync/atomic"
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

// findRegistered returns the *Job registered under the cron scheduler whose
// name matches the given job type, or nil if not found.
// Reusing the registered instance ensures that lastRun, RunCount, Status,
// and Latency written during execution are visible via StatusPage.
func findRegistered(name string) *Job {
	for _, entry := range MainCron.Entries() {
		if j, ok := entry.Job.(*Job); ok && j.Name == name {
			return j
		}
	}
	return nil
}

// Now runs the given job once, immediately.
//
// It pre-stamps lastRun synchronously before launching the goroutine so that
// any StatusPage call immediately after Now() returns sees the updated Prev —
// important for nudge-based WebSocket pushes. The goroutine's Run() will
// overwrite lastRun ~1ms later with the actual start time; both writes are
// atomic and the delta is negligible.
//
// It reuses the registered *Job instance when possible so that Status and
// Latency updates during the run are visible in StatusPage.
func Now(job cron.Job) {
	name := reflect.TypeOf(job).Name()
	j := findRegistered(name)
	if j == nil {
		j = New(job)
	}
	atomic.StoreInt64(&j.lastRun, time.Now().UnixNano())
	go j.Run()
}

// In runs the given job once, after the given duration.
func In(duration time.Duration, job cron.Job) {
	go func() {
		time.Sleep(duration)
		New(job).Run()
	}()
}
