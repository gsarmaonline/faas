package functions

import (
	"testing"

	"github.com/gsarmaonline/faas/faas/intf"
)

func TestLoggerAction_GetConfig(t *testing.T) {
	loggerAction := NewLoggerAction()
	config := loggerAction.GetConfig()

	if config.Name != "logger" {
		t.Errorf("Expected config name 'logger', got '%s'", config.Name)
	}
}

func TestLoggerAction_ParsePayload(t *testing.T) {
	tests := []struct {
		name    string
		payload intf.Payload
		want    LoggerInput
	}{
		{
			name: "payload with message",
			payload: intf.Payload{
				"message": "Hello from FAAS Logger!",
			},
			want: LoggerInput{
				Message: "Hello from FAAS Logger!",
			},
		},
		{
			name: "payload without message (empty)",
			payload: intf.Payload{},
			want: LoggerInput{
				Message: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loggerAction := NewLoggerAction()
			err := loggerAction.ParsePayload(tt.payload)

			if err != nil {
				t.Errorf("ParsePayload() error = %v", err)
				return
			}

			if loggerAction.Input.Message != tt.want.Message {
				t.Errorf("Message = %v, want %v", loggerAction.Input.Message, tt.want.Message)
			}
		})
	}
}

func TestLoggerAction_Validate(t *testing.T) {
	tests := []struct {
		name      string
		input     LoggerInput
		wantError bool
	}{
		{
			name: "valid input with message",
			input: LoggerInput{
				Message: "Hello from FAAS Logger!",
			},
			wantError: false,
		},
		{
			name: "valid input without message (message is optional)",
			input: LoggerInput{
				Message: "",
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loggerAction := LoggerAction{Input: tt.input}
			err := loggerAction.Validate()

			if (err != nil) != tt.wantError {
				t.Errorf("Validate() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}
