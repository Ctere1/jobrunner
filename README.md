[![Go Report Card](https://goreportcard.com/badge/github.com/Ctere1/jobrunner)](https://goreportcard.com/report/github.com/Ctere1/jobrunner) [![pkg.go.dev](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat)](https://pkg.go.dev/github.com/Ctere1/jobrunner)

# ![](./views/runclock.jpg) JobRunner

This is a fork of the [bamzi/jobrunner](https://github.com/bamzi/jobrunner) package.

JobRunner is framework for performing work asynchronously, outside of the request flow. It comes with cron to schedule and queue job functions for processing at specified time. 

It includes a live monitoring of current schedule and state of active jobs that can be outputed as JSON or Html template. 

## Install

`go get github.com/Ctere1/jobrunner`

### Setup

```go
package main

import "github.com/Ctere1/jobrunner"

func main() {
    jobrunner.Start() // optional: jobrunner.Start(pool int, concurrent int) (10, 1)
    jobrunner.Schedule("@every 5s", ReminderEmails{})
}

// Job Specific Functions
type ReminderEmails struct {
    // filtered
}

// ReminderEmails.Run() will get triggered automatically.
func (e ReminderEmails) Run() {
    // Queries the DB
    // Sends some email
    fmt.Printf("Every 5 sec send reminder emails \n")
}
```

### Live Monitoring
![](https://raw.githubusercontent.com/Ctere1/jobrunner/master/views/jobrunner-html.png)
```go

// Example of GIN micro framework
func main() {
    routes := gin.Default()

    // Resource to return the JSON data
    routes.GET("/jobrunner/json", JobJson)

    // Load template file location relative to the current working directory
    routes.LoadHTMLGlob("../github.com/Ctere1/jobrunner/views/Status.html")

    // Returns html page at given endpoint based on the loaded
    // template from above
    routes.GET("/jobrunner/html", JobHtml)

    routes.Run(":8080")
}

func JobJson(c *gin.Context) {
    // returns a map[string]interface{} that can be marshalled as JSON
    c.JSON(200, jobrunner.StatusJson())
}

func JobHtml(c *gin.Context) {
    // Returns the template data pre-parsed
    c.HTML(200, "", jobrunner.StatusPage())

}

```
## WHY?
To reduce our http response latency by 200+%

JobRunner was created to help us process functions unrelated to response without any delays to the http response. GoRoutines would timeout because response has already been processed, closing the instance window all together. 

Instead of creating a separate independent app, we wanted to save time and manageability of our current app by coupling-in the job processing. We did not want to have micro services. It's premature optimization.

If you have a web app or api service backend and want a job processing framework built into your app then JobRunner is for you. Once you hit mass growth and need to scale, you can simply decouple you JobRunners into a dedicated app.

## Use cases
Here are some examples of what we use JobRunner for:

* Send emails to new users after signup
* Sending push notification or emails based on specifics
* ReMarketing Engine - send invites, reminder emails, etc ...
* Clean DB, data or AMZ S3
* Sending Server stats to monitoring apps
* Send data stats at daily or weekly intervals

### Supported Featured
*All jobs are processed outside of the request flow*

* Now: process a job immediately
* In: processing a job one time, after a given time
* Every: process a recurring job after every given time gap
* Schedule: process a job (recurring or otherwise) at a given time


## Compatibility

JobRunner is designed to be framework agnostic. So it will work with pure Go apps as well as any framework written in Go Language. 

*Verified Supported Frameworks*

* Gin
* Echo
* Martini
* Beego
* Revel (JobRunner is modified version of revel's Jobs module)
* ...


## Basics

```go
    jobrunner.Schedule("* */5 * * * *", DoSomething{}) // every 5min do something
    jobrunner.Schedule("@every 1h30m10s", ReminderEmails{})
    jobrunner.Schedule("@midnight", DataStats{}) // every midnight do this..
    jobrunner.Every(16*time.Minute, CleanS3{}) // evey 16 min clean...
    jobrunner.In(10*time.Second, WelcomeEmail{}) // one time job. starts after 10sec
    jobrunner.Now(NowDo{}) // do the job as soon as it's triggered
```
[**More Detailed CRON Specs**](https://github.com/robfig/cron/blob/v2/doc.go)
