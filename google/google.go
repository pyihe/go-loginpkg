package google

import (
	"fmt"
	"net/http"

	"github.com/pyihe/go-loginpkg"
	"github.com/pyihe/go-pkg/https"
	"github.com/pyihe/go-pkg/serialize"
	jsonserialize "github.com/pyihe/go-pkg/serialize/json"
)

const Name = "google"

type Response struct {
	Iss     string `json:"iss"`     //
	Aud     string `json:"aud"`     //
	Sub     string `json:"sub"`     // 用户在google的唯一标示
	Name    string `json:"name"`    // 名字
	Picture string `json:"picture"` // 头像
	Iat     string `json:"iat"`     //
	Exp     string `json:"exp"`     // google token过期时间
}

type google struct {
}

func (g google) Auth(req interface{}) (result interface{}, err error) {
	var token, ok = req.(string)
	if !ok {
		err = loginpkg.ErrInvalidRequest
		return
	}

	var data Response
	var url = fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", token)
	err = https.GetWithObj(http.DefaultClient, url, serialize.Get(jsonserialize.Name), &data)
	result = data
	return
}

func init() {
	loginpkg.Register(Name, google{})
}
