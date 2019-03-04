package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"./config"
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
	c := config.LoadConfig()

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

	currentActivities := map[int]toggl.Activity{}
	zone, _ := time.LoadLocation(map[bool]string{true: c.Timezone, false: "UTC"}[c.Timezone != ""])

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

				now := time.Now()

				// Stop
				// current activity is running, and latest activity is stopped
				if currentActivity.Stop == "" && activity.Stop != "" {
					currentActivities[activity.UserID] = activity

					stoppedAt, _ := time.Parse(time.RFC3339, activity.Stop)
					startedAt := stoppedAt.Add(-time.Duration(activity.Duration) * time.Second)

					notifier.Notify(notifier.Information{
						Status:      "stopped",
						UserID:      activity.UserID,
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

				if now.After(started) && started.After(before) {
					currentActivities[activity.UserID] = activity
					notifier.Notify(notifier.Information{
						Status:      "started",
						UserID:      activity.UserID,
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
				return
			}
		}
	}
}
