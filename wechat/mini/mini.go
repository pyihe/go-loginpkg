package mini

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"

	"github.com/pyihe/go-pkg/https"
	"github.com/pyihe/go-pkg/serialize"
	jsonserialize "github.com/pyihe/go-pkg/serialize/json"

	"github.com/pyihe/go-loginpkg"
)

const (
	ParamAppId       = "app_id"
	ParamAppSecret   = "app_secret"
	ParamCode        = "code"
	ParamEncryptData = "encrypt_data"
	ParamIv          = "iv"
)

type Response struct {
	OpenId   string
	NickName string
	Gender   int
	UnionId  string
	Avatar   string
}

type weMini struct{}

func init() {
	loginpkg.Register(loginpkg.WechatMini, weMini{})
}

func (m weMini) Verify(req loginpkg.Request) (loginpkg.Response, error) {
	var (
		appId       = req.Get(ParamAppId)
		appSecret   = req.Get(ParamAppSecret)
		code        = req.Get(ParamCode)
		encryptData = req.Get(ParamEncryptData)
		iv          = req.Get(ParamIv)
		wxData      struct {
			OpenId    string `json:"openId,omitempty"`
			NickName  string `json:"nickName,omitempty"`
			Gender    int    `json:"gender,omitempty"`
			City      string `json:"city,omitempty"`
			Province  string `json:"province,omitempty"`
			Country   string `json:"country,omitempty"`
			AvatarUrl string `json:"avatarUrl,omitempty"`
			UnionId   string `json:"unionId,omitempty"`
			WaterMark struct {
				AppId string `json:"appid,omitempty"`
			} `json:"watermark,omitempty"`
		}
	)
	sessionKey, err := m.getSessionKeyAndOpenID(appId, appSecret, code)
	if err != nil {
		return loginpkg.NilResponse, err
	}

	encryptBytes, err := base64.StdEncoding.DecodeString(encryptData)
	if err != nil {
		return loginpkg.NilResponse, err
	}
	ivBytes, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return loginpkg.NilResponse, err
	}

	realData, err := aES128CBCDecrypt(encryptBytes, sessionKey, ivBytes)
	if err != nil {
		return loginpkg.NilResponse, err
	}

	if err = serialize.Get(jsonserialize.Name).Unmarshal(realData, &wxData); err != nil {
		return loginpkg.NilResponse, err
	}
	return loginpkg.Response{
		Avatar:   wxData.AvatarUrl,
		Gender:   wxData.Gender,
		Nickname: wxData.NickName,
		OpenId:   wxData.OpenId,
		UnionId:  wxData.UnionId,
	}, nil
}

func (m weMini) getSessionKeyAndOpenID(appId, secret, code string) (sessionKey []byte, err error) {
	var url = fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%v&secret=%v&js_code=%v&grant_type=authorization_code", appId, secret, code)
	var serializer = serialize.Get(jsonserialize.Name)
	var session struct {
		Key     string `json:"session_key,omitempty"`
		UnionId string `json:"unionid,omitempty"`
		ErrMsg  string `json:"errmsg,omitempty"`
		OpenId  string `json:"openid,omitempty"`
		ErrCode int    `json:"errcode,omitempty"`
	}

	if err = https.GetWithObj(http.DefaultClient, url, serializer, &session); err != nil {
		return
	}

	switch session.ErrCode {
	case 40029:
		err = errors.New("code无效")
	case 45011:
		err = errors.New("api minute-quota reach limit  mustslower  retry next minute")
	case 40226:
		err = errors.New("code blocked")
	case -1:
		err = errors.New("wechat system error")
	case 0:
	}

	if err != nil {
		return
	}

	sessionKey, err = base64.StdEncoding.DecodeString(session.Key)
	if err != nil {
		return
	}
	return
}

func aES128CBCDecrypt(encryptData, key, iv []byte) (origData []byte, err error) {
	defer func() {
		if pErr := recover(); pErr != nil {
			err = errors.New("key is not match data")
		}
	}()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData = make([]byte, len(encryptData))
	blockMode.CryptBlocks(origData, encryptData)
	origData = pKCS7UnPadding(origData)
	return origData, nil
}

func pKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
