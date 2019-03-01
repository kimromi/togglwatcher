package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	yaml "gopkg.in/yaml.v2"

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

	currentActivities := map[int]Activity{}

	for {
		select {
		case <-t.C:
			dashboard := FetchDashboard()

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
					fmt.Printf("%d stop!!!!\n", activity.UserID)
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
					fmt.Printf("%d start!!!!\n", activity.UserID)
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

func (d *Dashboard) LatestActivities() []Activity {
	activities := make([]Activity, 0)

	config := LoadConfig()
	for _, user := range config.Users {
		for _, activity := range d.Activities {
			if user.Id == activity.UserID {
				activities = append(activities, activity)
				break
			}
		}
	}
	return activities
}

type Config struct {
	Api   ApiConfig
	Users []User `yaml:"users"`
}

type ApiConfig struct {
	Token       string `yaml:"token"`
	WorkspaceId int    `yaml:"workspaceid"`
}

type User struct {
	Id   int    `yaml:"id"`
	Name string `yaml:"name"`
}

func LoadConfig() Config {
	config := Config{}

	buf, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		panic(err)
	}

	return config
}

type Dashboard struct {
	MostActiveUsers []MostActiveUser `json:"most_active_user"`
	Activities      []Activity       `json:"activity"`
}

type MostActiveUser struct {
	UserID   int `json:"user_id"`
	Duration int `json:"duration"`
}

type Activity struct {
	UserID      int    `json:"user_id"`
	ProjectID   int    `json:"project_id"`
	Duration    int64  `json:"duration"`
	Description string `json:"description"`
	Stop        string `json:"stop"`
	Tid         int    `json:"tid"`
}

func FetchDashboard() Dashboard {
	config := LoadConfig()

	client := &http.Client{}
	endpoint := fmt.Sprintf("%s%d", "https://www.toggl.com/api/v8/dashboard/", config.Api.WorkspaceId)
	request, err := http.NewRequest("GET", endpoint, nil)
	request.SetBasicAuth(config.Api.Token, "api_token")
	request.Header.Add("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	var dashboard Dashboard
	if err := decoder.Decode(&dashboard); err != nil {
		panic(err)
	}
	return dashboard
}
