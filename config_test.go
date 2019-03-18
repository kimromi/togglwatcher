package main

import "testing"

func TestLoadConfig(t *testing.T) {
	c, _ := LoadConfig("config.yaml.sample")

	if c.Timezone != "Asia/Tokyo" {
		t.Errorf("add failed. expect:%s, actual:%s", "Asia/Tokyo", c.Timezone)
	}
	if c.Interval != 10 {
		t.Errorf("add failed. expect:%d, actual:%d", 10, c.Interval)
	}
	if c.LogFile != "./togglwatcher.log" {
		t.Errorf("add failed. expect:%s, actual:%s", "./togglwatcher.log", c.LogFile)
	}

	if c.Api.Token != "dummy" {
		t.Errorf("add failed. expect:%s, actual:%s", "dummy", c.Api.Token)
	}
	if c.Api.DashboardId != 1234567 {
		t.Errorf("add failed. expect:%d, actual:%d", 1234567, c.Api.DashboardId)
	}

	for _, notification := range c.Notifications {
		if notification.Service == "slack" {
			if notification.Url != "https://hooks.slack.com/services/xxx" {
				t.Errorf("add failed. expect:%s, actual:%s", "https://hooks.slack.com/services/xxx", notification.Url)
			}
			if notification.Channel != "general" {
				t.Errorf("add failed. expect:%s, actual:%s", "general", notification.Channel)
			}
		}
	}
}
