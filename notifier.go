package main

import "strconv"

type NotifyInformation struct {
	Status      string
	UserID      int
	Description string
	StartedAt   string
	StoppedAt   string
}

func Notify(info NotifyInformation) {
	c, _ := LoadConfig()
	for _, n := range c.Notifications {
		switch n.Service {
		case "slack":
			NotifySlack(n, info)
		}
	}
}

func UserName(UserID int) string {
	c, _ := LoadConfig()
	for _, user := range c.Users {
		if user.Id == UserID {
			return user.Name
		}
	}
	return strconv.Itoa(UserID)
}
