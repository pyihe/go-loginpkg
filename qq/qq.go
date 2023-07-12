package qq

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/pyihe/go-pkg/https"
	"github.com/pyihe/go-pkg/serialize"
	jsonserialize "github.com/pyihe/go-pkg/serialize/json"

	"github.com/pyihe/go-loginpkg"
)

const (
	ParamClientId     = "client_id"
	ParamClientSecret = "client_secret"
	ParamRedirectURI  = "redirect_uri"
	ParamAuthCode     = "code"
)

type qqValidator struct{}

func (qq qqValidator) Verify(req loginpkg.Request) (loginpkg.Response, error) {
	// QQ 授权步骤
	// 1. 前段获取授权Code，传递给后端
	// 2. 后端根据获取的code，请求access_token
	// 3. 用access_token获取open_id
	// 4. 根据open_id获取用户信息

	const (
		tokenURL = "https://graph.qq.com/oauth2.0/token"
		idURL    = "https://graph.qq.com/oauth2.0/me"
		userURL  = "https://graph.qq.com/user/get_user_info"
	)

	var (
		clientId     = req.Get(ParamClientId)
		clientSecret = req.Get(ParamClientSecret)
		redirectURI  = req.Get(ParamRedirectURI)
		code         = req.Get(ParamAuthCode)
		encoder      = serialize.Get(jsonserialize.Name)
		urlValues    = url.Values{}
		// access_token API返回字段
		tokenResponse struct {
			AccessToken  string `json:"access_token"`
			ExpiresIn    int64  `json:"expires_in"`
			RefreshToken string `json:"refresh_token"`
		}
		// openid API返回字段
		idResponse struct {
			ClientId string `json:"client_id"`
			OpenId   string `json:"openid"`
		}
		// 用户信息API返回字段
		userInfo struct {
			NickName string `json:"nickname"`
			Gender   int    `json:"gender_type"`
			Avatar   string `json:"figureurl_qq_1"`
		}
	)

	// 获取access_token
	urlValues.Add("grant_type", "authorization_code")
	urlValues.Add(ParamClientId, clientId)
	urlValues.Add(ParamClientSecret, clientSecret)
	urlValues.Add(ParamAuthCode, code)
	urlValues.Add(ParamRedirectURI, redirectURI)
	err := https.GetWithObj(http.DefaultClient, fmt.Sprintf("%s?%s", tokenURL, urlValues.Encode()), encoder, &tokenResponse)
	if err != nil {
		return loginpkg.NilResponse, err
	}

	urlValues = url.Values{}
	urlValues.Add("access_token", tokenResponse.AccessToken)
	err = https.GetWithObj(http.DefaultClient, fmt.Sprintf("%s?%s", idURL, urlValues.Encode()), encoder, &idResponse)
	if err != nil {
		return loginpkg.NilResponse, err
	}

	urlValues = url.Values{}
	urlValues.Add("access_token", tokenResponse.AccessToken)
	urlValues.Add("oauth_consumer_key", clientId)
	urlValues.Add("openid", idResponse.OpenId)
	err = https.GetWithObj(http.DefaultClient, fmt.Sprintf("%s?%s", userURL, urlValues.Encode()), encoder, &userInfo)
	if err != nil {
		return loginpkg.NilResponse, err
	}
	return loginpkg.Response{
		Avatar:   userInfo.Avatar,
		Gender:   userInfo.Gender,
		Nickname: userInfo.NickName,
		OpenId:   idResponse.OpenId,
		UnionId:  "",
	}, nil
}
