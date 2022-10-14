package mOKX

import (
	"github.com/EasyGolang/goTools/mEncrypt"
	"github.com/EasyGolang/goTools/mFetch"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/EasyGolang/goTools/mTime"
)

/*
testCase.OKXWss()
*/

type LoginArgsType struct {
	APIKey     string `bson:"apiKey"`
	Passphrase string `bson:"passphrase"`
	Timestamp  string `bson:"timestamp"`
	Sign       string `bson:"sign"`
}
type LoginType struct {
	Op   string          `bson:"op"`
	Args []LoginArgsType `bson:"args"`
}

type OptWssOKX struct {
	FetchType int
	Event     func(string, any)
	OKXKey    TypeOkxKey
}

func WssOKX(opt OptWssOKX) (_this *mFetch.Wss) {
	WssOpt := mFetch.WssOpt{}
	WssOpt.Event = opt.Event
	if opt.FetchType == 0 {
		WssOpt.Url = "wss://ws.okx.com:8443/ws/v5/public"
	}
	if opt.FetchType == 1 {
		WssOpt.Url = "wss://ws.okx.com:8443/ws/v5/private"
	}
	_this = mFetch.NewWss(WssOpt)

	if opt.FetchType == 1 {
		Timestamp := mTime.EpochTime()
		SignStr := mStr.Join(
			Timestamp,
			"GET",
			"/users/self/verify",
		)
		Sign := mEncrypt.Sha256(SignStr, opt.OKXKey.SecretKey)
		LoginParam := LoginType{
			Op: "login",
			Args: []LoginArgsType{
				{
					APIKey:     opt.OKXKey.ApiKey,
					Passphrase: opt.OKXKey.Passphrase,
					Timestamp:  Timestamp,
					Sign:       Sign,
				},
			},
		}
		_this.Write(mJson.ToJson(LoginParam))
	}

	return
}
