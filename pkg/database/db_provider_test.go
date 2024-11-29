package database_test

import (
	"context"
	"errors"

	"github.com/bioform/go-web-app-template/pkg/database"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var _ = Describe("DbProvider", func() {
	var (
		ctx        context.Context
		db         *gorm.DB
		dbProvider *database.DbProvider
	)

	BeforeEach(func() {
		var err error
		db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		Expect(err).NotTo(HaveOccurred())

		ctx = context.Background()
		dbProvider = database.NewDbProvider(db)
	})

	Describe("DB method", func() {
		It("should return the initial *gorm.DB when context has no DB value", func() {
			result := dbProvider.DB(ctx)
			Expect(result).To(Equal(db))
		})

		It("should return the *gorm.DB from context when set", func() {
			tx := db.Session(&gorm.Session{NewDB: true})
			ctxWithDB := dbProvider.SetDB(ctx, tx)

			result := dbProvider.DB(ctxWithDB)
			Expect(result).To(Equal(tx))
		})
	})

	Describe("SetDB method", func() {
		It("should set the *gorm.DB in the context", func() {
			tx := db.Session(&gorm.Session{NewDB: true})
			ctxWithDB := dbProvider.SetDB(ctx, tx)

			result := dbProvider.DB(ctxWithDB)
			Expect(result).To(Equal(tx))
		})
	})

	Describe("Transaction method", func() {
		It("should execute the lambda within a transaction", func() {
			err := dbProvider.Transaction(ctx, func(txCtx context.Context) error {
				txDB := dbProvider.DB(txCtx)
				Expect(txDB).NotTo(Equal(db))

				// Verify the transaction is active
				Expect(txDB.Error).To(BeNil())

				// Rollback the transaction
				return errors.New("rollback")
			})
			Expect(err).To(MatchError("rollback"))
		})
	})

	Describe("String method", func() {
		It("should return string representation with dialect name", func() {
			str := dbProvider.String()
			Expect(str).To(ContainSubstring("db: sqlite"))
		})
	})
})
