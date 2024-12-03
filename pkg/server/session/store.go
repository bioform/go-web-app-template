package session

import (
	"log"

	"github.com/alexedwards/scs/gormstore"
	"github.com/alexedwards/scs/v2"
	"github.com/bioform/go-web-app-template/pkg/database"
)

var Manager *scs.SessionManager

type userKey string

const (
	UserIdKey string  = "user_id"
	UserKey   userKey = "user"
)

func init() {
	db := database.Default()
	var err error

	Manager = scs.New()
	Manager.Cookie.Name = "_session_id"

	if Manager.Store, err = gormstore.New(db); err != nil {
		log.Fatal(err)
	}
}
