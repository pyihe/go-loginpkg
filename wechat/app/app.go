package app

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/pyihe/go-loginpkg"
	"github.com/pyihe/go-pkg/https"
	"github.com/pyihe/go-pkg/serialize"
	jsonserialize "github.com/pyihe/go-pkg/serialize/json"
)

const Name = "wechat"

type Request struct {
	AppID     string
	AppSecret string
	Code      string
}

type Response struct {
	Sex         int    `json:"sex"`               // 性别, 0: 未知, 1: 男, 3: 女
	OpenId      string `json:"openid"`            // openid
	AccessToken string `json:"access_token"`      // access_token
	NickName    string `json:"nickname"`          // 昵称
	Avatar      string `json:"headimgurl"`        // 头像
	UnionId     string `json:"unionid,omitempty"` // unionid
}

type wechat struct {
}

func (w wechat) Auth(req interface{}) (result interface{}, err error) {
	var r, ok = req.(Request)
	if !ok {
		err = loginpkg.ErrInvalidRequest
		return
	}

	var data Response
	var url = fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%v&secret=%v&code=%v&grant_type=authorization_code", r.AppID, r.AppSecret, r.Code)
	err = https.GetWithObj(http.DefaultClient, url, serialize.Get(jsonserialize.Name), &data)
	if err != nil {
		return
	}

	if data.OpenId == "" || data.AccessToken == "" {
		err = errors.New("invalid code")
		return
	}

	url = fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%v&openid=%v&lang=zh_CN", data.AccessToken, data.OpenId)
	err = https.GetWithObj(http.DefaultClient, url, serialize.Get(jsonserialize.Name), &data)
	result = data
	return
}

func init() {
	loginpkg.Register(Name, wechat{})
}
