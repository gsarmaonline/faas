package faas

import (
	"context"
	"fmt"
	"testing"

	"github.com/gsarmaonline/faas/faas/functions"
	"github.com/gsarmaonline/faas/faas/intf"
)

// Mock function for testing
type MockFunction struct {
	name      string
	validated bool
	executed  bool
	shouldErr bool
}

func (m *MockFunction) GetConfig() intf.FunctionConfig {
	return intf.FunctionConfig{Name: m.name}
}

func (m *MockFunction) ParsePayload(payload intf.Payload) error {
	return nil
}

func (m *MockFunction) Validate() error {
	m.validated = true
	if m.shouldErr {
		return fmt.Errorf("validation error")
	}
	return nil
}

func (m *MockFunction) Execute() (intf.FunctionOutput, error) {
	m.executed = true
	if m.shouldErr {
		return nil, fmt.Errorf("execution error")
	}
	return nil, nil
}

func TestNewFaas(t *testing.T) {
	ctx := context.Background()
	faas, err := NewFaas(ctx)

	if err != nil {
		t.Errorf("NewFaas() error = %v, want nil", err)
		return
	}

	if faas == nil {
		t.Error("NewFaas() returned nil")
		return
	}

	if faas.ctx != ctx {
		t.Error("NewFaas() context not set correctly")
	}

	if faas.functions == nil {
		t.Error("NewFaas() functions map not initialized")
	}

	// Check that default functions are registered
	expectedDefaults := []string{"slack", "email", "docker_registry", "http", "logger", "github"}
	for _, name := range expectedDefaults {
		if _, exists := faas.functions[name]; !exists {
			t.Errorf("NewFaas() should register %s function by default", name)
		}
	}
}

func TestFaas_RegisterFunctions(t *testing.T) {
	ctx := context.Background()
	faas := &Faas{
		ctx:       ctx,
		functions: make(map[string]intf.Function),
	}

	tests := []struct {
		name      string
		functions []intf.Function
		wantErr   bool
		errMsg    string
	}{
		{
			name: "register single function",
			functions: []intf.Function{
				&MockFunction{name: "test1"},
			},
			wantErr: false,
		},
		{
			name: "register multiple functions",
			functions: []intf.Function{
				&MockFunction{name: "test2"},
				&MockFunction{name: "test3"},
			},
			wantErr: false,
		},
		{
			name: "register duplicate function",
			functions: []intf.Function{
				&MockFunction{name: "test1"}, // Already registered in first test
			},
			wantErr: true,
			errMsg:  "function with name test1 already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := faas.RegisterFunctions(tt.functions)

			if tt.wantErr {
				if err == nil {
					t.Errorf("RegisterFunctions() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("RegisterFunctions() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("RegisterFunctions() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				// Verify functions are registered
				for _, fn := range tt.functions {
					if _, exists := faas.functions[fn.GetConfig().Name]; !exists {
						t.Errorf("Function %s not registered", fn.GetConfig().Name)
					}
				}
			}
		})
	}
}

func TestFaas_ExecuteFunction(t *testing.T) {
	ctx := context.Background()
	faas := &Faas{
		ctx:       ctx,
		functions: make(map[string]intf.Function),
	}

	// Register test functions
	mockFunc1 := &MockFunction{name: "success_func", shouldErr: false}
	mockFunc2 := &MockFunction{name: "validation_error_func", shouldErr: true}
	mockFunc3 := &MockFunction{name: "execution_error_func", shouldErr: false}

	// Set execution_error_func to fail during execution
	mockFunc3.validated = true // This will bypass validation error
	faas.functions["success_func"] = mockFunc1
	faas.functions["validation_error_func"] = mockFunc2
	faas.functions["execution_error_func"] = mockFunc3

	tests := []struct {
		name         string
		functionName string
		wantErr      bool
		errContains  string
	}{
		{
			name:         "execute existing function successfully",
			functionName: "success_func",
			wantErr:      false,
		},
		{
			name:         "execute non-existing function",
			functionName: "non_existing",
			wantErr:      true,
			errContains:  "function with name non_existing does not exist",
		},
		{
			name:         "execute function with validation error",
			functionName: "validation_error_func",
			wantErr:      true,
			errContains:  "validation error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := faas.ExecuteFunction(tt.functionName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ExecuteFunction() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errContains != "" && err.Error() != tt.errContains {
					t.Errorf("ExecuteFunction() error = %v, want error containing %v", err.Error(), tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("ExecuteFunction() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				// Verify that validation and execution were called for successful case
				if tt.functionName == "success_func" {
					if !mockFunc1.validated {
						t.Error("Function validation was not called")
					}
					if !mockFunc1.executed {
						t.Error("Function execution was not called")
					}
				}
			}
		})
	}
}

func TestFaas_Integration_WithRealFunctions(t *testing.T) {
	ctx := context.Background()
	faas := &Faas{
		ctx:       ctx,
		functions: make(map[string]intf.Function),
	}

	// Register all the real functions
	realFunctions := []intf.Function{
		functions.NewSlack(),
		functions.NewDockerRegistryAction(),
		functions.NewHttpAction(),
		functions.NewLoggerAction(),
		functions.NewGithubAction(),
		functions.NewEmailAction(),
	}

	err := faas.RegisterFunctions(realFunctions)
	if err != nil {
		t.Errorf("Failed to register real functions: %v", err)
		return
	}

	// Test that all functions are registered
	expectedFunctions := []string{"slack", "docker_registry", "http", "logger", "github", "email"}
	for _, name := range expectedFunctions {
		if _, exists := faas.functions[name]; !exists {
			t.Errorf("Function %s not registered", name)
		}
	}

	// Test that we can't register duplicate functions
	err = faas.RegisterFunctions([]intf.Function{functions.NewSlack()})
	if err == nil {
		t.Error("Should not be able to register duplicate slack function")
	}
}

func TestFaas_ExecuteFunction_WithPayload(t *testing.T) {
	ctx := context.Background()
	faas := &Faas{
		ctx:       ctx,
		functions: make(map[string]intf.Function),
	}

	// Register logger function for testing (it has minimal validation requirements)
	loggerFunc := functions.NewLoggerAction()
	err := faas.RegisterFunctions([]intf.Function{loggerFunc})
	if err != nil {
		t.Errorf("Failed to register logger function: %v", err)
		return
	}

	// Test executing logger function
	_, err = faas.ExecuteFunction("logger")
	if err != nil {
		t.Errorf("ExecuteFunction('logger') error = %v, want nil", err)
	}
}

func TestFaas_GetRegisteredFunctions(t *testing.T) {
	ctx := context.Background()
	faas := &Faas{
		ctx:       ctx,
		functions: make(map[string]intf.Function),
	}

	// Register some test functions
	testFunctions := []intf.Function{
		&MockFunction{name: "func1"},
		&MockFunction{name: "func2"},
		&MockFunction{name: "func3"},
	}

	err := faas.RegisterFunctions(testFunctions)
	if err != nil {
		t.Errorf("Failed to register test functions: %v", err)
		return
	}

	// Verify the correct number of functions are registered
	if len(faas.functions) != 3 {
		t.Errorf("Expected 3 registered functions, got %d", len(faas.functions))
	}

	// Verify specific functions exist
	expectedNames := []string{"func1", "func2", "func3"}
	for _, name := range expectedNames {
		if _, exists := faas.functions[name]; !exists {
			t.Errorf("Function %s should be registered", name)
		}
	}
}
