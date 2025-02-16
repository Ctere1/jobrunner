package jobrunner

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

const DEFAULT_JOB_POOL_SIZE = 10 // Default number of concurrent jobs

var (
	MainCron       *cron.Cron            // Job scheduler singleton instance
	workPermits    chan struct{}         // Limits the number of concurrent jobs
	selfConcurrent bool                  // Whether a job can run concurrently with itself
	HideBanner     bool          = false // Flag to control banner display
)

// ANSI escape codes for colored terminal output
var (
	magenta = "\033[97;45m"
	reset   = "\033[0m"
)

// initWorkPermits initializes the work permits channel
func initWorkPermits(bufferCapacity int) {
	if bufferCapacity <= 0 {
		workPermits = make(chan struct{}, DEFAULT_JOB_POOL_SIZE)
	} else {
		workPermits = make(chan struct{}, bufferCapacity)
	}
}

// setSelfConcurrent sets the self-concurrency flag
func setSelfConcurrent(concurrencyFlag int) {
	selfConcurrent = concurrencyFlag > 0
}

func printBanner() {
	if !HideBanner {
		fmt.Printf("%s[JobRunner] %v Started... %s \n", magenta, time.Now().Format("2006/01/02 - 15:04:05"), reset)
	}
}

// Start initializes and starts the job scheduler
func Start(options ...int) {
	MainCron = cron.New()

	// Apply optional configurations
	if len(options) > 0 {
		initWorkPermits(options[0])
	}
	if len(options) > 1 {
		setSelfConcurrent(options[1])
	}

	MainCron.Start()

	printBanner()
}
