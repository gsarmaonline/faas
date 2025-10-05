package functions

import (
	"testing"

	"github.com/gsarmaonline/faas/faas/intf"
)

func TestGithubAction_GetConfig(t *testing.T) {
	githubAction := NewGithubAction()
	config := githubAction.GetConfig()

	if config.Name != "github" {
		t.Errorf("Expected config name 'github', got '%s'", config.Name)
	}
}

func TestGithubAction_ParsePayload(t *testing.T) {
	tests := []struct {
		name    string
		payload intf.Payload
		want    GithubInput
		envVars map[string]string
	}{
		{
			name: "secure payload without token (using environment variables)",
			payload: intf.Payload{
				"repository": "owner/repo",
				"action":     "create_issue",
			},
			envVars: map[string]string{
				"GITHUB_TOKEN": "ghp_test_token",
			},
			want: GithubInput{
				Repository: "owner/repo",
				Action:     "create_issue",
				Token:      "ghp_test_token",
			},
		},
		{
			name: "payload override token (testing override functionality)",
			payload: intf.Payload{
				"repository": "owner/repo",
				"action":     "create_issue",
				"token":      "ghp_override_token",
			},
			envVars: map[string]string{
				"GITHUB_TOKEN": "ghp_env_token",
			},
			want: GithubInput{
				Repository: "owner/repo",
				Action:     "create_issue",
				Token:      "ghp_override_token",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment variables for this test
			for key, value := range tt.envVars {
				t.Setenv(key, value)
			}

			githubAction := NewGithubAction()
			err := githubAction.ParsePayload(tt.payload)

			if err != nil {
				t.Errorf("ParsePayload() error = %v", err)
				return
			}

			if githubAction.Input.Repository != tt.want.Repository {
				t.Errorf("Repository = %v, want %v", githubAction.Input.Repository, tt.want.Repository)
			}
			if githubAction.Input.Action != tt.want.Action {
				t.Errorf("Action = %v, want %v", githubAction.Input.Action, tt.want.Action)
			}
			if githubAction.Input.Token != tt.want.Token {
				t.Errorf("Token = %v, want %v", githubAction.Input.Token, tt.want.Token)
			}
		})
	}
}

func TestGithubAction_Validate(t *testing.T) {
	tests := []struct {
		name      string
		input     GithubInput
		wantError bool
	}{
		{
			name: "valid input with token",
			input: GithubInput{
				Repository: "owner/repo",
				Action:     "create_issue",
				Token:      "ghp_test_token",
			},
			wantError: false,
		},
		{
			name: "valid input without token (token is optional)",
			input: GithubInput{
				Repository: "owner/repo",
				Action:     "create_issue",
			},
			wantError: false,
		},
		{
			name: "missing repository",
			input: GithubInput{
				Action: "create_issue",
				Token:  "ghp_test_token",
			},
			wantError: true,
		},
		{
			name: "missing action",
			input: GithubInput{
				Repository: "owner/repo",
				Token:      "ghp_test_token",
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			githubAction := GithubAction{Input: tt.input}
			err := githubAction.Validate()

			if (err != nil) != tt.wantError {
				t.Errorf("Validate() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}
