package taptap

import (
	"crypto"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pyihe/secret"

	"github.com/pyihe/go-loginpkg"
)

const (
	ParamClientId    = "client_id"
	ParamMacKey      = "mac_key"
	ParamAccessToken = "access_token"
)

type validator struct {
	hasher secret.Hasher
}

func newValidator() loginpkg.Validator {
	return validator{hasher: secret.NewHasher()}
}

func (c validator) Verify(req loginpkg.Request) (result loginpkg.Response, err error) {
	const domain = "openapi.taptap.com"

	uri := fmt.Sprintf("/account/profile/v1?client_id=%s", req.Get(ParamClientId))
	now := time.Now()
	unix := now.Unix()
	nonce := now.UnixNano()
	signStr := fmt.Sprintf("%d\n%d\n%s\n%s\n%s\n443\n\n", unix, nonce, http.MethodGet, uri, domain)
	mac := c.hasher.MAC(crypto.SHA1, []byte(signStr), []byte(req.Get(ParamMacKey)))
	macToken := base64.StdEncoding.EncodeToString(mac)
	auth := fmt.Sprintf("MAC id=\"%s\",ts=\"%d\",nonce=\"%d\",mac=\"%s\"", req.Get(ParamAccessToken), unix, nonce, macToken)

	httpRequest, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://%s%s", domain, uri), nil)
	if err != nil {
		return
	}
	httpRequest.Header.Set("Authorization", auth)

	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return
	}
	defer httpResponse.Body.Close()

	var tapResponse struct {
		Now     int64 `json:"now,omitempty"`
		Success bool  `json:"success,omitempty"`
		Data    struct {
			Avatar  string `json:"avatar,omitempty"`
			Gender  string `json:"gender,omitempty"`
			Name    string `json:"name,omitempty"`
			OpenId  string `json:"openid,omitempty"`
			UnionId string `json:"unionid,omitempty"`
		} `json:"data,omitempty"`
	}
	err = json.NewDecoder(httpResponse.Body).Decode(&tapResponse)
	if err != nil {
		return
	}

	sex := 0
	switch strings.ToLower(tapResponse.Data.Gender) {
	case "1", "male":
		sex = 1
	case "2", "female":
		sex = 2
	}
	result = loginpkg.Response{
		Nickname: tapResponse.Data.Name,
		Avatar:   tapResponse.Data.Avatar,
		Gender:   sex,
		OpenId:   tapResponse.Data.OpenId,
		UnionId:  tapResponse.Data.UnionId,
	}
	return
}

func init() {
	loginpkg.Register(loginpkg.TapTap, newValidator())
}
