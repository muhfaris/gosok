package provider

import (
	uuid "github.com/satori/go.uuid"
	"golang.org/x/oauth2"
)

type (
	InitOption struct {
		GoogleConfig *oauth2.Config
	}
)

var (
	gc *oauth2.Config
)

func Init(opt InitOption) {
	googleConfig = &oauth2.Config{
		RedirectURL:  opt.GoogleConfig.RedirectURL,
		Endpoint:     opt.GoogleConfig.Endpoint,
		Scopes:       opt.GoogleConfig.Scopes,
		ClientID:     opt.GoogleConfig.ClientID,
		ClientSecret: opt.GoogleConfig.ClientSecret,
	}
	googleStateString = uuid.NewV4()
}
