package reqDataCenter

import (
	"time"

	"CoinAI.net/server/global/config"
	"github.com/EasyGolang/goTools/mEncrypt"
	"github.com/EasyGolang/goTools/mFetch"
)

/*

resData := reqDataCenter.NewRest(reqDataCenter.NewRestOpt{
	Path:   "/private/get_user_info",
	Method: "GET",
	Data:   map[string]any{},
})
fmt.Println(mStr.ToStr(resData))


*/

type RestOpt struct {
	Origin string
	Path   string
	UserID string
	Method string
	Data   map[string]any
}

func NewRest(opt RestOpt) (resData []byte, resErr error) {
	Token := mEncrypt.NewToken(mEncrypt.NewTokenOpt{
		SecretKey: config.SecretKey,              // key
		ExpiresAt: time.Now().Add(time.Hour / 2), // 过期时间 半小时
		Message:   opt.UserID,
		Issuer:    "AITrade.net",
		Subject:   "UserToken",
	}).Generate()

	UserAgent := "CoinAI.net"

	fetch := mFetch.NewHttp(mFetch.HttpOpt{
		Origin: opt.Origin,
		Path:   opt.Path,
		Data:   opt.Data,
		Header: map[string]string{
			"Auth-Encrypt": config.ClientEncrypt(opt.Path + UserAgent),
			"Auth-Token":   Token,
			"User-Agent":   UserAgent,
		},
	})

	if opt.Method == "GET" {
		return fetch.Get()
	} else {
		return fetch.Post()
	}
}
