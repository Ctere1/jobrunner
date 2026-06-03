package jobrunner

import (
	"bytes"
	"log"
	"reflect"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"github.com/robfig/cron/v3"
)

// Job represents a scheduled task within the JobRunner system.
// It wraps an inner cron.Job and maintains execution metadata such as
// name, status, latency, last-run time, and total run count.
//
// lastRun and runCount are placed first in the struct to guarantee 64-bit
// atomic alignment on 32-bit platforms (Go spec requirement).
type Job struct {
	lastRun  int64      `json:"-"` // unix nanoseconds of last execution start; 0 = never run; atomic
	runCount int64      `json:"-"` // total number of executions started; atomic
	Name     string     `json:"name"`
	inner    cron.Job   `json:"-"`
	status   uint32     `json:"-"`
	Status   string     `json:"status"`
	Latency  string     `json:"latency"`
	running  sync.Mutex `json:"-"`
}

const (
	UNNAMED_JOB = "(unnamed)" // Default name for unnamed jobs

	JOB_RUNNING_STATUS = "RUNNING" // Job status constants
	JOB_IDLE_STATUS    = "IDLE"    // Job status constants
)

// New creates a new Job instance from a given cron.Job
func New(job cron.Job) *Job {
	name := reflect.TypeOf(job).Name()
	if name == "Func" {
		name = UNNAMED_JOB
	}
	return &Job{
		Name:  name,
		inner: job,
	}
}

// StatusUpdate updates the job status based on the atomic status flag
func (j *Job) StatusUpdate() string {
	if atomic.LoadUint32(&j.status) > 0 {
		j.Status = JOB_RUNNING_STATUS
	} else {
		j.Status = JOB_IDLE_STATUS
	}
	return j.Status
}

// LastRun returns the time this job last began executing (scheduled or manual).
// Returns the zero time if the job has never run.
func (j *Job) LastRun() time.Time {
	ns := atomic.LoadInt64(&j.lastRun)
	if ns == 0 {
		return time.Time{}
	}
	return time.Unix(0, ns)
}

// RunCount returns the total number of times this job has started executing.
func (j *Job) RunCount() int64 {
	return atomic.LoadInt64(&j.runCount)
}

// Run executes the job and updates the job status, latency, last-run time, and run count.
func (j *Job) Run() {
	start := time.Now()
	atomic.StoreInt64(&j.lastRun, start.UnixNano())
	atomic.AddInt64(&j.runCount, 1)

	// If the job panics, just print a stack trace.
	// Don't let the whole process die.
	defer func() {
		if err := recover(); err != nil {
			var buf bytes.Buffer
			logger := log.New(&buf, "JobRunner Log: ", log.Lshortfile)
			logger.Panic(err, "\n", string(debug.Stack()))
		}
	}()

	if !selfConcurrent {
		j.running.Lock()
		defer j.running.Unlock()
	}

	if workPermits != nil {
		workPermits <- struct{}{}
		defer func() { <-workPermits }()
	}

	atomic.StoreUint32(&j.status, 1)
	j.StatusUpdate()
	defer func() {
		atomic.StoreUint32(&j.status, 0)
		j.StatusUpdate()
	}()

	j.inner.Run()

	j.Latency = time.Since(start).String()
}
