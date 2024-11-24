package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/bioform/go-web-app-template/internal/user/model"
	"github.com/bioform/go-web-app-template/pkg/action"
	_ "github.com/bioform/go-web-app-template/pkg/database" // init() in database.go
	"github.com/bioform/go-web-app-template/pkg/dbaction"
	"github.com/bioform/go-web-app-template/pkg/logging"
	validator "github.com/rezakhademix/govalidator/v2"
	"gorm.io/gorm"
)

// Define a specific type embedding BaseAction.
type MyAction struct {
	dbaction.BaseAction

	SomeAttr string
}

// Implement the specific behavior for MyAction.
func (a *MyAction) Perform(ctx context.Context) error {
	log := logging.Get(ctx)
	db := a.DB(ctx)

	log.Info("MyAction-specific perform logic")

	user := &model.User{Name: "L1212", Email: "mmm@example.com", Password: "123456"}
	err := db.Create(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return &model.EmailDuplicateError{Email: user.Email}
		}
		return err
	}

	var users []model.User
	db.Raw("SELECT * FROM users").Scan(&users)

	log.Info("Query users", "users", users)
	return nil
}

func (a *MyAction) IsValid(ctx context.Context) (bool, action.ErrorMap) {
	v := validator.New()
	v.RequiredString(a.SomeAttr, "SomeAttr", "required")
	return v.IsPassed(), v.Errors()
}

func main() {
	ctx := context.TODO()
	log := logging.Get(ctx)

	a := &MyAction{
		SomeAttr: "mmm", // not invalid length
	}

	// Call the template method, which handles the shared logic implicitly.
	ok, err := action.New(a).Perform(ctx)
	if !ok {
		fmt.Println("Error message: ", err)

		var validationError *action.ValidationError
		if errors.As(err, &validationError) {
			log.Error("Error details", "error", validationError.Errors())
		}
	}
}
