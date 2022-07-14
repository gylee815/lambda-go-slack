package main

import (
	"fmt"
	slack "lambda/Slack"
)

func main() {
	// tmp := slack.TestFunc()
	// fmt.Println(tmp)

	payload := slack.SlackPayload{Text: "test message", Username: "gylee"}

	rst := slack.PostMessage("https://hooks.slack.com/services/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", payload)
	fmt.Println(rst)
}