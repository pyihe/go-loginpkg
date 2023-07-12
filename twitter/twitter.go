package twitter

import (
	"context"

	twittersdk "github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"

	"github.com/pyihe/go-loginpkg"
)

const (
	ParamClientId     = "client_id"
	ParamClientSecret = "client_secret"
	ParamToken        = "token"
	ParamTokenSecret  = "token_secret"
)

type twitter struct{}

func (t twitter) Verify(req loginpkg.Request) (loginpkg.Response, error) {
	var oauthConfig = oauth1.NewConfig(req.Get(ParamClientId), req.Get(ParamClientSecret))
	var tokenConfig = oauth1.NewToken(req.Get(ParamToken), req.Get(ParamTokenSecret))
	var client = twittersdk.NewClient(oauthConfig.Client(context.Background(), tokenConfig))
	var param = &twittersdk.AccountVerifyParams{
		IncludeEntities: twittersdk.Bool(true),
		SkipStatus:      twittersdk.Bool(true),
		IncludeEmail:    twittersdk.Bool(false),
	}

	info, _, err := client.Accounts.VerifyCredentials(param)
	if err != nil {
		return loginpkg.NilResponse, err
	}
	if info == nil {
		return loginpkg.NilResponse, loginpkg.ErrExpired
	}
	return loginpkg.Response{
		Avatar:   info.ProfileImageURL,
		Gender:   0,
		Nickname: info.Name,
		OpenId:   info.IDStr,
		UnionId:  "",
	}, nil
}

func init() {
	loginpkg.Register(loginpkg.Twitter, twitter{})
}
