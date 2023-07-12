package google

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/pyihe/go-pkg/https"
	"github.com/pyihe/go-pkg/serialize"
	jsonserialize "github.com/pyihe/go-pkg/serialize/json"

	"github.com/pyihe/go-loginpkg"
)

const (
	ParamIdToken = "id_token"
)

type Response struct {
	Iss     string `json:"iss"`     //
	Aud     string `json:"aud"`     //
	Sub     string `json:"sub"`     // 用户在google的唯一标示
	Name    string `json:"name"`    // 名字
	Picture string `json:"picture"` // 头像
	Iat     string `json:"iat"`     //
	Exp     string `json:"exp"`     // google token过期时间
}

type google struct{}

func (g google) Verify(req loginpkg.Request) (loginpkg.Response, error) {
	var (
		url  = fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", req.Get(ParamIdToken))
		data struct {
			Iss     string `json:"iss"`     //
			Aud     string `json:"aud"`     //
			Sub     string `json:"sub"`     // 用户在google的唯一标示
			Name    string `json:"name"`    // 名字
			Picture string `json:"picture"` // 头像
			Iat     string `json:"iat"`     //
			Exp     string `json:"exp"`     // google token过期时间
		}
	)
	err := https.GetWithObj(http.DefaultClient, url, serialize.Get(jsonserialize.Name), &data)
	if err != nil {
		return loginpkg.NilResponse, err
	}
	exp, err := strconv.ParseInt(data.Exp, 10, 64)
	if err != nil {
		return loginpkg.NilResponse, err
	}
	if exp <= time.Now().Unix() {
		return loginpkg.NilResponse, loginpkg.ErrExpired
	}
	return loginpkg.Response{
		Avatar:   data.Picture,
		Gender:   0,
		Nickname: data.Name,
		OpenId:   data.Sub,
		UnionId:  "",
	}, nil
}

func init() {
	loginpkg.Register(loginpkg.Google, google{})
}
