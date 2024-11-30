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

// func TestActionPerformer(t *testing.T) {
// 	RegisterFailHandler(Fail)
// 	RunSpecs(t, "ActionPerformer Suite")
// }

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
		mockAction.On("TransactionProvider").Return(mockTransactionProvider).Maybe()

		call := mockTransactionProvider.EXPECT().Transaction(mock.Anything, mock.Anything).Maybe()
		call.Run(func(args mock.Arguments) {
			ctx := args.Get(0).(context.Context)
			lambda := args.Get(1).(func(context.Context) error)
			err := lambda(ctx)
			call.Return(err)
		})

		performer = action.New(mockAction)
	})

	Describe("Action", func() {
		It("should return the action", func() {
			Expect(performer.Action()).To(Equal(mockAction))
		})
	})

	Describe("Perform", func() {
		It("should perform action successfully", func() {
			mockAction.EXPECT().IsAllowed(ctx).Return(true, nil)
			mockAction.EXPECT().IsEnabled(ctx).Return(true, nil)
			mockAction.EXPECT().IsValid(ctx).Return(true, nil)
			mockAction.EXPECT().Perform(ctx).Return(nil)

			ok, err := performer.Perform(ctx)
			Expect(ok).To(BeTrue())
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return error when action is not enabled", func() {
			mockAction.EXPECT().IsAllowed(ctx).Return(true, nil)
			mockAction.EXPECT().IsEnabled(ctx).Return(false, action.ErrorMap{"error": "action not enabled"})

			ok, err := performer.Perform(ctx)
			Expect(ok).To(BeFalse())
			Expect(err).To(HaveOccurred())
		})

		It("should return error when action is not valid", func() {
			mockAction.EXPECT().IsAllowed(ctx).Return(true, nil)
			mockAction.EXPECT().IsEnabled(ctx).Return(true, nil)
			mockAction.EXPECT().IsValid(ctx).Return(false, action.ErrorMap{"error": "action not valid"})

			ok, err := performer.Perform(ctx)
			Expect(ok).To(BeFalse())
			Expect(err).To(HaveOccurred())
		})

		It("should return error when perform fails", func() {
			mockAction.EXPECT().IsAllowed(ctx).Return(true, nil)
			mockAction.EXPECT().IsEnabled(ctx).Return(true, nil)
			mockAction.EXPECT().IsValid(ctx).Return(true, nil)
			mockAction.EXPECT().Perform(ctx).Return(errors.New("perform failed"))

			ok, err := performer.Perform(ctx)
			Expect(ok).To(BeFalse())
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("PerformIfEnabled", func() {
		It("should perform action successfully when enabled", func() {
			mockAction.EXPECT().IsAllowed(ctx).Return(true, nil)
			mockAction.EXPECT().IsEnabled(ctx).Return(true, nil)
			mockAction.EXPECT().IsValid(ctx).Return(true, nil)
			mockAction.EXPECT().Perform(ctx).Return(nil)

			ok, err := performer.PerformIfEnabled(ctx)
			Expect(ok).To(BeTrue())
			Expect(err).ToNot(HaveOccurred())
		})

		It("should perform action without error when not enabled but ifEnabled is true", func() {
			mockAction.EXPECT().IsAllowed(ctx).Return(true, nil)
			mockAction.EXPECT().IsEnabled(ctx).Return(false, action.ErrorMap{"error": "action not enabled"})

			ok, err := performer.PerformIfEnabled(ctx)
			Expect(ok).To(BeFalse())
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return error when action is not valid", func() {
			mockAction.EXPECT().IsAllowed(ctx).Return(true, nil)
			mockAction.EXPECT().IsEnabled(ctx).Return(true, nil)
			mockAction.EXPECT().IsValid(ctx).Return(false, action.ErrorMap{"error": "action not valid"})

			ok, err := performer.PerformIfEnabled(ctx)
			Expect(ok).To(BeFalse())
			Expect(err).To(HaveOccurred())
		})

		It("should return error when perform fails", func() {
			mockAction.EXPECT().IsAllowed(ctx).Return(true, nil)
			mockAction.EXPECT().IsEnabled(ctx).Return(true, nil)
			mockAction.EXPECT().IsValid(ctx).Return(true, nil)
			mockAction.EXPECT().Perform(ctx).Return(errors.New("perform failed"))

			ok, err := performer.PerformIfEnabled(ctx)
			Expect(ok).To(BeFalse())
			Expect(err).To(HaveOccurred())
		})
	})
})
