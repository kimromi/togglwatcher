package notifier

import (
	"fmt"

	"../config"
	"github.com/nlopes/slack"
)

func NotifySlack(config config.Notification, info Information) {
	attachment := slack.Attachment{
		Color: "good",
		Title: UserName(info.UserID),
		Text:  fmt.Sprintf("%s `%s` :+1:\n", info.Status, info.Description),
	}
	msg := slack.WebhookMessage{
		Username:    map[bool]string{true: config.Name, false: "TogglWatcher"}[config.Name != ""],
		Channel:     config.Channel,
		IconURL:     "https://toggl.com/common/images/share/favicon/favicon-32x32-664ade47fe47cfb253492cb043e3ffeb.png",
		Attachments: []slack.Attachment{attachment},
	}

	err := slack.PostWebhook(config.Url, &msg)
	if err != nil {
		fmt.Println(err)
	}
}
