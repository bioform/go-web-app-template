package main

import (
	"context"
	"fmt"
)

// Define the base type with the shared template method.
type ActionPerformer struct{}

// Template method that orchestrates the behavior.
func (ap *ActionPerformer) Call(action Performable) {
	// Shared pre-execution logic
	fmt.Println("Shared logic before perform")

	// Delegate the specific action to the embedding type
	action.Perform()

	// Shared post-execution logic
	fmt.Println("Shared logic after perform")
}

type Performable interface {
	Perform()
}

// Define an interface that the embedding types must implement.
type Action interface {
	Perform()
	Call(action Performable)
}

func CallAction(action Action) {
	action.Call(action)
}

// Define a specific type embedding ActionPerformer.
type MyAction struct {
	ActionPerformer
	ctx context.Context
}

// Implement the specific behavior for MyAction.
func (action *MyAction) Perform() {
	fmt.Println("MyAction-specific perform logic")
}

func main() {
	action := &MyAction{
		ctx: context.TODO(),
	}

	// Call the template method, which handles the shared logic implicitly.
	CallAction(action)
}
