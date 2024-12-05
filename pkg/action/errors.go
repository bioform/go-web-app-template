package action

import "fmt"

type ErrorMap map[string]string

type ActionError struct {
	action Action
	error
}

type CallbackError struct {
	ActionError
}

type AuthorizationError struct {
	ActionError
}

type DisabledError struct {
	ActionError
}

type ValidationError struct {
	ActionError
}

func (e ErrorMap) Error() string {
	return fmt.Sprintf("%v", map[string]string(e))
}

func NewActionError(action Action, err error) *ActionError {
	return &ActionError{action: action, error: err}
}

func NewCallbackError(action Action, err error) *CallbackError {
	return &CallbackError{ActionError{action: action, error: err}}
}

func NewAuthorizationError(action Action) *AuthorizationError {
	return &AuthorizationError{ActionError: ActionError{action: action}}
}

func NewDisabledError(action Action, errs ErrorMap) *DisabledError {
	return &DisabledError{
		ActionError: ActionError{action: action, error: errs},
	}
}

func NewValidationError(action Action, errs ErrorMap) *ValidationError {
	return &ValidationError{
		ActionError: ActionError{action: action, error: errs},
	}
}

func (e ActionError) Unwrap() error {
	return e.error
}

func (e CallbackError) Unwrap() error {
	return e.ActionError
}

func (e AuthorizationError) Unwrap() error {
	return e.ActionError
}

func (e DisabledError) Unwrap() error {
	return e.ActionError
}

func (e ValidationError) Unwrap() error {
	return e.ActionError
}

func (e ActionError) Error() string {
	return fmt.Sprintf("action: %T, performer: %v, error: %v", e.action, e.action.Performer(), e.error)
}

func (e CallbackError) Error() string {
	return fmt.Sprintf("callback: %v", e.ActionError)
}

func (e AuthorizationError) Error() string {
	return fmt.Sprintf("authorization: %v", e.ActionError)
}

func (e DisabledError) Error() string {
	return fmt.Sprintf("not enabled: %v", e.ActionError)
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation failed: %v", e.ActionError)
}

func (e DisabledError) Errors() ErrorMap {
	return e.error.(ErrorMap)
}

func (e ValidationError) Errors() ErrorMap {
	return e.error.(ErrorMap)
}
