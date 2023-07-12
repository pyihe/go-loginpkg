package apple

import (
	"github.com/pyihe/apple_validator"

	"github.com/pyihe/go-loginpkg"
)

const (
	ParamIdentityToken = "identity_token"
)

type validator struct {
	*apple_validator.Validator
}

func newValidator() loginpkg.Validator {
	return validator{Validator: apple_validator.NewValidator()}
}

func (a validator) Verify(req loginpkg.Request) (loginpkg.Response, error) {
	jwtToken, err := a.CheckIdentityToken(req.Get(ParamIdentityToken))
	if err != nil {
		return loginpkg.NilResponse, err
	}
	if ok, err := jwtToken.IsValid(); err != nil {
		return loginpkg.NilResponse, err
	} else if !ok {
		return loginpkg.NilResponse, loginpkg.ErrExpired
	}
	return loginpkg.Response{
		Avatar:   "",
		Gender:   0,
		Nickname: "",
		OpenId:   jwtToken.Sub(),
		UnionId:  "",
	}, nil
}

func init() {
	loginpkg.Register(loginpkg.Apple, newValidator())
}
