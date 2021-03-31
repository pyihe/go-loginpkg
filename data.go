package loginpkg

/**************Wechat****************/
type WechatResponse struct {
	Sex         int    `json:"sex"`          // 性别, 0: 未知, 1: 男, 3: 女
	OpenId      string `json:"openid"`       // openid
	AccessToken string `json:"access_token"` // access_token
	NickName    string `json:"nickname"`     // 昵称
	Avatar      string `json:"headimgurl"`   // 头像
}

/**************Google****************/
type GoogleResponse struct {
	Iss     string `json:"iss"`     //
	Aud     string `json:"aud"`     //
	Sub     string `json:"sub"`     // 用户在google的唯一标示
	Name    string `json:"name"`    // 名字
	Picture string `json:"picture"` // 头像
	Iat     string `json:"iat"`     //
	Exp     string `json:"exp"`     // google token过期时间
}

/**************Apple****************/
type AppleToken struct {
	headerStr string
	claimsStr string
	sign      string       //签名
	header    *AppleHeader //header
	claims    *AppleClaim  //claims
}
type AppleHeader struct {
	Kid string `json:"kid"` //apple公钥的密钥ID
	Alg string `json:"alg"` //签名token的算法
}

type AppleClaim struct {
	Iss            string `json:"iss"`   //签发者，固定值: https://appleid.apple.com
	Sub            string `json:"sub"`   //用户唯一标识
	Aud            string `json:"aud"`   //App ID
	Iat            int64  `json:"iat"`   //token生成时间
	Exp            int64  `json:"exp"`   //token过期时间
	Nonce          string `json:"nonce"` //客户端设置的随机值
	NonceSupported bool   `json:"nonce_supported"`
	Email          string `json:"email"` //邮件
	EmailVerified  bool   `json:"email_verified"`
	IsPrivateEmail bool   `json:"is_private_email"`
	RealUserStatus int    `json:"real_user_status"`
	CHash          string `json:"c_hash"`    //
	AuthTime       int64  `json:"auth_time"` //验证时间
}

/**************Facebook****************/
type FacebookResponse struct {
	Name    string `json:"name"`
	Id      string `json:"id"`
	Picture struct {
		Data struct {
			Height       int    `json:"height"`
			Width        int    `json:"width"`
			IsSilhouette bool   `json:"is_silhouette"`
			Url          string `json:"url"`
		} `json:"data"`
	} `json:"picture"`
}

/**************Instagram****************/
type InstagramResponse struct {
	AccessToken string `json:"access_token"`
	UserId      int64  `json:"user_id"`
	Username    string `json:"username"`
}
