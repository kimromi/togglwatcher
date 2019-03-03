package notifier

import (
	"strconv"

	"../config"
)

type Information struct {
	Status      string
	UserID      int
	Description string
	StartedAt   string
	StoppedAt   string
}

func Notify(info Information) {
	for _, n := range config.LoadConfig().Notifications {
		switch n.Service {
		case "slack":
			NotifySlack(n, info)
		}
	}
}

func UserName(UserID int) string {
	for _, user := range config.LoadConfig().Users {
		if user.Id == UserID {
			return user.Name
		}
	}
	return strconv.Itoa(UserID)
}
