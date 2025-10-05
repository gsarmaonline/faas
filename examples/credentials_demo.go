package main

import (
	"context"
	"fmt"

	"github.com/gsarmaonline/faas/faas"
	"github.com/gsarmaonline/faas/faas/helpers"
	"github.com/gsarmaonline/faas/faas/intf"
)

func main() {
	fmt.Println("🔐 FAAS Credential Management Demo")
	fmt.Println("==================================")

	// Demonstrate environment variable validation
	fmt.Println("\n1. Checking environment variables...")
	requiredVars := []string{
		helpers.EnvSendGridAPIKey,
		helpers.EnvSlackAPIToken,
	}

	missing := helpers.ValidateEnvironmentVars(requiredVars)
	if len(missing) > 0 {
		fmt.Printf("⚠️  Missing environment variables: %v\n", missing)
		fmt.Println("Set them for secure credential management:")
		for _, envVar := range missing {
			fmt.Printf("   export %s=your-secret-key\n", envVar)
		}
	} else {
		fmt.Println("✅ All required environment variables are set!")
	}

	// Demonstrate credential manager usage
	fmt.Println("\n2. Testing credential manager...")
	credManager := helpers.NewCredentialManager()

	// Test with environment variable
	if apiKey := credManager.GetCredential(nil, helpers.EnvSendGridAPIKey); apiKey != "" {
		fmt.Printf("✅ SendGrid API key loaded from environment (first 10 chars): %s...\n", apiKey[:min(10, len(apiKey))])
	} else {
		fmt.Println("❌ No SendGrid API key found in environment")
	}

	// Test payload override
	payloadAPIKey := "payload-override-key"
	overrideKey := credManager.GetCredential(payloadAPIKey, helpers.EnvSendGridAPIKey)
	fmt.Printf("✅ Payload override works: %s\n", overrideKey)

	// Demonstrate FAAS framework initialization
	fmt.Println("\n3. Initializing FAAS framework...")
	faasFramework, err := faas.NewFaas(context.Background())
	if err != nil {
		fmt.Printf("❌ Error initializing FAAS: %v\n", err)
		return
	}

	fmt.Println("✅ FAAS framework initialized with secure credential management!")
	fmt.Println("\n📋 Available functions with environment-based credentials:")
	fmt.Println("   • email    (uses SENDGRID_API_KEY)")
	fmt.Println("   • slack    (uses SLACK_API_TOKEN)")
	fmt.Println("   • sms      (uses TWILIO_ACCOUNT_SID, TWILIO_AUTH_TOKEN)")
	fmt.Println("   • docker   (uses DOCKER_REGISTRY_USERNAME, DOCKER_REGISTRY_PASSWORD)")
	fmt.Println("   • github   (uses GITHUB_TOKEN)")
	fmt.Println("   • http     (no credentials required)")
	fmt.Println("   • logger   (no credentials required)")

	// Demonstrate secure payload (no credentials exposed)
	fmt.Println("\n4. Example secure payload (no API keys exposed):")
	securePayload := intf.Payload{
		"from_email": "noreply@example.com",
		"to_email":   "user@example.com",
		"subject":    "Secure FAAS Email",
		"plain_text": "This email was sent securely without exposing API keys in the payload!",
	}

	fmt.Printf("   %+v\n", securePayload)
	fmt.Println("\n✨ API keys are safely managed through environment variables!")

	// Show what would happen without environment variables
	if len(missing) > 0 {
		fmt.Println("\n⚠️  To actually send emails, set the SENDGRID_API_KEY environment variable")
		fmt.Println("   and then the email function will work without requiring api_key in payload")
	}

	_ = faasFramework // Use the variable to avoid linting error
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
