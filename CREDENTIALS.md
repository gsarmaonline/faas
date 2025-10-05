# Environment-Based Credential Management

This document describes the improved credential management system that reads sensitive data from environment variables instead of requiring them in payloads.

## Benefits

- **🔒 Enhanced Security**: API keys and tokens are not exposed in payloads
- **🌍 Environment Flexibility**: Easy to use different credentials for dev/staging/prod
- **📝 Better Error Messages**: Clear guidance on where to set credentials
- **🔄 Backward Compatibility**: Payload values still work but env vars take precedence

## How It Works

The credential manager follows this priority order:

1. **Payload value** (if provided and not empty)
2. **Environment variable** (fallback)
3. **Error** (if neither is available for required credentials)

## Example: Email Function

### Before (Less Secure)

```json
{
  "api_key": "SG.very-secret-key-exposed-in-payload",
  "from_email": "sender@example.com",
  "to_email": "recipient@example.com",
  "subject": "Test",
  "plain_text": "Hello world"
}
```

### After (More Secure)

```bash
# Set environment variable once
export SENDGRID_API_KEY="SG.your-secret-key-here"
```

```json
{
  "from_email": "sender@example.com",
  "to_email": "recipient@example.com",
  "subject": "Test",
  "plain_text": "Hello world"
}
```

## Supported Environment Variables

| Function | Environment Variable                                   | Description                  |
| -------- | ------------------------------------------------------ | ---------------------------- |
| Email    | `SENDGRID_API_KEY`                                     | SendGrid API key             |
| Slack    | `SLACK_API_TOKEN`                                      | Slack bot token              |
| SMS      | `TWILIO_ACCOUNT_SID`, `TWILIO_AUTH_TOKEN`              | Twilio credentials           |
| Docker   | `DOCKER_REGISTRY_USERNAME`, `DOCKER_REGISTRY_PASSWORD` | Registry auth                |
| GitHub   | `GITHUB_TOKEN`                                         | GitHub personal access token |

## Development Setup

Create a `.env` file in your project root:

```bash
# .env file (never commit this!)
SENDGRID_API_KEY=SG.your-development-key
SLACK_API_TOKEN=xoxb-your-slack-token
TWILIO_ACCOUNT_SID=your-twilio-sid
TWILIO_AUTH_TOKEN=your-twilio-token
GITHUB_TOKEN=ghp_your-github-token
```

## Production Deployment

Set environment variables in your deployment environment:

```bash
# Kubernetes
kubectl create secret generic faas-credentials \
  --from-literal=SENDGRID_API_KEY=SG.prod-key \
  --from-literal=SLACK_API_TOKEN=xoxb-prod-token

# Docker
docker run -e SENDGRID_API_KEY=SG.prod-key myapp

# Heroku
heroku config:set SENDGRID_API_KEY=SG.prod-key

# AWS Lambda
# Set in Lambda environment variables through console or terraform
```

## Error Handling

When credentials are missing, you get helpful error messages:

```
missing required credential: api_key (provide in payload or set SENDGRID_API_KEY environment variable)
```

## Implementation for New Functions

When creating new functions, use the credential manager:

```go
func (myFunc *MyFunction) ParsePayload(payload intf.Payload) error {
    credManager := helpers.NewCredentialManager()

    processedInput := MyInput{
        APIKey: credManager.GetCredential(payload["api_key"], helpers.EnvMyServiceAPIKey),
        // ... other fields
    }

    myFunc.Input = processedInput
    return nil
}
```
