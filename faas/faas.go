package faas

import (
	"context"
	"fmt"

	"github.com/gsarmaonline/faas/faas/functions"
	"github.com/gsarmaonline/faas/faas/intf"
)

type (
	Faas struct {
		ctx context.Context

		functions map[string]intf.Function
	}
)

func NewFaas(ctx context.Context) (*Faas, error) {
	faas := &Faas{
		ctx:       ctx,
		functions: make(map[string]intf.Function),
	}
	if err := faas.RegisterFunctions([]intf.Function{
		functions.NewSlack(),
		functions.NewEmailAction(),
		functions.NewSmsAction(),
		functions.NewDockerRegistryAction(),
		functions.NewHttpAction(),
		functions.NewLoggerAction(),
		functions.NewGithubAction(),
	}); err != nil {
		return nil, err
	}
	return faas, nil
}

func (faas *Faas) RegisterFunctions(functions []intf.Function) (err error) {
	for _, function := range functions {
		if _, exists := faas.functions[function.GetConfig().Name]; exists {
			err = fmt.Errorf("function with name %s already exists", function.GetConfig().Name)
			return
		}
		faas.functions[function.GetConfig().Name] = function
	}
	return
}

func (faas *Faas) ExecuteFunction(name string) (output intf.FunctionOutput, err error) {
	function, exists := faas.functions[name]
	if !exists {
		err = fmt.Errorf("function with name %s does not exist", name)
		return
	}
	if err = function.Validate(); err != nil {
		return
	}
	if output, err = function.Execute(); err != nil {
		return
	}
	return
}
