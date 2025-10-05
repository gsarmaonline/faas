package functions

import (
	"testing"

	"github.com/gsarmaonline/faas/faas/intf"
)

func TestDockerAction_GetConfig(t *testing.T) {
	dockerAction := NewDockerRegistryAction()
	config := dockerAction.GetConfig()

	if config.Name != "docker_registry" {
		t.Errorf("Expected config name 'docker_registry', got '%s'", config.Name)
	}
}

func TestDockerAction_ParsePayload(t *testing.T) {
	tests := []struct {
		name    string
		payload intf.Payload
		want    DockerRegistryInput
	}{
		{
			name: "payload with only image",
			payload: intf.Payload{
				"image": "nginx",
			},
			want: DockerRegistryInput{
				Image: "nginx",
			},
		},
		{
			name: "payload with all fields",
			payload: intf.Payload{
				"image":             "nginx",
				"registry":          "registry.example.com",
				"registry_username": "user",
				"registry_password": "pass",
			},
			want: DockerRegistryInput{
				Image:            "nginx",
				Registry:         "registry.example.com",
				RegistryUsername: "user",
				RegistryPassword: "pass",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dockerAction := NewDockerRegistryAction()
			err := dockerAction.ParsePayload(tt.payload)

			if err != nil {
				t.Errorf("ParsePayload() error = %v", err)
				return
			}

			if dockerAction.Input.Image != tt.want.Image {
				t.Errorf("Image = %v, want %v", dockerAction.Input.Image, tt.want.Image)
			}
			if dockerAction.Input.Registry != tt.want.Registry {
				t.Errorf("Registry = %v, want %v", dockerAction.Input.Registry, tt.want.Registry)
			}
			if dockerAction.Input.RegistryUsername != tt.want.RegistryUsername {
				t.Errorf("RegistryUsername = %v, want %v", dockerAction.Input.RegistryUsername, tt.want.RegistryUsername)
			}
			if dockerAction.Input.RegistryPassword != tt.want.RegistryPassword {
				t.Errorf("RegistryPassword = %v, want %v", dockerAction.Input.RegistryPassword, tt.want.RegistryPassword)
			}
		})
	}
}

func TestDockerAction_Validate(t *testing.T) {
	tests := []struct {
		name    string
		input   DockerRegistryInput
		wantErr bool
	}{
		{
			name: "valid with only image",
			input: DockerRegistryInput{
				Image: "nginx",
			},
			wantErr: false,
		},
		{
			name: "valid with image and registry",
			input: DockerRegistryInput{
				Image:    "nginx",
				Registry: "registry.example.com",
			},
			wantErr: false,
		},
		{
			name: "valid with all fields",
			input: DockerRegistryInput{
				Image:            "nginx",
				Registry:         "registry.example.com",
				RegistryUsername: "user",
				RegistryPassword: "pass",
			},
			wantErr: false,
		},
		{
			name: "missing image",
			input: DockerRegistryInput{
				Registry: "registry.example.com",
			},
			wantErr: true,
		},
		{
			name: "username without password",
			input: DockerRegistryInput{
				Image:            "nginx",
				Registry:         "registry.example.com",
				RegistryUsername: "user",
			},
			wantErr: true,
		},
		{
			name: "password without username",
			input: DockerRegistryInput{
				Image:            "nginx",
				Registry:         "registry.example.com",
				RegistryPassword: "pass",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dockerAction := NewDockerRegistryAction()
			dockerAction.Input = tt.input
			err := dockerAction.Validate()

			if (err != nil) != tt.wantErr {
				t.Errorf("DockerAction.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
