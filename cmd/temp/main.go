package main

import (
	"errors"
	"fmt"
	"sync"
)

// ErrorHandler is a generic type for error handler functions with a specific error type.
type ErrorHandler[T error] func(err T)

// Internal storage for error handlers
var (
	errorHandlers []func(error)
	mu            sync.Mutex
)

// RegisterErrorHandler registers a new error handler for a specific error type.
func RegisterErrorHandler[T error](handler ErrorHandler[T]) {
	mu.Lock()
	defer mu.Unlock()
	// Wrap the handler in a type-agnostic function and store it
	errorHandlers = append(errorHandlers, func(err error) {
		var specificErr T
		if errors.As(err, &specificErr) {
			handler(specificErr)
		}
	})
}

// RunErrorHandlers calls all registered error handlers for the provided error.
func RunErrorHandlers(err error) {
	mu.Lock()
	defer mu.Unlock()
	for _, handler := range errorHandlers {
		handler(err)
	}
}

// Custom error types implementing the error interface
type MyError struct {
	msg string
}

func (e MyError) Error() string {
	return e.msg
}

type AnotherError struct {
	code int
	msg  string
}

func (e AnotherError) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s", e.code, e.msg)
}

func main() {
	// Register handlers for specific error types
	RegisterErrorHandler(func(err MyError) {
		fmt.Println("Handled MyError:", err.msg)
	})
	RegisterErrorHandler(func(err AnotherError) {
		fmt.Printf("Handled AnotherError: %s\n", err.Error())
	})
	RegisterErrorHandler(func(err error) {
		fmt.Println("Generic handler:", err.Error())
	})

	// Trigger handlers with different error types
	myErr := MyError{msg: "This is MyError"}
	anotherErr := AnotherError{code: 404, msg: "Not Found"}
	genericErr := errors.New("a generic error")

	fmt.Println("Running handlers for MyError:")
	RunErrorHandlers(myErr)

	fmt.Println("\nRunning handlers for AnotherError:")
	RunErrorHandlers(anotherErr)

	fmt.Println("\nRunning handlers for a generic error:")
	RunErrorHandlers(genericErr)
}
