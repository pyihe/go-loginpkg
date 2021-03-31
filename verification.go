package loginpkg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/pyihe/apple_validator"
)

// LoginByWechat 使用微信第三方登录； 返回微信用户基本信息
func LoginByWechat(appId, appSecret, code string) (*WechatResponse, error) {
	var url = fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%v&secret=%v&code=%v&grant_type=authorization_code", appId, appSecret, code)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	var data *WechatResponse
	if err = json.NewDecoder(response.Body).Decode(&data); err != nil {
		return nil, err
	}
	response.Body.Close()

	url = fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%v&openid=%v&lang=zh_CN", data.AccessToken, data.OpenId)
	response, err = http.Get(url)
	if err != nil {
		return nil, err
	}

	if err = json.NewDecoder(response.Body).Decode(&data); err != nil {
		return nil, err
	}
	response.Body.Close()

	return data, nil
}

// LoginByGoogle 使用Google第三方登录，返回Google IdToken验证结果
func LoginByGoogle(idToken string) (*GoogleResponse, error) {
	if idToken == "" {
		return nil, errors.New("empty google idToken")
	}

	var url = fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", idToken)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(response.Status)
	}

	var data *GoogleResponse
	if err = json.NewDecoder(response.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}

// LoginByApple 使用苹果第三方登录， 返回验证结果
func LoginByApple(appleToken string) (apple_validator.JWTToken, error) {
	if appleToken == "" {
		return nil, errors.New("empty apple token")
	}

	var validator = apple_validator.NewValidator()
	var jwtToken apple_validator.JWTToken
	var err error

	if jwtToken, err = validator.CheckIdentityToken(appleToken); err != nil {
		return nil, err
	}

	if ok, err := jwtToken.IsValid(); err != nil {
		return nil, err
	} else if !ok {
		return nil, errors.New("invalid apple token")
	}

	return jwtToken, nil
}

// LoginByFacebook 使用facebook第三方登录，返回用户基本信息
func LoginByFacebook(clientId, clientSecret, fbToken string) (*FacebookResponse, error) {
	var fbAccessToken = fmt.Sprintf("%s|%s", clientId, clientSecret)
	var url = fmt.Sprintf("https://graph.facebook.com/v9.0/debug_token?access_token=%s&input_token=%s", fbAccessToken, fbToken)
	var profileUrl = func(uid string) string {
		return fmt.Sprintf("https://graph.facebook.com/%s?fields=name,picture&access_token=%s", uid, fbAccessToken)
	}

	var result struct {
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

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(response.Status)
	}

	if err = json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}
	response.Body.Close()

	if result.Data.IsValid == false {
		return nil, errors.New("auth fail")
	}

	response, err = http.Get(profileUrl(result.Data.UserId))
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, errors.New(response.Status)
	}

	var fbUser *FacebookResponse
	if err = json.NewDecoder(response.Body).Decode(&fbUser); err != nil {
		return nil, err
	}
	response.Body.Close()

	return fbUser, nil
}

// LoginByInstagram Instagram第三方登录验证
func LoginByInstagram(clientId, clientSecret, redirectUri, code string) (*InstagramResponse, error) {
	//这里code需要去掉最后两个字符
	code = strings.TrimRight(code, "#_")

	var url = "https://api.instagram.com/oauth/access_token"
	var param = fmt.Sprintf("client_id=%s&client_secret=%s&grant_type=authorization_code&redirect_uri=%s&code=%s", clientId, clientSecret, redirectUri, code)

	var data *InstagramResponse

	response, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(param))
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, errors.New(response.Status)
	}
	if err = json.NewDecoder(response.Body).Decode(&data); err != nil {
		return nil, err
	}
	response.Body.Close()

	response, err = http.Get(fmt.Sprintf("https://graph.instagram.com/%d?fields=username&access_token=%s", data.UserId, data.AccessToken))
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, errors.New(response.Status)
	}
	if err = json.NewDecoder(response.Body).Decode(&data); err != nil {
		return nil, err
	}
	response.Body.Close()

	return data, nil
}

// LoginByTwitter twitter第三方登录验证，返回验证信息
func LoginByTwitter(clientId, clientSecret, twitterToken, twitterTokenSecret string) (*twitter.User, error) {
	oauthConfig := oauth1.NewConfig(clientId, clientSecret)
	tokenConfig := oauth1.NewToken(twitterToken, twitterTokenSecret)
	twitterClient := twitter.NewClient(oauthConfig.Client(context.Background(), tokenConfig))

	param := &twitter.AccountVerifyParams{
		IncludeEntities: twitter.Bool(true),
		SkipStatus:      twitter.Bool(true),
		IncludeEmail:    twitter.Bool(false),
	}

	twitterUser, _, err := twitterClient.Accounts.VerifyCredentials(param)
	if err != nil {
		return nil, err
	}

	return twitterUser, nil
}
