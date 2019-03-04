package main

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Timezone      string `yaml:"timezone"`
	Interval      int    `yaml:"interval"`
	LogFile       string `yaml:"logfile"`
	Api           ApiConfig
	Users         []User         `yaml:"users"`
	Notifications []Notification `yaml:"notifications"`
}

type ApiConfig struct {
	Token       string `yaml:"token"`
	DashboardId int    `yaml:"dashboardid"`
}

type User struct {
	Id   int    `yaml:"id"`
	Name string `yaml:"name"`
}

type Notification struct {
	Service string `yaml:"service"`
	Url     string `yaml:"webhook_url"`
	Channel string `yaml:"channel"`
	Name    string `yaml:"name"`
}

func LoadConfig() (*Config, error) {
	config := Config{}

	buf, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
