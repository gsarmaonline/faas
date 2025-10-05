package functions

import (
	"context"
	"fmt"
	"log"

	"github.com/gsarmaonline/faas/faas/intf"
	"github.com/slack-go/slack"
)

type (
	SlackInput struct {
		ApiToken  string `json:"api_token"`
		Message   string `json:"message"`
		ChannelID string `json:"channel_id"`
	}
	Slack struct {
		Input SlackInput
	}
)

func NewSlack() (slackFunc *Slack) {
	return &Slack{Input: SlackInput{}}
}

func (slackFunc Slack) GetConfig() intf.FunctionConfig {
	return intf.FunctionConfig{Name: "slack"}
}

func (slackFunc Slack) ParsePayload(payload map[string]interface{}) (SlackInput, error) {
	processedSlackInput := SlackInput{
		ApiToken:  payload["api_token"].(string),
		Message:   payload["message"].(string),
		ChannelID: payload["channel_id"].(string),
	}
	return processedSlackInput, nil
}

func (slackFunc Slack) Validate() (err error) {
	if slackFunc.Input.ApiToken == "" || slackFunc.Input.ChannelID == "" || slackFunc.Input.Message == "" {
		err = fmt.Errorf("missing required fields in slack input")
		return
	}
	return
}

func (slackFunc Slack) Execute() (output Output, err error) {
	client := slack.New(slackFunc.Input.ApiToken)

	if _, _, err = client.PostMessageContext(context.Background(),
		slackFunc.Input.ChannelID,
		slack.MsgOptionText(slackFunc.Input.Message, false),
	); err != nil {
		log.Println("Failed to post slack message", err)
	}
	return
}
