package functions

import (
	"fmt"
	"log"

	"github.com/gsarmaonline/faas/faas/helpers"
	"github.com/gsarmaonline/faas/faas/intf"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type (
	SmsInput struct {
		AccountSid string `json:"account_sid"`
		AuthToken  string `json:"auth_token"`
		From       string `json:"from"`
		To         string `json:"to"`
		Body       string `json:"body"`
		MediaUrl   string `json:"media_url,omitempty"`
	}
	SmsAction struct {
		Input SmsInput
	}
)

func NewSmsAction() (smsAction *SmsAction) {
	return &SmsAction{}
}

func (smsAction SmsAction) GetConfig() intf.FunctionConfig {
	return intf.FunctionConfig{Name: "sms"}
}

func (smsAction *SmsAction) ParsePayload(payload intf.Payload) error {
	credManager := helpers.NewCredentialManager()
	
	processedInput := SmsInput{
		AccountSid: credManager.GetCredential(payload["account_sid"], helpers.EnvTwilioAccountSID),
		AuthToken:  credManager.GetCredential(payload["auth_token"], helpers.EnvTwilioAuthToken),
		From:       payload["from"].(string),
		To:         payload["to"].(string),
		Body:       payload["body"].(string),
	}

	// Optional media URL for MMS
	if mediaUrl, exists := payload["media_url"]; exists && mediaUrl != nil {
		processedInput.MediaUrl = mediaUrl.(string)
	}

	smsAction.Input = processedInput
	return nil
}

func (smsAction SmsAction) Validate() (err error) {
	if smsAction.Input.AccountSid == "" {
		return fmt.Errorf("missing required field: account_sid")
	}
	if smsAction.Input.AuthToken == "" {
		return fmt.Errorf("missing required field: auth_token")
	}
	if smsAction.Input.From == "" {
		return fmt.Errorf("missing required field: from")
	}
	if smsAction.Input.To == "" {
		return fmt.Errorf("missing required field: to")
	}
	if smsAction.Input.Body == "" && smsAction.Input.MediaUrl == "" {
		return fmt.Errorf("at least one of body or media_url must be provided")
	}
	return nil
}

func (smsAction SmsAction) Execute() (output intf.FunctionOutput, err error) {
	// Create Twilio client
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: smsAction.Input.AccountSid,
		Password: smsAction.Input.AuthToken,
	})

	// Prepare message parameters
	params := &twilioApi.CreateMessageParams{}
	params.SetFrom(smsAction.Input.From)
	params.SetTo(smsAction.Input.To)

	// Set body if provided
	if smsAction.Input.Body != "" {
		params.SetBody(smsAction.Input.Body)
	}

	// Set media URL if provided (for MMS)
	if smsAction.Input.MediaUrl != "" {
		params.SetMediaUrl([]string{smsAction.Input.MediaUrl})
	}

	// Send the message
	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		log.Printf("Failed to send SMS: %v", err)
		return nil, err
	}

	// Log successful send
	if resp.Sid != nil {
		log.Printf("SMS sent successfully. SID: %s, Status: %s", *resp.Sid, *resp.Status)
	} else {
		log.Printf("SMS sent successfully. Status: %s", *resp.Status)
	}

	return nil, nil
}
