package loginpkg

import (
	"errors"
	"sync"

	"github.com/pyihe/go-pkg/strings"
)

type Checker interface {
	Auth(interface{}) (interface{}, error)
}

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrAuthFail       = errors.New("auth fail")
)

var (
	lock = sync.Mutex{}
	m    = make(map[string]Checker)
)

func Register(p string, handler Checker) {
	lock.Lock()
	m[strings.ToLower(p)] = handler
	lock.Unlock()
}

func GetChecker(p string) Checker {
	return m[strings.ToLower(p)]
}
