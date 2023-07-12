package loginpkg

import (
	"errors"
	"sync"
)

const (
	WeChatApp  = iota + 1 // 微信公众号
	WechatMini            // 微信小程序
	QQ                    // QQ登录
	TapTap                // TapTap登录
	Apple                 // Apple登录
	Google                // Google登录
	Facebook              // Facebook登录
	Instagram             // ins登录
	Twitter               // Twitter登录
)

type Request map[string]string

func NewRequest(params ...map[string]string) Request {
	req := make(Request)
	if len(params) > 0 {
		for _, param := range params {
			for name, value := range param {
				req[name] = value
			}
		}
	}
	return req
}

func (req Request) Add(name, value string) {
	req[name] = value
}

func (req Request) Del(name string) {
	delete(req, name)
}

func (req Request) Get(name string) (value string) {
	return req[name]
}

type Response struct {
	Avatar   string // 头像
	Gender   int    // 性别
	Nickname string // 昵称
	OpenId   string // 第三方唯一标识
	UnionId  string // 联合ID
}

type Validator interface {
	Verify(Request) (Response, error)
}

var (
	ErrExpired     = errors.New("authorization expired")
	ErrUnknownCode = errors.New("unknown code")
)

var (
	locker      = sync.RWMutex{}
	m           = make(map[int]Validator)
	NilResponse = Response{}
)

func Register(t int, validator Validator) {
	locker.Lock()
	m[t] = validator
	locker.Unlock()
}

func GetValidator(t int) (v Validator) {
	locker.RLock()
	v = m[t]
	locker.RUnlock()
	return
}
