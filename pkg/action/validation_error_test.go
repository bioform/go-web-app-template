package action

import (
	"reflect"
	"testing"
)

func TestNewValidationError(t *testing.T) {
	errs := ErrorMap{
		"field1": "error message 1",
		"field2": "error message 2",
	}

	ve := NewValidationError(nil, errs)
	if ve == nil {
		t.Fatalf("Expected NewValidationError to return a non-nil value")
	}

	if !reflect.DeepEqual(ve.ErrorMap, errs) {
		t.Errorf("Expected Errors() to return %v, got %v", errs, ve.ErrorMap)
	}
}

func TestValidationError_Error(t *testing.T) {
	errs := ErrorMap{
		"field": "error message",
	}

	ve := NewValidationError(nil, errs)
	expectedErrMsg := "validation failed"
	if ve.Error() != expectedErrMsg {
		t.Errorf("Expected Error() to return %q, got %q", expectedErrMsg, ve.Error())
	}
}
