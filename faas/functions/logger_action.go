package functions

import (
	"log"

	"github.com/gsarmaonline/faas/faas/intf"
)

type (
	LoggerInput struct {
		Message string `json:"message"`
	}
	LoggerAction struct {
		Input LoggerInput
	}
)

func NewLoggerAction() (loggerAction *LoggerAction) {
	return &LoggerAction{}
}

func (loggerAction LoggerAction) GetConfig() intf.FunctionConfig {
	return intf.FunctionConfig{Name: "logger"}
}

func (loggerAction *LoggerAction) ParsePayload(payload intf.Payload) error {
	processedInput := LoggerInput{}

	// Optional message field
	if message, exists := payload["message"]; exists && message != nil {
		processedInput.Message = message.(string)
	}

	loggerAction.Input = processedInput
	return nil
}

func (loggerAction LoggerAction) Validate() (err error) {
	// Logger doesn't require any specific validation
	return nil
}

func (loggerAction LoggerAction) Execute() (output intf.FunctionOutput, err error) {
	if loggerAction.Input.Message != "" {
		log.Println("From Logger action:", loggerAction.Input.Message)
	} else {
		log.Println("From Logger action: (no message provided)")
	}
	return
}
