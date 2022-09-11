package facebook

import (
	"fmt"
	"net/http"

	"github.com/pyihe/go-loginpkg"
	"github.com/pyihe/go-pkg/https"
	"github.com/pyihe/go-pkg/serialize"
	jsonserialize "github.com/pyihe/go-pkg/serialize/json"
)

const Name = "facebook"

type Request struct {
	ClientID     string
	ClientSecret string
	Token        string
}

type Response struct {
	Name    string `json:"name"`
	Id      string `json:"id"`
	Picture struct {
		Data struct {
			Height       int    `json:"height"`
			Width        int    `json:"width"`
			IsSilhouette bool   `json:"is_silhouette"`
			Url          string `json:"url"`
		} `json:"data"`
	} `json:"picture"`
}

type facebook struct {
}

func (f facebook) Auth(req interface{}) (result interface{}, err error) {
	var r, ok = req.(Request)
	if !ok {
		err = loginpkg.ErrInvalidRequest
		return
	}
	var fbAccessToken = fmt.Sprintf("%s|%s", r.ClientID, r.ClientSecret)
	var url = fmt.Sprintf("https://graph.facebook.com/v9.0/debug_token?access_token=%s&input_token=%s", fbAccessToken, r.Token)
	var getProfile = func(uid string) string {
		return fmt.Sprintf("https://graph.facebook.com/%s?fields=name,picture&access_token=%s", uid, fbAccessToken)
	}
	var data struct {
		Data struct {
			AppId               string `json:"app_id"`
			Type                string `json:"type"`
			Application         string `json:"application"`
			DataAccessExpiresAt int64  `json:"data_access_expires_at"`
			ExpiresAt           int64  `json:"expires_at"`
			IsValid             bool   `json:"is_valid"`
			IssuedAt            int64  `json:"issued_at"`
			Metadata            struct {
				AuthType string `json:"auth_type"`
				Sso      string `json:"sso"`
			}
			Scopes []string `json:"scopes"`
			UserId string   `json:"user_id"`
		} `json:"data"`
	}

	err = https.GetWithObj(http.DefaultClient, url, serialize.Get(jsonserialize.Name), &data)
	if err != nil {
		return
	}
	if !data.Data.IsValid {
		err = loginpkg.ErrAuthFail
		return
	}

	url = getProfile(data.Data.UserId)

	var rsp Response
	err = https.GetWithObj(http.DefaultClient, url, serialize.Get(jsonserialize.Name), &rsp)
	result = rsp
	return
}

func init() {
	loginpkg.Register(Name, facebook{})
}
