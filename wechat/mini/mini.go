package mini

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"

	"github.com/pyihe/go-loginpkg"
	"github.com/pyihe/go-pkg/https"
	"github.com/pyihe/go-pkg/serialize"
	jsonserialize "github.com/pyihe/go-pkg/serialize/json"
)

const Name = "wechat_mini"

type Request struct {
	AppID       string
	AppSecret   string
	Code        string
	EncryptData string
	Iv          string
}

type Response struct {
	OpenId   string
	NickName string
	Gender   int
	UnionId  string
	Avatar   string
}

type weMini struct {
}

func init() {
	loginpkg.Register(Name, weMini{})
}

func (m weMini) Auth(arg interface{}) (result interface{}, err error) {
	var wxData struct {
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
	var r, ok = arg.(Request)
	if !ok {
		err = loginpkg.ErrInvalidRequest
		return
	}

	sessionKey, err := m.getSessionKeyAndOpenID(r.AppID, r.AppSecret, r.Code)
	if err != nil {
		return
	}

	encryptData, err := base64.StdEncoding.DecodeString(r.EncryptData)
	if err != nil {
		return
	}
	iv, err := base64.StdEncoding.DecodeString(r.Iv)
	if err != nil {
		return
	}

	realData, err := aES128CBCDecrypt(encryptData, sessionKey, iv)
	if err != nil {
		return
	}

	if err = serialize.Get(jsonserialize.Name).Unmarshal(realData, &wxData); err != nil {
		return
	}
	response := Response{
		OpenId:   wxData.OpenId,
		NickName: wxData.NickName,
		Gender:   wxData.Gender,
		UnionId:  wxData.UnionId,
		Avatar:   wxData.AvatarUrl,
	}
	return response, nil
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
