package router

import (
	"database/sql"

	"github.com/gorilla/sessions"
	"github.com/spf13/viper"
	"github.com/thedevsaddam/renderer"
)

type (
	SessionOptions struct {
		Domain   string
		Path     string
		MaxAge   int
		HttpOnly bool
	}

	InitOption struct {
		Rnd     *renderer.Render
		Session *SessionOptions
	}
)

var (
	sessOpt *SessionOptions
	rnd     *renderer.Render
	dbPool  *sql.DB

	store = sessions.NewCookieStore([]byte(viper.GetString("session.key")))
)

func Init(db *sql.DB, opt InitOption) {
	rnd = opt.Rnd
	dbPool = db

	store.Options = &sessions.Options{
		Domain:   opt.Session.Domain,
		Path:     opt.Session.Path,
		MaxAge:   opt.Session.MaxAge,
		HttpOnly: opt.Session.HttpOnly,
	}

}
