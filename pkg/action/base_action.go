package action

import (
	"context"
	//validator "github.com/rezakhademix/govalidator/v2"
)

// Define the base type with the shared template method.
type BaseAction struct {
	ctx context.Context
}

func (ba *BaseAction) SetContext(ctx context.Context) {
	ba.ctx = ctx
}

func (ba *BaseAction) Context() context.Context {
	return ba.ctx
}

func (ba *BaseAction) Performer() any {
	return nil
}

func (ba *BaseAction) IsAllowed() (bool, error) {
	// See https://github.com/rezakhademix/govalidator
	// For example:
	// v := validator.New()
	// return v.IsPassed(), v.Errors()
	return true, nil
}

func (ba *BaseAction) IsEnabled() (bool, error) {
	// See https://github.com/rezakhademix/govalidator
	// For example:
	// v := validator.New()
	// return v.IsPassed(), v.Errors()
	return true, nil
}

func (ba *BaseAction) IsValid(ctx context.Context) (bool, error) {
	// See https://github.com/rezakhademix/govalidator
	// For example:
	// v := validator.New()
	// return v.IsPassed(), v.Errors()
	return true, nil
}
