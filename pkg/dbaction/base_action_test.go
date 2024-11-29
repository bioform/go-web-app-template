package dbaction_test

import (
	"context"
	"sync"
	"testing"

	"github.com/bioform/go-web-app-template/pkg/action"
	"github.com/bioform/go-web-app-template/pkg/database"
	"github.com/bioform/go-web-app-template/pkg/dbaction"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type MyAction struct {
	dbaction.BaseAction
}

func (a *MyAction) Perform(ctx context.Context) error {
	return nil
}

func TestMyActionTransaction(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MyAction Transaction Suite")
}

var _ = Describe("MyAction", func() {
	var (
		ctx              context.Context
		db               *gorm.DB
		spyDbProvider    *SpyDbProvider
		transactionCount int
		mu               sync.Mutex
	)

	BeforeEach(func(specContext SpecContext) {
		var err error
		db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		Expect(err).NotTo(HaveOccurred())

		transactionCount = 0

		spyDbProvider = &SpyDbProvider{
			DbProvider: database.NewDbProvider(db),
			OnTransaction: func() {
				mu.Lock()
				transactionCount++
				mu.Unlock()
			},
		}

		ctx = specContext
	})

	It("should perform each action in a separate transaction", func() {
		performAction := func() {
			actionInstance := &MyAction{}
			actionInstance.SetTransactionProvider(spyDbProvider)

			performer := action.New(actionInstance)
			ok, err := performer.Perform(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(ok).To(BeTrue())
		}

		// Perform the action twice
		performAction()
		performAction()

		// Verify that two separate transactions were used
		mu.Lock()
		count := transactionCount
		mu.Unlock()
		Expect(count).To(Equal(2))
	})
})

// SpyDbProvider wraps DbProvider to count Transaction calls
type SpyDbProvider struct {
	*database.DbProvider
	OnTransaction func()
}

func (s *SpyDbProvider) Transaction(ctx context.Context, fn func(newCtx context.Context) error) error {
	if s.OnTransaction != nil {
		s.OnTransaction()
	}
	return s.DbProvider.Transaction(ctx, fn)
}
