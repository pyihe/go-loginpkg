package facebook

import (
	"fmt"
	"net/http"

	"github.com/pyihe/go-pkg/https"
	"github.com/pyihe/go-pkg/serialize"
	jsonserialize "github.com/pyihe/go-pkg/serialize/json"

	"github.com/pyihe/go-loginpkg"
)

const (
	ParamClientId     = "client_id"
	ParamClientSecret = "client_secret"
	ParamInputToken   = "input_token"
)

type facebook struct{}

func (f facebook) Verify(req loginpkg.Request) (loginpkg.Response, error) {
	var fbAccessToken = fmt.Sprintf("%s|%s", req.Get(ParamClientId), req.Get(ParamClientSecret))
	var url = fmt.Sprintf("https://graph.facebook.com/v9.0/debug_token?access_token=%s&input_token=%s", fbAccessToken, req.Get(ParamInputToken))
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

	err := https.GetWithObj(http.DefaultClient, url, serialize.Get(jsonserialize.Name), &data)
	if err != nil {
		return loginpkg.NilResponse, err
	}
	if !data.Data.IsValid {
		return loginpkg.NilResponse, loginpkg.ErrExpired
	}

	url = getProfile(data.Data.UserId)

	var rsp struct {
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
	err = https.GetWithObj(http.DefaultClient, url, serialize.Get(jsonserialize.Name), &rsp)
	if err != nil {
		return loginpkg.NilResponse, err
	}

	return loginpkg.Response{
		Avatar:   rsp.Picture.Data.Url,
		Gender:   0,
		Nickname: rsp.Name,
		OpenId:   data.Data.UserId,
		UnionId:  "",
	}, nil
}

func init() {
	loginpkg.Register(loginpkg.Facebook, facebook{})
}
