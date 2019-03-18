package main

import (
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	log = logrus.New()
)

func init() {
	log.SetFormatter(&logrus.JSONFormatter{})

	c, err := LoadConfig()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	logFile := map[bool]string{true: c.LogFile, false: "./togglwatcher.log"}[c.LogFile != ""]
	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Error("logFile opening failed")
		os.Exit(1)
	}
	mw := io.MultiWriter(os.Stdout, file)
	log.SetOutput(mw)
}

func main() {
	app := cli.NewApp()

	app.Name = "togglwatcher"
	app.Version = "0.1.0"

	app.Action = func(c *cli.Context) {
		Watch()
	}

	app.Run(os.Args)
}

func Watch() {
	c, _ := LoadConfig()

	interval := map[bool]int{true: c.Interval, false: 10}[c.Interval > 0]
	t := time.NewTicker(time.Duration(interval) * time.Second)
	defer t.Stop()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer signal.Stop(sig)

	users, err := FetchUsers()
	if err != nil {
		log.Error("users fetching failed")
		os.Exit(1)
	}

	currentActivities := map[int]Activity{}
	zone, _ := time.LoadLocation(map[bool]string{true: c.Timezone, false: "UTC"}[c.Timezone != ""])

	for {
		select {
		case <-t.C:
			dashboard, err := FetchDashboard()
			if err != nil {
				log.Error("dashboard fetching failed")
				continue
			}

			for _, activity := range dashboard.LatestActivities(users) {
				currentActivity, exist := currentActivities[activity.UserID]
				if !exist {
					currentActivities[activity.UserID] = activity
					continue
				}

				now := time.Now()

				// Stop
				// current activity is running, and latest activity is stopped
				if currentActivity.Stop == "" && activity.Stop != "" {
					currentActivities[activity.UserID] = activity

					stoppedAt, _ := time.Parse(time.RFC3339, activity.Stop)
					startedAt := stoppedAt.Add(-time.Duration(activity.Duration) * time.Second)

					Notify(NotifyInformation{
						Status:      "stopped",
						User:        users[activity.UserID],
						Description: activity.Description,
						StartedAt:   startedAt.In(zone).Format("1/2 15:04"),
						StoppedAt:   stoppedAt.In(zone).Format("1/2 15:04"),
					})
					continue
				}

				// Start
				// start time is between now and last time check before
				// now <------> start time <------> last time check before
				t, _ := time.Parse("2006-01-02", "1970-01-01")
				started := time.Unix(t.Unix()-activity.Duration, 0)
				before := time.Now().Add(-time.Duration(float64(interval)*1.5) * time.Second)

				if now.After(started) && started.After(before) && currentActivity.Description != activity.Description {
					currentActivities[activity.UserID] = activity
					Notify(NotifyInformation{
						Status:      "started",
						User:        users[activity.UserID],
						Description: activity.Description,
						StartedAt:   started.In(zone).Format("1/2 15:04"),
						StoppedAt:   "-",
					})
					continue
				}
			}

		case s := <-sig:
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				log.Info("togglwatcher terminated")
				return
			}
		}
	}
}
