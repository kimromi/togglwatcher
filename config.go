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
	Notifications []Notification `yaml:"notifications"`
}

type ApiConfig struct {
	Token       string `yaml:"token"`
	DashboardId int    `yaml:"dashboardid"`
}

type Notification struct {
	Service string `yaml:"service"`
	Url     string `yaml:"webhook_url"`
	Channel string `yaml:"channel"`
	Name    string `yaml:"name"`
}

func LoadConfig(params ...string) (*Config, error) {
	config := Config{}
	file := "./config.yaml"
	if len(params) > 0 {
		file = params[0]
	}

	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
