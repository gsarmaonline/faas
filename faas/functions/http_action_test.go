package functions

import (
	"testing"

	"github.com/gsarmaonline/faas/faas/intf"
)

func TestHttpAction_GetConfig(t *testing.T) {
	httpAction := NewHttpAction()
	config := httpAction.GetConfig()

	if config.Name != "http" {
		t.Errorf("Expected config name 'http', got '%s'", config.Name)
	}
}

func TestHttpAction_ParsePayload(t *testing.T) {
	tests := []struct {
		name    string
		payload intf.Payload
		want    HttpInput
	}{
		{
			name: "GET request without body",
			payload: intf.Payload{
				"url":    "https://api.example.com/users",
				"method": "GET",
			},
			want: HttpInput{
				Url:    "https://api.example.com/users",
				Method: GetHttpMethod,
			},
		},
		{
			name: "POST request with body",
			payload: intf.Payload{
				"url":          "https://api.example.com/users",
				"method":       "POST",
				"request_body": map[string]interface{}{"name": "John", "email": "john@example.com"},
			},
			want: HttpInput{
				Url:         "https://api.example.com/users",
				Method:      PostHttpMethod,
				RequestBody: map[string]interface{}{"name": "John", "email": "john@example.com"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpAction := NewHttpAction()
			err := httpAction.ParsePayload(tt.payload)

			if err != nil {
				t.Errorf("ParsePayload() error = %v", err)
				return
			}

			if httpAction.Input.Url != tt.want.Url {
				t.Errorf("Url = %v, want %v", httpAction.Input.Url, tt.want.Url)
			}
			if httpAction.Input.Method != tt.want.Method {
				t.Errorf("Method = %v, want %v", httpAction.Input.Method, tt.want.Method)
			}
		})
	}
}

func TestHttpAction_Validate(t *testing.T) {
	tests := []struct {
		name      string
		input     HttpInput
		wantError bool
	}{
		{
			name: "valid GET request",
			input: HttpInput{
				Url:    "https://api.example.com/users",
				Method: GetHttpMethod,
			},
			wantError: false,
		},
		{
			name: "valid POST request",
			input: HttpInput{
				Url:         "https://api.example.com/users",
				Method:      PostHttpMethod,
				RequestBody: map[string]interface{}{"name": "John"},
			},
			wantError: false,
		},
		{
			name: "missing URL",
			input: HttpInput{
				Method: GetHttpMethod,
			},
			wantError: true,
		},
		{
			name: "missing method",
			input: HttpInput{
				Url: "https://api.example.com/users",
			},
			wantError: true,
		},
		{
			name: "valid with custom method (validation doesn't check method validity)",
			input: HttpInput{
				Url:    "https://api.example.com/users",
				Method: HttpMethodT("CUSTOM"),
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpAction := HttpAction{Input: tt.input}
			err := httpAction.Validate()

			if (err != nil) != tt.wantError {
				t.Errorf("Validate() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}
