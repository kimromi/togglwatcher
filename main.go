package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "togglwatcher"
	app.Version = "0.1.0"

	app.Action = func(c *cli.Context) error {
		dashboard := FetchDashboard()

		// iroiro
		fmt.Println(dashboard)

		return nil
	}

	app.Run(os.Args)
}

type Config struct {
	Api ApiConfig
}

type ApiConfig struct {
	Token       string `toml:"token"`
	WorkspaceId int    `toml:"workspaceid"`
}

func LoadConfig() Config {
	var config Config
	_, err := toml.DecodeFile("config.toml", &config)
	if err != nil {
		panic(err)
	}
	return config
}

type Dashboard struct {
	MostActiveUser []struct {
		UserID   int `json:"user_id"`
		Duration int `json:"duration"`
	} `json:"most_active_user"`
	Activity []struct {
		UserID      int         `json:"user_id"`
		ProjectID   int         `json:"project_id"`
		Duration    int         `json:"duration"`
		Description string      `json:"description"`
		Stop        interface{} `json:"stop"`
		Tid         interface{} `json:"tid"`
	} `json:"activity"`
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
