package intf

type (
	Payload map[string]interface{}

	FunctionConfig struct {
		Name string `json:"name"`
	}

	Function interface {
		GetConfig() FunctionConfig
		ParsePayload(Payload) error
		Validate() error
		Execute() (FunctionOutput, error)
	}
	FunctionOutput interface {
		GetPayload() (Payload, error)
	}
)
