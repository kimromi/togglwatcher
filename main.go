package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"./notifier"
	"./toggl"
	"github.com/urfave/cli"
)

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
	interval := 5
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

	currentActivities := map[int]toggl.Activity{}

	for {
		select {
		case <-t.C:
			dashboard := toggl.FetchDashboard()

			for _, activity := range dashboard.LatestActivities() {
				currentActivity, exist := currentActivities[activity.UserID]
				if !exist {
					currentActivities[activity.UserID] = activity
					continue
				}

				// Stop
				// current activity is running, and latest activity is stopped
				if currentActivity.Stop == "" && activity.Stop != "" {
					currentActivities[activity.UserID] = activity
					notifier.Notify(notifier.Information{
						Status:      "stopped",
						UserID:      activity.UserID,
						Description: activity.Description,
					})
					continue
				}

				// Start
				// start time is between now and last time check before
				// now <------> start time <------> last time check before
				now := time.Now()
				t, _ := time.Parse("2006-01-02", "1970-01-01")
				start := time.Unix(t.Unix()-activity.Duration, 0)
				before := time.Now().Add(-time.Duration(float64(interval)*1.5) * time.Second)

				if now.After(start) && start.After(before) {
					currentActivities[activity.UserID] = activity
					notifier.Notify(notifier.Information{
						Status:      "started",
						UserID:      activity.UserID,
						Description: activity.Description,
					})
					continue
				}
			}

		case s := <-sig:
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				return
			}
		}
	}
}
