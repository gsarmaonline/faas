package functions

import (
	"testing"

	"github.com/gsarmaonline/faas/faas/intf"
)

func TestSmsAction_GetConfig(t *testing.T) {
	smsAction := NewSmsAction()
	config := smsAction.GetConfig()

	if config.Name != "sms" {
		t.Errorf("Expected config name 'sms', got '%s'", config.Name)
	}
}

func TestSmsAction_ParsePayload(t *testing.T) {
	tests := []struct {
		name    string
		payload intf.Payload
		want    SmsInput
	}{
		{
			name: "payload with required fields only",
			payload: intf.Payload{
				"account_sid": "AC123456789",
				"auth_token":  "test-auth-token",
				"from":        "+1234567890",
				"to":          "+0987654321",
				"body":        "Hello from FAAS!",
			},
			want: SmsInput{
				AccountSid: "AC123456789",
				AuthToken:  "test-auth-token",
				From:       "+1234567890",
				To:         "+0987654321",
				Body:       "Hello from FAAS!",
			},
		},
		{
			name: "payload with all fields",
			payload: intf.Payload{
				"account_sid": "AC123456789",
				"auth_token":  "test-auth-token",
				"from":        "+1234567890",
				"to":          "+0987654321",
				"body":        "Hello from FAAS!",
				"media_url":   "https://example.com/image.jpg",
			},
			want: SmsInput{
				AccountSid: "AC123456789",
				AuthToken:  "test-auth-token",
				From:       "+1234567890",
				To:         "+0987654321",
				Body:       "Hello from FAAS!",
				MediaUrl:   "https://example.com/image.jpg",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			smsAction := NewSmsAction()
			err := smsAction.ParsePayload(tt.payload)

			if err != nil {
				t.Errorf("ParsePayload() error = %v", err)
				return
			}

			if smsAction.Input.AccountSid != tt.want.AccountSid {
				t.Errorf("AccountSid = %v, want %v", smsAction.Input.AccountSid, tt.want.AccountSid)
			}
			if smsAction.Input.AuthToken != tt.want.AuthToken {
				t.Errorf("AuthToken = %v, want %v", smsAction.Input.AuthToken, tt.want.AuthToken)
			}
			if smsAction.Input.From != tt.want.From {
				t.Errorf("From = %v, want %v", smsAction.Input.From, tt.want.From)
			}
			if smsAction.Input.To != tt.want.To {
				t.Errorf("To = %v, want %v", smsAction.Input.To, tt.want.To)
			}
			if smsAction.Input.Body != tt.want.Body {
				t.Errorf("Body = %v, want %v", smsAction.Input.Body, tt.want.Body)
			}
			if smsAction.Input.MediaUrl != tt.want.MediaUrl {
				t.Errorf("MediaUrl = %v, want %v", smsAction.Input.MediaUrl, tt.want.MediaUrl)
			}
		})
	}
}

func TestSmsAction_Validate(t *testing.T) {
	tests := []struct {
		name    string
		input   SmsInput
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid SMS with body",
			input: SmsInput{
				AccountSid: "AC123456789",
				AuthToken:  "test-auth-token",
				From:       "+1234567890",
				To:         "+0987654321",
				Body:       "Hello from FAAS!",
			},
			wantErr: false,
		},
		{
			name: "valid MMS with media URL only",
			input: SmsInput{
				AccountSid: "AC123456789",
				AuthToken:  "test-auth-token",
				From:       "+1234567890",
				To:         "+0987654321",
				MediaUrl:   "https://example.com/image.jpg",
			},
			wantErr: false,
		},
		{
			name: "valid with both body and media URL",
			input: SmsInput{
				AccountSid: "AC123456789",
				AuthToken:  "test-auth-token",
				From:       "+1234567890",
				To:         "+0987654321",
				Body:       "Check this out!",
				MediaUrl:   "https://example.com/image.jpg",
			},
			wantErr: false,
		},
		{
			name: "missing account_sid",
			input: SmsInput{
				AuthToken: "test-auth-token",
				From:      "+1234567890",
				To:        "+0987654321",
				Body:      "Hello from FAAS!",
			},
			wantErr: true,
			errMsg:  "missing required field: account_sid",
		},
		{
			name: "missing auth_token",
			input: SmsInput{
				AccountSid: "AC123456789",
				From:       "+1234567890",
				To:         "+0987654321",
				Body:       "Hello from FAAS!",
			},
			wantErr: true,
			errMsg:  "missing required field: auth_token",
		},
		{
			name: "missing from",
			input: SmsInput{
				AccountSid: "AC123456789",
				AuthToken:  "test-auth-token",
				To:         "+0987654321",
				Body:       "Hello from FAAS!",
			},
			wantErr: true,
			errMsg:  "missing required field: from",
		},
		{
			name: "missing to",
			input: SmsInput{
				AccountSid: "AC123456789",
				AuthToken:  "test-auth-token",
				From:       "+1234567890",
				Body:       "Hello from FAAS!",
			},
			wantErr: true,
			errMsg:  "missing required field: to",
		},
		{
			name: "missing both body and media_url",
			input: SmsInput{
				AccountSid: "AC123456789",
				AuthToken:  "test-auth-token",
				From:       "+1234567890",
				To:         "+0987654321",
			},
			wantErr: true,
			errMsg:  "at least one of body or media_url must be provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			smsAction := NewSmsAction()
			smsAction.Input = tt.input
			err := smsAction.Validate()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Validate() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("Validate() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestSmsAction_Execute_Integration(t *testing.T) {
	// This test requires real Twilio credentials and should be run with integration tests
	// For now, we'll skip it unless TWILIO_ACCOUNT_SID and TWILIO_AUTH_TOKEN environment variables are set
	t.Skip("Skipping integration test - requires TWILIO_ACCOUNT_SID and TWILIO_AUTH_TOKEN environment variables")

	// Uncomment and modify for actual integration testing:
	/*
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	from := os.Getenv("TWILIO_FROM_NUMBER")
	to := os.Getenv("TWILIO_TO_NUMBER")

	if accountSid == "" || authToken == "" || from == "" || to == "" {
		t.Skip("Required Twilio environment variables not set")
	}

	smsAction := NewSmsAction()
	smsAction.Input = SmsInput{
		AccountSid: accountSid,
		AuthToken:  authToken,
		From:       from,
		To:         to,
		Body:       "Test SMS from FAAS framework",
	}

	err := smsAction.Validate()
	if err != nil {
		t.Errorf("Validate() error = %v", err)
		return
	}

	_, err = smsAction.Execute()
	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}
	*/
}
