package instagram

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/pyihe/go-pkg/https"
	"github.com/pyihe/go-pkg/serialize"
	jsonserialize "github.com/pyihe/go-pkg/serialize/json"

	"github.com/pyihe/go-loginpkg"
)

const (
	ParamClientId     = "client_id"
	ParamClientSecret = "client_secret"
	ParamRedirectURI  = "redirect_uri"
	ParamCode         = "code"
)

type instagram struct{}

func (ins instagram) Verify(req loginpkg.Request) (loginpkg.Response, error) {
	var (
		clientId     = req.Get(ParamClientId)
		clientSecret = req.Get(ParamClientSecret)
		redirectURI  = req.Get(ParamRedirectURI)
		code         = req.Get(ParamCode)
		url          = "https://api.instagram.com/oauth/access_token"
		param        = fmt.Sprintf("client_id=%s&client_secret=%s&grant_type=authorization_code&redirect_uri=%s&code=%s", clientId, clientSecret, redirectURI, strings.TrimRight(code, "#_"))
		data         struct {
			AccessToken string `json:"access_token"`
			UserId      int64  `json:"user_id"`
			Username    string `json:"username"`
		}
	)

	err := https.PostWithObj(http.DefaultClient, url, "application/x-www-form-urlencoded", strings.NewReader(param), serialize.Get(jsonserialize.Name), &data)
	if err != nil {
		return loginpkg.NilResponse, err
	}
	url = fmt.Sprintf("https://graph.instagram.com/%d?fields=username&access_token=%s", data.UserId, data.AccessToken)
	err = https.GetWithObj(http.DefaultClient, url, serialize.Get(jsonserialize.Name), &data)
	if err != nil {
		return loginpkg.NilResponse, err
	}
	return loginpkg.Response{
		Avatar:   "",
		Gender:   0,
		Nickname: data.Username,
		OpenId:   strconv.FormatInt(data.UserId, 10),
		UnionId:  "",
	}, nil
}

func init() {
	loginpkg.Register(loginpkg.Instagram, instagram{})
}
