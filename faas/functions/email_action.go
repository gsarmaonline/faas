package functions

import (
	"fmt"
	"log"

	"github.com/gsarmaonline/faas/faas/intf"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type (
	EmailInput struct {
		ApiKey    string `json:"api_key"`
		FromEmail string `json:"from_email"`
		FromName  string `json:"from_name"`
		ToEmail   string `json:"to_email"`
		ToName    string `json:"to_name"`
		Subject   string `json:"subject"`
		PlainText string `json:"plain_text"`
		HtmlText  string `json:"html_text"`
	}
	EmailAction struct {
		Input EmailInput
	}
)

func NewEmailAction() (emailAction *EmailAction) {
	return &EmailAction{}
}

func (emailAction EmailAction) GetConfig() intf.FunctionConfig {
	return intf.FunctionConfig{Name: "email"}
}

func (emailAction *EmailAction) ParsePayload(payload intf.Payload) error {
	processedInput := EmailInput{
		ApiKey:    payload["api_key"].(string),
		FromEmail: payload["from_email"].(string),
		ToEmail:   payload["to_email"].(string),
		Subject:   payload["subject"].(string),
	}

	// Optional fields
	if fromName, exists := payload["from_name"]; exists && fromName != nil {
		processedInput.FromName = fromName.(string)
	}

	if toName, exists := payload["to_name"]; exists && toName != nil {
		processedInput.ToName = toName.(string)
	}

	if plainText, exists := payload["plain_text"]; exists && plainText != nil {
		processedInput.PlainText = plainText.(string)
	}

	if htmlText, exists := payload["html_text"]; exists && htmlText != nil {
		processedInput.HtmlText = htmlText.(string)
	}

	emailAction.Input = processedInput
	return nil
}

func (emailAction EmailAction) Validate() (err error) {
	if emailAction.Input.ApiKey == "" {
		return fmt.Errorf("missing required field: api_key")
	}
	if emailAction.Input.FromEmail == "" {
		return fmt.Errorf("missing required field: from_email")
	}
	if emailAction.Input.ToEmail == "" {
		return fmt.Errorf("missing required field: to_email")
	}
	if emailAction.Input.Subject == "" {
		return fmt.Errorf("missing required field: subject")
	}
	if emailAction.Input.PlainText == "" && emailAction.Input.HtmlText == "" {
		return fmt.Errorf("at least one of plain_text or html_text must be provided")
	}
	return nil
}

func (emailAction EmailAction) Execute() (output intf.FunctionOutput, err error) {
	// Create sender and recipient
	from := mail.NewEmail(emailAction.Input.FromName, emailAction.Input.FromEmail)
	to := mail.NewEmail(emailAction.Input.ToName, emailAction.Input.ToEmail)

	// Create the email message
	var message *mail.SGMailV3

	if emailAction.Input.HtmlText != "" {
		// If HTML content is provided, create email with both plain text and HTML
		message = mail.NewSingleEmail(from, emailAction.Input.Subject, to, emailAction.Input.PlainText, emailAction.Input.HtmlText)
	} else {
		// If only plain text is provided
		message = mail.NewSingleEmail(from, emailAction.Input.Subject, to, emailAction.Input.PlainText, "")
	}

	// Create SendGrid client
	client := sendgrid.NewSendClient(emailAction.Input.ApiKey)

	// Send the email
	response, err := client.Send(message)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return nil, err
	}

	// Check if the response indicates success
	if response.StatusCode >= 400 {
		err = fmt.Errorf("sendgrid API error: status code %d, body: %s", response.StatusCode, response.Body)
		log.Printf("SendGrid API error: %v", err)
		return nil, err
	}

	log.Printf("Email sent successfully. Status: %d", response.StatusCode)
	return nil, nil
}
