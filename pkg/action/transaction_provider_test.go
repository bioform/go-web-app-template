package action

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type mockTransactionProvider struct{}

func (m *mockTransactionProvider) Transaction(currentContext context.Context, executeInTransaction func(newContext context.Context) error) error {
	return nil
}

var _ = Describe("TransactionProvider", func() {
	var (
		provider *mockTransactionProvider
	)

	BeforeEach(func() {
		provider = &mockTransactionProvider{}
	})

	Describe("SetTransactionProvider", func() {
		It("should set the transaction provider", func() {
			SetTransactionProvider(provider)
			Expect(transactionProvider).To(Equal(provider))
		})
	})
})
