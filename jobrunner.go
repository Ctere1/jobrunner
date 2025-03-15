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
// name, status, and execution latency.
//
// Fields:
// - Name: The name of the job.
// - inner: The actual job implementation that executes.
// - status: Atomic status flag to indicate if the job is running.
// - Status: A human-readable representation of the job status ("RUNNING" or "IDLE").
// - Latency: The time taken for the last execution of the job.
// - running: A mutex to prevent concurrent execution if self-concurrency is disabled.
type Job struct {
	Name    string     `json:"name"`
	inner   cron.Job   `json:"-"`
	status  uint32     `json:"-"`
	Status  string     `json:"status"`
	Latency string     `json:"latency"`
	running sync.Mutex `json:"-"`
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

// Run executes the job and updates the job status and latency
func (j *Job) Run() {
	start := time.Now()
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
