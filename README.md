### go-loginpkg
third party login auth check package written in go.

### usage
```go
package main

import (
    "github.com/pyihe/go-loginpkg"
    "github.com/pyihe/go-loginpkg/wechat"
)

func main() {
    var request = wechat.Request{
        AppID:     "your appid",
        AppSecret: "your secret",
        Code:      "your auth code",
    }
    var rsp, err = loginpkg.GetChecker(wechat.Name).Auth(request)
    if err != nil {
        //handle err
    }
    var reply, ok = rsp.(wechat.Response)
    
    // handle logic
    _, _ = reply, ok
}
```
