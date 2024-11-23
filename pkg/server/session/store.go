package session

import (
	"context"
	"log"

	"github.com/alexedwards/scs/gormstore"
	"github.com/alexedwards/scs/v2"
	"github.com/bioform/go-web-app-template/internal/user/model"
	"github.com/bioform/go-web-app-template/pkg/database"
)

var Manager *scs.SessionManager

type userKey string

const (
	UserIdKey string  = "user_id"
	UserKey   userKey = "user"
)

func init() {
	db := database.GetDefault(context.Background())
	var err error

	Manager = scs.New()
	Manager.Cookie.Name = "_session_id"

	if Manager.Store, err = gormstore.New(db); err != nil {
		log.Fatal(err)
	}
}

func GetUser(ctx context.Context) *model.User {
	user := ctx.Value(UserKey)
	if user == nil {
		return nil
	}
	return user.(*model.User)
}
