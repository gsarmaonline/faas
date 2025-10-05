package helpers

import (
	"os"
	"testing"
)

func TestCredentialManager_GetCredential(t *testing.T) {
	cm := NewCredentialManager()

	tests := []struct {
		name         string
		payloadValue interface{}
		envVarName   string
		envVarValue  string
		expected     string
	}{
		{
			name:         "payload value takes priority",
			payloadValue: "payload_value",
			envVarName:   "TEST_ENV_VAR",
			envVarValue:  "env_value",
			expected:     "payload_value",
		},
		{
			name:         "falls back to env var when payload is nil",
			payloadValue: nil,
			envVarName:   "TEST_ENV_VAR",
			envVarValue:  "env_value",
			expected:     "env_value",
		},
		{
			name:         "falls back to env var when payload is empty string",
			payloadValue: "",
			envVarName:   "TEST_ENV_VAR",
			envVarValue:  "env_value",
			expected:     "env_value",
		},
		{
			name:         "returns empty when neither is available",
			payloadValue: nil,
			envVarName:   "NON_EXISTENT_VAR",
			envVarValue:  "",
			expected:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment variable if provided
			if tt.envVarValue != "" {
				os.Setenv(tt.envVarName, tt.envVarValue)
				defer os.Unsetenv(tt.envVarName)
			}

			result := cm.GetCredential(tt.payloadValue, tt.envVarName)
			if result != tt.expected {
				t.Errorf("GetCredential() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCredentialManager_GetRequiredCredential(t *testing.T) {
	cm := NewCredentialManager()

	t.Run("returns value when available", func(t *testing.T) {
		os.Setenv("TEST_REQUIRED_VAR", "test_value")
		defer os.Unsetenv("TEST_REQUIRED_VAR")

		result, err := cm.GetRequiredCredential(nil, "TEST_REQUIRED_VAR", "test_field")
		if err != nil {
			t.Errorf("GetRequiredCredential() error = %v, want nil", err)
		}
		if result != "test_value" {
			t.Errorf("GetRequiredCredential() = %v, want %v", result, "test_value")
		}
	})

	t.Run("returns error when not available", func(t *testing.T) {
		_, err := cm.GetRequiredCredential(nil, "NON_EXISTENT_VAR", "test_field")
		if err == nil {
			t.Error("GetRequiredCredential() error = nil, want error")
		}
		expectedErrMsg := "missing required credential: test_field (provide in payload or set NON_EXISTENT_VAR environment variable)"
		if err.Error() != expectedErrMsg {
			t.Errorf("GetRequiredCredential() error = %v, want %v", err.Error(), expectedErrMsg)
		}
	})
}

func TestValidateEnvironmentVars(t *testing.T) {
	// Set up test environment
	os.Setenv("EXISTING_VAR", "value")
	defer os.Unsetenv("EXISTING_VAR")

	tests := []struct {
		name         string
		requiredVars []string
		expected     []string
	}{
		{
			name:         "all variables exist",
			requiredVars: []string{"EXISTING_VAR"},
			expected:     []string{},
		},
		{
			name:         "some variables missing",
			requiredVars: []string{"EXISTING_VAR", "MISSING_VAR"},
			expected:     []string{"MISSING_VAR"},
		},
		{
			name:         "all variables missing",
			requiredVars: []string{"MISSING_VAR_1", "MISSING_VAR_2"},
			expected:     []string{"MISSING_VAR_1", "MISSING_VAR_2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateEnvironmentVars(tt.requiredVars)
			if len(result) != len(tt.expected) {
				t.Errorf("ValidateEnvironmentVars() = %v, want %v", result, tt.expected)
				return
			}
			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("ValidateEnvironmentVars()[%d] = %v, want %v", i, result[i], expected)
				}
			}
		})
	}
}
