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
	return &AuthorizationError{
		ActionError: ActionError{performer: performer},
	}
}

func NewDisabledError(performer any, errs ErrorMap) *DisabledError {
	return &DisabledError{
		ActionError: ActionError{performer: performer},
		ErrorMap:    errs,
	}
}

func NewValidationError(performer any, errs ErrorMap) *ValidationError {
	return &ValidationError{
		ActionError: ActionError{performer: performer},
		ErrorMap:    errs,
	}
}

func (*ActionError) Error() string {
	return "authorization failed"
}

func (*DisabledError) Error() string {
	return "action is not enabled"
}

func (*ValidationError) Error() string {
	return "validation failed"
}
