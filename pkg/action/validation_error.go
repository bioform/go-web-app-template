package action

import "fmt"

type ErrorMap map[string]string

type ActionError struct {
	errs ErrorMap
}

type DisabledError struct {
	ActionError
}

type ValidationError struct {
	ActionError
}

func NewDisabledError(errs ErrorMap) *DisabledError {
	return &DisabledError{ActionError: ActionError{errs: errs}}
}

func NewValidationError(errs ErrorMap) *ValidationError {
	return &ValidationError{ActionError: ActionError{errs: errs}}
}

func (v *ActionError) Error() string {
	return fmt.Sprintf("validation failed: %v", v.errs)
}

func (v *ActionError) Errors() ErrorMap {
	return v.errs
}
