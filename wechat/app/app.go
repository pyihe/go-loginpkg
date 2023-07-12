package app

import (
	"fmt"
	"net/http"

	"github.com/pyihe/go-pkg/https"
	"github.com/pyihe/go-pkg/serialize"
	jsonserialize "github.com/pyihe/go-pkg/serialize/json"

	"github.com/pyihe/go-loginpkg"
)

const (
	ParamAppId     = "app_id"
	ParamAppSecret = "app_secret"
	ParamCode      = "code"
)

type wechat struct{}

func (w wechat) Verify(req loginpkg.Request) (loginpkg.Response, error) {
	var (
		appID     = req.Get(ParamAppId)
		appSecret = req.Get(ParamAppSecret)
		code      = req.Get(ParamCode)
		url       = fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%v&secret=%v&code=%v&grant_type=authorization_code", appID, appSecret, code)
		data      struct {
			Sex         int    `json:"sex"`               // 性别, 0: 未知, 1: 男, 3: 女
			OpenId      string `json:"openid"`            // openid
			AccessToken string `json:"access_token"`      // access_token
			NickName    string `json:"nickname"`          // 昵称
			Avatar      string `json:"headimgurl"`        // 头像
			UnionId     string `json:"unionid,omitempty"` // unionid
		}
	)
	err := https.GetWithObj(http.DefaultClient, url, serialize.Get(jsonserialize.Name), &data)
	if err != nil {
		return loginpkg.NilResponse, err
	}

	if data.OpenId == "" || data.AccessToken == "" {
		return loginpkg.NilResponse, loginpkg.ErrUnknownCode
	}

	url = fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%v&openid=%v&lang=zh_CN", data.AccessToken, data.OpenId)
	err = https.GetWithObj(http.DefaultClient, url, serialize.Get(jsonserialize.Name), &data)
	if err != nil {
		return loginpkg.NilResponse, err
	}
	return loginpkg.Response{
		Avatar:   data.Avatar,
		Gender:   data.Sex,
		Nickname: data.NickName,
		OpenId:   data.OpenId,
		UnionId:  data.UnionId,
	}, nil
}

func init() {
	loginpkg.Register(loginpkg.WeChatApp, wechat{})
}
