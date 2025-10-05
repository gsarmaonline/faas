package faas

import (
	"context"
	"fmt"

	"github.com/gsarmaonline/faas/faas/intf"
)

type (
	Faas struct {
		ctx context.Context

		functions map[string]intf.Function
	}
)

func NewFaas(ctx context.Context) *Faas {
	return &Faas{
		ctx:       ctx,
		functions: make(map[string]intf.Function),
	}
}

func (faas *Faas) RegisterFunction(function intf.Function) (err error) {
	if _, exists := faas.functions[function.GetConfig().Name]; exists {
		err = fmt.Errorf("function with name %s already exists", function.GetConfig().Name)
		return
	}
	faas.functions[function.GetConfig().Name] = function
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
