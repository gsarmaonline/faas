package functions

import (
	"testing"

	"github.com/gsarmaonline/faas/faas/intf"
)

func TestEmailAction_GetConfig(t *testing.T) {
	emailAction := NewEmailAction()
	config := emailAction.GetConfig()

	if config.Name != "email" {
		t.Errorf("Expected config name 'email', got '%s'", config.Name)
	}
}

func TestEmailAction_ParsePayload(t *testing.T) {
	tests := []struct {
		name    string
		payload intf.Payload
		want    EmailInput
	}{
		{
			name: "payload with required fields only",
			payload: intf.Payload{
				"api_key":    "test-api-key",
				"from_email": "sender@example.com",
				"to_email":   "recipient@example.com",
				"subject":    "Test Subject",
				"plain_text": "Test message body",
			},
			want: EmailInput{
				ApiKey:    "test-api-key",
				FromEmail: "sender@example.com",
				ToEmail:   "recipient@example.com",
				Subject:   "Test Subject",
				PlainText: "Test message body",
			},
		},
		{
			name: "payload with all fields",
			payload: intf.Payload{
				"api_key":    "test-api-key",
				"from_email": "sender@example.com",
				"from_name":  "Test Sender",
				"to_email":   "recipient@example.com",
				"to_name":    "Test Recipient",
				"subject":    "Test Subject",
				"plain_text": "Test message body",
				"html_text":  "<h1>Test HTML</h1>",
			},
			want: EmailInput{
				ApiKey:    "test-api-key",
				FromEmail: "sender@example.com",
				FromName:  "Test Sender",
				ToEmail:   "recipient@example.com",
				ToName:    "Test Recipient",
				Subject:   "Test Subject",
				PlainText: "Test message body",
				HtmlText:  "<h1>Test HTML</h1>",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			emailAction := NewEmailAction()
			err := emailAction.ParsePayload(tt.payload)

			if err != nil {
				t.Errorf("ParsePayload() error = %v", err)
				return
			}

			if emailAction.Input.ApiKey != tt.want.ApiKey {
				t.Errorf("ApiKey = %v, want %v", emailAction.Input.ApiKey, tt.want.ApiKey)
			}
			if emailAction.Input.FromEmail != tt.want.FromEmail {
				t.Errorf("FromEmail = %v, want %v", emailAction.Input.FromEmail, tt.want.FromEmail)
			}
			if emailAction.Input.FromName != tt.want.FromName {
				t.Errorf("FromName = %v, want %v", emailAction.Input.FromName, tt.want.FromName)
			}
			if emailAction.Input.ToEmail != tt.want.ToEmail {
				t.Errorf("ToEmail = %v, want %v", emailAction.Input.ToEmail, tt.want.ToEmail)
			}
			if emailAction.Input.ToName != tt.want.ToName {
				t.Errorf("ToName = %v, want %v", emailAction.Input.ToName, tt.want.ToName)
			}
			if emailAction.Input.Subject != tt.want.Subject {
				t.Errorf("Subject = %v, want %v", emailAction.Input.Subject, tt.want.Subject)
			}
			if emailAction.Input.PlainText != tt.want.PlainText {
				t.Errorf("PlainText = %v, want %v", emailAction.Input.PlainText, tt.want.PlainText)
			}
			if emailAction.Input.HtmlText != tt.want.HtmlText {
				t.Errorf("HtmlText = %v, want %v", emailAction.Input.HtmlText, tt.want.HtmlText)
			}
		})
	}
}

func TestEmailAction_Validate(t *testing.T) {
	tests := []struct {
		name    string
		input   EmailInput
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid with plain text only",
			input: EmailInput{
				ApiKey:    "test-api-key",
				FromEmail: "sender@example.com",
				ToEmail:   "recipient@example.com",
				Subject:   "Test Subject",
				PlainText: "Test message",
			},
			wantErr: false,
		},
		{
			name: "valid with HTML text only",
			input: EmailInput{
				ApiKey:    "test-api-key",
				FromEmail: "sender@example.com",
				ToEmail:   "recipient@example.com",
				Subject:   "Test Subject",
				HtmlText:  "<h1>Test HTML</h1>",
			},
			wantErr: false,
		},
		{
			name: "valid with both plain and HTML text",
			input: EmailInput{
				ApiKey:    "test-api-key",
				FromEmail: "sender@example.com",
				ToEmail:   "recipient@example.com",
				Subject:   "Test Subject",
				PlainText: "Test message",
				HtmlText:  "<h1>Test HTML</h1>",
			},
			wantErr: false,
		},
		{
			name: "missing api_key",
			input: EmailInput{
				FromEmail: "sender@example.com",
				ToEmail:   "recipient@example.com",
				Subject:   "Test Subject",
				PlainText: "Test message",
			},
			wantErr: true,
			errMsg:  "missing required field: api_key",
		},
		{
			name: "missing from_email",
			input: EmailInput{
				ApiKey:    "test-api-key",
				ToEmail:   "recipient@example.com",
				Subject:   "Test Subject",
				PlainText: "Test message",
			},
			wantErr: true,
			errMsg:  "missing required field: from_email",
		},
		{
			name: "missing to_email",
			input: EmailInput{
				ApiKey:    "test-api-key",
				FromEmail: "sender@example.com",
				Subject:   "Test Subject",
				PlainText: "Test message",
			},
			wantErr: true,
			errMsg:  "missing required field: to_email",
		},
		{
			name: "missing subject",
			input: EmailInput{
				ApiKey:    "test-api-key",
				FromEmail: "sender@example.com",
				ToEmail:   "recipient@example.com",
				PlainText: "Test message",
			},
			wantErr: true,
			errMsg:  "missing required field: subject",
		},
		{
			name: "missing both plain_text and html_text",
			input: EmailInput{
				ApiKey:    "test-api-key",
				FromEmail: "sender@example.com",
				ToEmail:   "recipient@example.com",
				Subject:   "Test Subject",
			},
			wantErr: true,
			errMsg:  "at least one of plain_text or html_text must be provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			emailAction := NewEmailAction()
			emailAction.Input = tt.input
			err := emailAction.Validate()

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

func TestEmailAction_Execute_Integration(t *testing.T) {
	// This test requires a real SendGrid API key and should be run with integration tests
	// For now, we'll skip it unless SENDGRID_API_KEY environment variable is set
	t.Skip("Skipping integration test - requires SENDGRID_API_KEY environment variable")

	// Uncomment and modify for actual integration testing:
	/*
		apiKey := os.Getenv("SENDGRID_API_KEY")
		if apiKey == "" {
			t.Skip("SENDGRID_API_KEY environment variable not set")
		}

		emailAction := NewEmailAction()
		emailAction.Input = EmailInput{
			ApiKey:    apiKey,
			FromEmail: "test@example.com",
			ToEmail:   "recipient@example.com",
			Subject:   "Test Email",
			PlainText: "This is a test email from the FAAS framework",
		}

		err := emailAction.Validate()
		if err != nil {
			t.Errorf("Validate() error = %v", err)
			return
		}

		_, err = emailAction.Execute()
		if err != nil {
			t.Errorf("Execute() error = %v", err)
		}
	*/
}
