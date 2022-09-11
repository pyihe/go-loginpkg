package instagram

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pyihe/go-loginpkg"
	"github.com/pyihe/go-pkg/https"
	"github.com/pyihe/go-pkg/serialize"
	jsonserialize "github.com/pyihe/go-pkg/serialize/json"
)

const Name = "instagram"

type Request struct {
	ClientID     string
	ClientSecret string
	RedirectUri  string
	Code         string
}

type Response struct {
	AccessToken string `json:"access_token"`
	UserId      int64  `json:"user_id"`
	Username    string `json:"username"`
}

type instagram struct {
}

func (ins instagram) Auth(req interface{}) (result interface{}, err error) {
	var r, ok = req.(Request)
	if !ok {
		err = loginpkg.ErrInvalidRequest
		return
	}

	var data Response
	var url = "https://api.instagram.com/oauth/access_token"
	var param = fmt.Sprintf("client_id=%s&client_secret=%s&grant_type=authorization_code&redirect_uri=%s&code=%s", r.ClientID, r.ClientSecret, r.RedirectUri, strings.TrimRight(r.Code, "#_"))

	err = https.PostWithObj(http.DefaultClient, url, "application/x-www-form-urlencoded", strings.NewReader(param), serialize.Get(jsonserialize.Name), &data)
	if err != nil {
		return
	}

	url = fmt.Sprintf("https://graph.instagram.com/%d?fields=username&access_token=%s", data.UserId, data.AccessToken)
	err = https.GetWithObj(http.DefaultClient, url, serialize.Get(jsonserialize.Name), &data)
	result = data
	return
}

func init() {
	loginpkg.Register(Name, instagram{})
}
