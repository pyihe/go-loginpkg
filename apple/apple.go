package apple

import (
	"errors"

	"github.com/pyihe/apple_validator"
	"github.com/pyihe/go-loginpkg"
)

const Name = "apple"

type Response = apple_validator.JWTToken

type apple struct {
	validator *apple_validator.Validator
}

func (a apple) Auth(req interface{}) (result interface{}, err error) {
	var jwtToken apple_validator.JWTToken
	var token, ok = req.(string)
	if !ok {
		err = loginpkg.ErrInvalidRequest
		return
	}

	if jwtToken, err = a.validator.CheckIdentityToken(token); err != nil {
		return
	}
	if ok, err = jwtToken.IsValid(); err != nil {
		return
	} else if !ok {
		err = errors.New("invalid token")
		return
	}
	result = jwtToken
	return
}

func init() {
	loginpkg.Register(Name, apple{validator: apple_validator.NewValidator()})
}
