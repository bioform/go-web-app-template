package main

import (
	"context"
	"errors"

	"github.com/bioform/go-web-app-template/internal/user/model"
	"github.com/bioform/go-web-app-template/pkg/action"
	"github.com/bioform/go-web-app-template/pkg/api"
	"github.com/bioform/go-web-app-template/pkg/database"
	_ "github.com/bioform/go-web-app-template/pkg/database" // init() in database.go
	"github.com/bioform/go-web-app-template/pkg/dbaction"
	validator "github.com/rezakhademix/govalidator/v2"
	"gorm.io/gorm"

	"github.com/bioform/go-web-app-template/pkg/logging"
)

func main() {
	ctx := api.New(database.Default()).AddTo(context.Background())

	log := logging.Logger(ctx)

	// Prepare the action.
	ap := action.New(ctx, &MyAction{
		SomeAttr: "", // Set the action-specific attribute.
	})

	// Perform the action.
	ok, err := ap.Perform()
	if !ok {

		var actionError action.ActionError
		if errors.As(err, &actionError) {
			log.Error("Action error", "error", err)
		} else {
			log.Error("Error", "error", err)
		}
	} else {
		log.Info("Action performed successfully", "SomeAttr", ap.Action().SomeAttr)
	}
}

/////////////////////////////////////////////////////////////////////////////////////////
// Example Action implementation
/////////////////////////////////////////////////////////////////////////////////////////

// Define a specific type embedding BaseAction.
type MyAction struct {
	dbaction.BaseAction

	SomeAttr string
}

// Implement the specific behavior for MyAction.
func (a *MyAction) Perform() error {
	ctx := a.Context()
	// Put your business logic here.
	log := logging.Logger(ctx)
	api, err := api.From(ctx)
	if err != nil {
		return err
	}
	db := api.DB()

	log.Info("MyAction-specific perform logic")

	user := &model.User{Name: "L1212", Email: "mmm@example.com", Password: "123456"}
	err = db.Create(user).Error
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

func (a *MyAction) IsValid() (bool, error) {
	v := validator.New()
	v.RequiredString(a.SomeAttr, "SomeAttr", "required")
	return v.IsPassed(), action.ErrorMap(v.Errors())
}
