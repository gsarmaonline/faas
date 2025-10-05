package functions

import (
	"fmt"

	"github.com/gsarmaonline/faas/faas/intf"
)

type (
	GithubInput struct {
		Repository string `json:"repository"`
		Action     string `json:"action"`
		Token      string `json:"token"`
	}
	GithubAction struct {
		Input GithubInput
	}
)

func NewGithubAction() (githubAction *GithubAction) {
	return &GithubAction{}
}

func (githubAction GithubAction) GetConfig() intf.FunctionConfig {
	return intf.FunctionConfig{Name: "github"}
}

func (githubAction *GithubAction) ParsePayload(payload intf.Payload) error {
	processedInput := GithubInput{
		Repository: payload["repository"].(string),
		Action:     payload["action"].(string),
	}

	// Optional token field
	if token, exists := payload["token"]; exists && token != nil {
		processedInput.Token = token.(string)
	}

	githubAction.Input = processedInput
	return nil
}

func (githubAction GithubAction) Validate() (err error) {
	if githubAction.Input.Repository == "" {
		return fmt.Errorf("missing required field: repository")
	}
	if githubAction.Input.Action == "" {
		return fmt.Errorf("missing required field: action")
	}
	return nil
}

func (githubAction GithubAction) Execute() (output intf.FunctionOutput, err error) {
	// TODO: Implement actual GitHub API calls here
	// For now, this is a placeholder implementation
	return
}
