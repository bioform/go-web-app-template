package action

type ErrorMap map[string]string

type ActionError struct {
	performer any
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

func NewAuthorizationError(performer any) *AuthorizationError {
	err := &AuthorizationError{
		ActionError: ActionError{performer: performer},
	}

	return err
}

func NewDisabledError(performer any, errs ErrorMap) *DisabledError {
	err := &DisabledError{
		ActionError: ActionError{performer: performer},
		ErrorMap:    errs,
	}

	return err
}

func NewValidationError(performer any, errs ErrorMap) *ValidationError {
	err := &ValidationError{
		ActionError: ActionError{performer: performer},
		ErrorMap:    errs,
	}

	return err
}

func (*ActionError) Error() string {
	return "authorization failed"
}

func (e *ActionError) Performer() any {
	return e.performer
}

func (*DisabledError) Error() string {
	return "action is not enabled"
}

func (*ValidationError) Error() string {
	return "validation failed"
}
