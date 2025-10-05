package helpers

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// CredentialManager handles reading credentials from environment variables
// with optional payload overrides
type CredentialManager struct{}

// GetCredential reads a credential value with the following priority:
// 1. Payload value (if provided)
// 2. Environment variable
// 3. Empty string (if neither is available)
func (cm *CredentialManager) GetCredential(payloadValue interface{}, envVarName string) string {
	// First check if payload has the value
	if payloadValue != nil {
		if str, ok := payloadValue.(string); ok && str != "" {
			return str
		}
	}

	// Fall back to environment variable
	return os.Getenv(envVarName)
}

// GetRequiredCredential is like GetCredential but returns an error if no value is found
func (cm *CredentialManager) GetRequiredCredential(payloadValue interface{}, envVarName, fieldName string) (string, error) {
	value := cm.GetCredential(payloadValue, envVarName)
	if value == "" {
		return "", fmt.Errorf("missing required credential: %s (provide in payload or set %s environment variable)", fieldName, envVarName)
	}
	return value, nil
}

// Standard environment variable names
const (
	// Email/SendGrid
	EnvSendGridAPIKey = "SENDGRID_API_KEY"

	// Slack
	EnvSlackAPIToken = "SLACK_API_TOKEN"

	// SMS/Twilio
	EnvTwilioAccountSID = "TWILIO_ACCOUNT_SID"
	EnvTwilioAuthToken  = "TWILIO_AUTH_TOKEN"

	// Docker Registry
	EnvDockerRegistryUsername = "DOCKER_REGISTRY_USERNAME"
	EnvDockerRegistryPassword = "DOCKER_REGISTRY_PASSWORD"

	// GitHub
	EnvGitHubToken = "GITHUB_TOKEN"
)

// NewCredentialManager creates a new credential manager instance
func NewCredentialManager() *CredentialManager {
	return &CredentialManager{}
}

// ValidateEnvironmentVars checks if required environment variables are set
// and returns a list of missing ones
func ValidateEnvironmentVars(requiredVars []string) []string {
	var missing []string
	for _, envVar := range requiredVars {
		if os.Getenv(envVar) == "" {
			missing = append(missing, envVar)
		}
	}
	return missing
}

// LoadEnvFile loads environment variables from a .env file (if it exists)
// This is useful for development
func LoadEnvFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		// File doesn't exist, which is okay
		return nil
	}
	defer file.Close()

	// Simple .env file parser
	// In production, you might want to use a library like godotenv
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// Remove quotes if present
			if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
				(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
				value = value[1 : len(value)-1]
			}
			os.Setenv(key, value)
		}
	}

	return scanner.Err()
}
