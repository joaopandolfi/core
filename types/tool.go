package types

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
)

// Tool represents a function that an AI agent can use
type Tool struct {
	// The name of the tool
	Name string

	// The description of the tool
	Description string

	// WrappedToolFunction is the in-code function to be called by the AI agent
	// when using this tool "wrapped" by WrapFunction.
	//
	// A tool's function expects 2 arguments: a context and a byte slice.
	// The byte slice should be the raw arguments provided by the LLM. The wrapped
	// function will then automatically unmarshal those arguments to the underlying
	// function.
	//
	// WrappedToolFunction should have 2 returns: an interface and an error. The interface
	// may be anything defined by the wrapped function (a struct, a string, a number, etc.).
	WrappedToolFunction func(ctx context.Context, args []byte) (interface{}, error)

	// JSONSchema is the raw JSON schema data as a byte slice that will be provided
	// to a tool calling LLM for argument validation.
	JSONSchema []byte
}

// WrapToolFunction dynamically, at runtime, converts the input function to a "WrappedToolFunction"
// that can be used as part of Tool.WrappedToolFunction - i.e., a function of type:
// func(context.Context []byte) (interface{}, error)
func WrapToolFunction(fn interface{}) (func(context.Context, []byte) (interface{}, error), error) {
	fnValue := reflect.ValueOf(fn)

	if fnValue.Kind() != reflect.Func {
		panic("fn must be a function")
	}

	fnType := fnValue.Type()

	// Validate function signature
	if fnType.NumIn() != 2 || fnType.NumOut() != 2 {
		panic("function must have two parameters and two return values")
	}
	if fnType.In(0) != reflect.TypeOf((*context.Context)(nil)).Elem() {
		panic("first parameter must be context.Context")
	}
	argType := fnType.In(1)
	if argType.Kind() != reflect.Ptr || argType.Elem().Kind() != reflect.Struct {
		panic("second parameter must be a pointer to a struct")
	}

	return func(ctx context.Context, args []byte) (interface{}, error) {
		fmt.Printf("In wrapped func\nargs: %s\n", args)
		// Create a new instance of the target struct
		target := reflect.New(argType.Elem()).Interface()
		if err := json.Unmarshal(args, target); err != nil {
			return nil, fmt.Errorf("error unmarshaling args: %v", err)
		}

		// Call the original function
		results := fnValue.Call([]reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(target),
		})

		// Extract return values
		var result interface{}
		if !results[0].IsNil() {
			result = results[0].Interface()
		}

		var errResult error
		if !results[1].IsNil() {
			errResult = results[1].Interface().(error)
		}

		return result, errResult
	}, nil
}
