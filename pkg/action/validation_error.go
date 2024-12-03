package action

import "fmt"

type ErrorMap map[string]string

type ActionError struct {
	action Action
}

type AuthorizationError struct {
	ActionError
}

type DisabledError struct {
	ActionError
	ErrorMap
}

type ValidationError struct {
	ActionError
	ErrorMap
}

func (e ErrorMap) Error() string {
	return fmt.Sprintf("%v", map[string]string(e))
}

func NewAuthorizationError(action Action) *AuthorizationError {
	err := &AuthorizationError{
		ActionError: ActionError{action: action},
	}

	return err
}

func NewDisabledError(action Action, errs ErrorMap) *DisabledError {
	err := &DisabledError{
		ActionError: ActionError{action: action},
		ErrorMap:    errs,
	}

	return err
}

func NewValidationError(action Action, errs ErrorMap) *ValidationError {
	err := &ValidationError{
		ActionError: ActionError{action: action},
		ErrorMap:    errs,
	}

	return err
}

func (e ActionError) Error() string {
	return fmt.Sprintf("performer: %v", e.action.Performer())
}

func (e DisabledError) Error() string {
	return fmt.Sprintf("%s, action is not enabled: %v", e.ActionError, e.ErrorMap)
}

func (e ValidationError) Unwrap() []error {
	return []error{e.ActionError, e.ErrorMap}
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s, validation failed: %v", e.ActionError, e.ErrorMap)
}
