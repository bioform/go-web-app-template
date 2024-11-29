package action_test

import (
	"github.com/bioform/go-web-app-template/pkg/action"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Action Errors", func() {
	var performer any

	BeforeEach(func() {
		performer = "test performer"
	})

	Describe("AuthorizationError", func() {
		It("should return the correct error message and performer", func() {
			err := action.NewAuthorizationError(performer)
			Expect(err.Error()).To(Equal("authorization failed"))
			Expect(err.Performer()).To(Equal(performer))
		})
	})

	Describe("DisabledError", func() {
		It("should return the correct error message, performer, and error map", func() {
			errs := action.ErrorMap{"feature": "disabled"}
			err := action.NewDisabledError(performer, errs)
			Expect(err.Error()).To(Equal("action is not enabled"))
			Expect(err.Performer()).To(Equal(performer))
			Expect(err.ErrorMap).To(Equal(errs))
		})
	})

	Describe("ValidationError", func() {
		It("should return the correct error message, performer, and error map", func() {
			errs := action.ErrorMap{"field": "invalid"}
			err := action.NewValidationError(performer, errs)
			Expect(err.Error()).To(Equal("validation failed"))
			Expect(err.Performer()).To(Equal(performer))
			Expect(err.ErrorMap).To(Equal(errs))
		})
	})
})
