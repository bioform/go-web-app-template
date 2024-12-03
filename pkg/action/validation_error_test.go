package action_test

import (
	"github.com/bioform/go-web-app-template/pkg/action"
	"github.com/bioform/go-web-app-template/pkg/action/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Action Errors", func() {
	var mockAction *mocks.Action
	var performer any

	BeforeEach(func() {
		performer = "test performer"
		mockAction = &mocks.Action{}
		mockAction.On("Performer").Return(performer)
	})

	Describe("AuthorizationError", func() {
		It("should return the correct error message and performer", func() {
			err := action.NewAuthorizationError(mockAction)
			Expect(err.Error()).To(Equal("performer: test performer"))
		})
	})

	Describe("DisabledError", func() {
		It("should return the correct error message, performer, and error map", func() {
			errs := action.ErrorMap{"feature": "disabled"}
			err := action.NewDisabledError(mockAction, errs)
			Expect(err.Error()).To(Equal("performer: test performer, action is not enabled: map[feature:disabled]"))
			Expect(err.ErrorMap).To(Equal(errs))
		})
	})

	Describe("ValidationError", func() {
		It("should return the correct error message, performer, and error map", func() {
			errs := action.ErrorMap{"field": "invalid"}
			err := action.NewValidationError(mockAction, errs)
			Expect(err.Error()).To(Equal("performer: test performer, validation failed: map[field:invalid]"))
			Expect(err.ErrorMap).To(Equal(errs))
		})
	})
})
