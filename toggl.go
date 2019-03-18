package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	Id   int    `json:"id"`
	Uid  int    `json:"uid"`
	Name string `json:"name"`
}

func FetchUsers() (map[int]User, error) {
	c, _ := LoadConfig()

	client := &http.Client{}
	endpoint := fmt.Sprintf("%s%d%s", "https://www.toggl.com/api/v8/workspaces/", c.Api.DashboardId, "/workspace_users")
	request, err := http.NewRequest("GET", endpoint, nil)
	request.SetBasicAuth(c.Api.Token, "api_token")
	request.Header.Add("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	var users []User
	if err := decoder.Decode(&users); err != nil {
		return nil, err
	}

	teamUsers := map[int]User{}
	for _, user := range users {
		teamUsers[user.Uid] = user
	}
	return teamUsers, nil
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

func FetchDashboard() (*Dashboard, error) {
	c, _ := LoadConfig()

	client := &http.Client{}
	endpoint := fmt.Sprintf("%s%d", "https://www.toggl.com/api/v8/dashboard/", c.Api.DashboardId)
	request, err := http.NewRequest("GET", endpoint, nil)
	request.SetBasicAuth(c.Api.Token, "api_token")
	request.Header.Add("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	var dashboard Dashboard
	if err := decoder.Decode(&dashboard); err != nil {
		return nil, err
	}
	return &dashboard, nil
}

func (d *Dashboard) LatestActivities(users map[int]User) []Activity {
	activities := make([]Activity, 0)

	for _, activity := range d.Activities {
		if _, exists := users[activity.UserID]; exists {
			activities = append(activities, activity)
			break
		}
	}
	return activities
}
