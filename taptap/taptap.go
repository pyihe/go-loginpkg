package taptap

import (
	"crypto"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pyihe/go-loginpkg"
	"github.com/pyihe/secret"
)

type Response struct {
	Name    string `json:"name,omitempty"`
	Avatar  string `json:"avatar,omitempty"`
	OpenId  string `json:"openid,omitempty"`
	UnionId string `json:"unionid,omitempty"`
}

type Request struct {
	ClientId    string
	AccessToken string
	MacKey      string
}

const Name = "taptap"

type checker struct {
}

func (c checker) Auth(req interface{}) (result interface{}, err error) {
	request, ok := req.(Request)
	if !ok {
		err = loginpkg.ErrInvalidRequest
		return
	}

	const domain = "openapi.taptap.com"

	uri := fmt.Sprintf("/account/profile/v1?client_id=%s", request.ClientId)
	now := time.Now()
	unix := now.Unix()
	nonce := now.UnixNano()
	signStr := fmt.Sprintf("%d\n%d\n%s\n%s\n%s\n443\n\n", unix, nonce, http.MethodGet, uri, domain)
	mac := secret.NewHasher().MAC(crypto.SHA1, []byte(signStr), []byte(request.MacKey))
	macToken := base64.StdEncoding.EncodeToString(mac)
	auth := fmt.Sprintf("MAC id=\"%s\",ts=\"%d\",nonce=\"%d\",mac=\"%s\"", request.AccessToken, unix, nonce, macToken)

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

	response := Response{}
	err = json.NewDecoder(httpResponse.Body).Decode(&response)
	result = response
	return
}

func init() {
	loginpkg.Register(Name, checker{})
}
