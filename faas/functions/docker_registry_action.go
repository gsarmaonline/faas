package functions

import (
	"fmt"

	"github.com/gsarmaonline/faas/faas/helpers"
	"github.com/gsarmaonline/faas/faas/intf"
)

type (
	DockerRegistryInput struct {
		Image            string `json:"image"`
		Registry         string `json:"registry"`
		RegistryUsername string `json:"registry_username"`
		RegistryPassword string `json:"registry_password"`
	}
	DockerRegistryAction struct {
		Input DockerRegistryInput
	}
)

func NewDockerRegistryAction() (dockerAction *DockerRegistryAction) {
	return &DockerRegistryAction{}
}

func (dockerAction DockerRegistryAction) GetConfig() intf.FunctionConfig {
	return intf.FunctionConfig{Name: "docker_registry"}
}

func (dockerAction *DockerRegistryAction) ParsePayload(payload intf.Payload) error {
	credManager := helpers.NewCredentialManager()
	
	processedInput := DockerRegistryInput{
		Image: payload["image"].(string),
	}

	// Optional fields
	if registry, exists := payload["registry"]; exists && registry != nil {
		processedInput.Registry = registry.(string)
	}

	// Credential fields with fallback to environment variables
	processedInput.RegistryUsername = credManager.GetCredential(payload["registry_username"], helpers.EnvDockerRegistryUsername)
	processedInput.RegistryPassword = credManager.GetCredential(payload["registry_password"], helpers.EnvDockerRegistryPassword)

	dockerAction.Input = processedInput
	return nil
}

func (dockerAction DockerRegistryAction) Validate() (err error) {
	// Check required fields
	if dockerAction.Input.Image == "" {
		return fmt.Errorf("missing required field: image")
	}

	// If registry credentials are provided, validate them
	if dockerAction.Input.Registry != "" {
		// If one credential is provided but not the other, return an error
		hasUsername := dockerAction.Input.RegistryUsername != ""
		hasPassword := dockerAction.Input.RegistryPassword != ""

		if (hasUsername && !hasPassword) || (!hasUsername && hasPassword) {
			return fmt.Errorf("when providing registry credentials, both username and password must be provided")
		}
	}

	return nil
}

func (dockerAction DockerRegistryAction) Execute() (output intf.FunctionOutput, err error) {
	var dockerExecutor *helpers.DockerExecutor

	// Create Docker executor with all parameters
	if dockerExecutor, err = helpers.NewDockerExecutor(
		dockerAction.Input.Image,
		dockerAction.Input.Registry,
		dockerAction.Input.RegistryUsername,
		dockerAction.Input.RegistryPassword,
	); err != nil {
		return
	}

	// Execute the Docker container
	if _, err = dockerExecutor.Execute(); err != nil {
		return
	}

	return
}
