package twitter

import (
	"context"

	twittersdk "github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/pyihe/go-loginpkg"
)

const Name = "twitter"

type Response = twittersdk.User

type Request struct {
	ClientID     string
	ClientSecret string
	Token        string
	TokenSecret  string
}

type twitter struct {
}

func (t twitter) Auth(req interface{}) (result interface{}, err error) {
	var r, ok = req.(Request)
	if !ok {
		err = loginpkg.ErrInvalidRequest
		return
	}
	var oauthConfig = oauth1.NewConfig(r.ClientID, r.ClientSecret)
	var tokenConfig = oauth1.NewToken(r.Token, r.TokenSecret)
	var client = twittersdk.NewClient(oauthConfig.Client(context.Background(), tokenConfig))
	var param = &twittersdk.AccountVerifyParams{
		IncludeEntities: twittersdk.Bool(true),
		SkipStatus:      twittersdk.Bool(true),
		IncludeEmail:    twittersdk.Bool(false),
	}

	var user *twittersdk.User
	user, _, err = client.Accounts.VerifyCredentials(param)
	result = user
	return
}

func init() {
	loginpkg.Register(Name, twitter{})
}
