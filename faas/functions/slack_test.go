package functions

import (
	"testing"

	"github.com/gsarmaonline/faas/faas/intf"
)

func TestSlack_GetConfig(t *testing.T) {
	slack := NewSlack()
	config := slack.GetConfig()

	if config.Name != "slack" {
		t.Errorf("Expected config name 'slack', got '%s'", config.Name)
	}
}

func TestSlack_ParsePayload(t *testing.T) {
	tests := []struct {
		name    string
		payload intf.Payload
		want    SlackInput
		envVars map[string]string
	}{
		{
			name: "secure payload without api_token (using environment variables)",
			payload: intf.Payload{
				"message":    "Hello from FAAS!",
				"channel_id": "C1234567890",
			},
			envVars: map[string]string{
				"SLACK_API_TOKEN": "xoxb-test-token",
			},
			want: SlackInput{
				ApiToken:  "xoxb-test-token",
				Message:   "Hello from FAAS!",
				ChannelID: "C1234567890",
			},
		},
		{
			name: "payload override api_token (testing override functionality)",
			payload: intf.Payload{
				"api_token":  "xoxb-override-token",
				"message":    "Hello from FAAS!",
				"channel_id": "C1234567890",
			},
			envVars: map[string]string{
				"SLACK_API_TOKEN": "xoxb-env-token",
			},
			want: SlackInput{
				ApiToken:  "xoxb-override-token",
				Message:   "Hello from FAAS!",
				ChannelID: "C1234567890",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment variables for this test
			for key, value := range tt.envVars {
				t.Setenv(key, value)
			}

			slack := NewSlack()
			err := slack.ParsePayload(tt.payload)

			if err != nil {
				t.Errorf("ParsePayload() error = %v", err)
				return
			}

			if slack.Input.ApiToken != tt.want.ApiToken {
				t.Errorf("ApiToken = %v, want %v", slack.Input.ApiToken, tt.want.ApiToken)
			}
			if slack.Input.Message != tt.want.Message {
				t.Errorf("Message = %v, want %v", slack.Input.Message, tt.want.Message)
			}
			if slack.Input.ChannelID != tt.want.ChannelID {
				t.Errorf("ChannelID = %v, want %v", slack.Input.ChannelID, tt.want.ChannelID)
			}
		})
	}
}

func TestSlack_Validate(t *testing.T) {
	tests := []struct {
		name      string
		input     SlackInput
		wantError bool
	}{
		{
			name: "valid input",
			input: SlackInput{
				ApiToken:  "xoxb-test-token",
				Message:   "Hello from FAAS!",
				ChannelID: "C1234567890",
			},
			wantError: false,
		},
		{
			name: "missing api_token",
			input: SlackInput{
				Message:   "Hello from FAAS!",
				ChannelID: "C1234567890",
			},
			wantError: true,
		},
		{
			name: "missing message",
			input: SlackInput{
				ApiToken:  "xoxb-test-token",
				ChannelID: "C1234567890",
			},
			wantError: true,
		},
		{
			name: "missing channel_id",
			input: SlackInput{
				ApiToken: "xoxb-test-token",
				Message:  "Hello from FAAS!",
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slack := Slack{Input: tt.input}
			err := slack.Validate()

			if (err != nil) != tt.wantError {
				t.Errorf("Validate() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}
