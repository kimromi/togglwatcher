package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Api   ApiConfig
	Users []User `yaml:"users"`
}

type ApiConfig struct {
	Token       string `yaml:"token"`
	DashboardId int    `yaml:"dashboardid"`
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
