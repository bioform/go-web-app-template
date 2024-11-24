package action

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewValidationError(t *testing.T) {
	errs := ErrorMap{
		"field1": "error message 1",
		"field2": "error message 2",
	}

	ve := NewValidationError(errs)
	if ve == nil {
		t.Fatalf("Expected NewValidationError to return a non-nil value")
	}

	if !reflect.DeepEqual(ve.Errors(), errs) {
		t.Errorf("Expected Errors() to return %v, got %v", errs, ve.Errors())
	}
}

func TestValidationError_Error(t *testing.T) {
	errs := ErrorMap{
		"field": "error message",
	}

	ve := NewValidationError(errs)
	expectedErrMsg := fmt.Sprintf("validation failed: %v", errs)
	if ve.Error() != expectedErrMsg {
		t.Errorf("Expected Error() to return %q, got %q", expectedErrMsg, ve.Error())
	}
}
