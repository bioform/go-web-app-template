package action_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/bioform/go-web-app-template/pkg/action"
	"github.com/bioform/go-web-app-template/pkg/action/mocks"
)

var _ = Describe("ActionPerformer", func() {
	var (
		ctx                     context.Context
		performer               *action.ActionPerformer[*mocks.Action]
		mockAction              *mocks.Action
		mockTransactionProvider *mocks.TransactionProvider
	)

	BeforeEach(func() {
		ctx = context.Background()
		mockAction = mocks.NewAction(GinkgoT())
		mockTransactionProvider = mocks.NewTransactionProvider(GinkgoT())

		// Mock SetContext and TransactionProvider
		mockAction.EXPECT().SetContext(ctx).Maybe()
		mockAction.EXPECT().TransactionProvider().Return(mockTransactionProvider).Maybe()

		// Mock the transaction
		call := mockTransactionProvider.EXPECT().Transaction(mock.Anything, mock.Anything).Maybe()
		call.Run(func(args mock.Arguments) {
			ctx := args.Get(0).(context.Context)
			lambda := args.Get(1).(func(context.Context) error)
			err := lambda(ctx)
			call.Return(err)
		})

		// Update the New function call to include the context
		performer = action.New(ctx, mockAction)
	})

	Describe("Action", func() {
		It("should return the action", func() {
			Expect(performer.Action()).To(Equal(mockAction))
		})
	})

	Describe("Perform", func() {
		It("should perform action successfully", func() {
			mockAction.EXPECT().Context().Return(ctx)

			mockAction.EXPECT().IsAllowed().Return(true, nil)
			mockAction.EXPECT().IsEnabled().Return(true, nil)
			mockAction.EXPECT().IsValid().Return(true, nil)
			mockAction.EXPECT().Perform().Return(nil)

			ok, err := performer.Perform()
			Expect(ok).To(BeTrue())
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return error when action is not allowed", func() {
			mockAction.EXPECT().Context().Return(ctx)

			mockAction.EXPECT().IsAllowed().Return(false, nil)
			//mockAction.EXPECT().Performer().Return("TestPerformer")

			ok, err := performer.Perform()
			Expect(ok).To(BeFalse())
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(&action.AuthorizationError{}))
		})

		It("should return error when action is not enabled", func() {
			mockAction.EXPECT().Context().Return(ctx)

			mockAction.EXPECT().IsAllowed().Return(true, nil)
			mockAction.EXPECT().IsEnabled().Return(false, action.ErrorMap{"error": "action not enabled"})

			ok, err := performer.Perform()
			Expect(ok).To(BeFalse())
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(&action.DisabledError{}))
		})

		It("should return error when action is not valid", func() {
			mockAction.EXPECT().Context().Return(ctx)

			mockAction.EXPECT().IsAllowed().Return(true, nil)
			mockAction.EXPECT().IsEnabled().Return(true, nil)
			mockAction.EXPECT().IsValid().Return(false, action.ErrorMap{"error": "action not valid"})

			ok, err := performer.Perform()
			Expect(ok).To(BeFalse())
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(&action.ValidationError{}))
		})

		It("should return error when perform fails", func() {
			mockAction.EXPECT().Context().Return(ctx)

			mockAction.EXPECT().IsAllowed().Return(true, nil)
			mockAction.EXPECT().IsEnabled().Return(true, nil)
			mockAction.EXPECT().IsValid().Return(true, nil)
			mockAction.EXPECT().Perform().Return(errors.New("perform failed"))

			ok, err := performer.Perform()
			Expect(ok).To(BeFalse())
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("PerformIfEnabled", func() {
		It("should perform action successfully when enabled", func() {
			mockAction.EXPECT().Context().Return(ctx)

			mockAction.EXPECT().IsAllowed().Return(true, nil)
			mockAction.EXPECT().IsEnabled().Return(true, nil)
			mockAction.EXPECT().IsValid().Return(true, nil)
			mockAction.EXPECT().Perform().Return(nil)

			ok, err := performer.Try()
			Expect(ok).To(BeTrue())
			Expect(err).ToNot(HaveOccurred())
		})

		It("should skip action without error when not enabled", func() {
			mockAction.EXPECT().Context().Return(ctx)

			mockAction.EXPECT().IsAllowed().Return(true, nil)
			mockAction.EXPECT().IsEnabled().Return(false, action.ErrorMap{"error": "action not enabled"})

			ok, err := performer.Try()
			Expect(ok).To(BeFalse())
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return error when action is not valid", func() {
			mockAction.EXPECT().Context().Return(ctx)

			mockAction.EXPECT().IsAllowed().Return(true, nil)
			mockAction.EXPECT().IsEnabled().Return(true, nil)
			mockAction.EXPECT().IsValid().Return(false, action.ErrorMap{"error": "action not valid"})

			ok, err := performer.Try()
			Expect(ok).To(BeFalse())
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(&action.ValidationError{}))
		})

		It("should return error when perform fails", func() {
			mockAction.EXPECT().Context().Return(ctx)

			mockAction.EXPECT().IsAllowed().Return(true, nil)
			mockAction.EXPECT().IsEnabled().Return(true, nil)
			mockAction.EXPECT().IsValid().Return(true, nil)
			mockAction.EXPECT().Perform().Return(errors.New("perform failed"))

			ok, err := performer.Try()
			Expect(ok).To(BeFalse())
			Expect(err).To(HaveOccurred())
		})
	})
})
