package main

import (
	"fmt"

	"github.com/nlopes/slack"
)

func NotifySlack(config Notification, info NotifyInformation) {
	fields := make([]slack.AttachmentField, 0)
	if info.Status == "stopped" {
		start := slack.AttachmentField{
			Title: "Started at",
			Value: info.StartedAt,
			Short: true,
		}
		fields = append(fields, start)
		stop := slack.AttachmentField{
			Title: "Stopped at",
			Value: info.StoppedAt,
			Short: true,
		}
		fields = append(fields, stop)
	}

	emoji := map[bool]string{true: ":running::dash:", false: ":tada:"}[info.Status == "started"]
	attachment := slack.Attachment{
		Color:  "good",
		Text:   fmt.Sprintf("%s %s `%s` %s\n", UserName(info.UserID), info.Status, info.Description, emoji),
		Fields: fields,
	}
	msg := slack.WebhookMessage{
		Username:    map[bool]string{true: config.Name, false: "togglwatcher"}[config.Name != ""],
		Channel:     config.Channel,
		IconURL:     "https://toggl.com/common/images/share/favicon/favicon-32x32-664ade47fe47cfb253492cb043e3ffeb.png",
		Attachments: []slack.Attachment{attachment},
	}

	err := slack.PostWebhook(config.Url, &msg)
	if err != nil {
		fmt.Println(err)
	}
}
