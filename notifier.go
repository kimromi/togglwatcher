package main

type NotifyInformation struct {
	Status      string
	User        User
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
