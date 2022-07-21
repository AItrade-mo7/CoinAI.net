package restApi

import (
	"fmt"
	"strings"

	"CoinAI.net/server/global/dbData"
	"github.com/EasyGolang/goTools/mEncrypt"
	"github.com/EasyGolang/goTools/mFetch"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/EasyGolang/goTools/mTime"
	"github.com/EasyGolang/goTools/mUrl"
)

/*
	resData, err := restApi.Fetch(restApi.FetchOpt{
		Path: "/api/v5/account/balance",
		Data: map[string]any{
			"ccy": "USDT",
		},
		Method: "get",
		Event: func(s string, a any) {
			fmt.Println("Event", s, a)
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(mStr.ToStr(resData))

*/

type FetchOpt struct {
	Path   string
	Data   map[string]any
	Method string
	Event  func(string, any)
}

func Fetch(opt FetchOpt) (resData []byte, resErr error) {
	if len(opt.Method) < 1 {
		opt.Method = "GET"
	}

	// 处理 Header 和 加密信息
	Method := strings.ToUpper(opt.Method)
	Timestamp := mTime.IsoTime(true)
	ApiKey := dbData.OkxKey.ApiKey
	SecretKey := dbData.OkxKey.SecretKey
	Passphrase := dbData.OkxKey.Passphrase
	Body := mJson.ToJson(opt.Data)

	SignStr := mStr.Join(
		Timestamp,
		Method,
		opt.Path,
		string(Body),
	)

	if Method == "GET" {
		Body = []byte("")
		urlO := mUrl.InitUrl(opt.Path)
		for key, val := range opt.Data {
			v := fmt.Sprintf("%+v", val)
			urlO.AddParam(key, v)
		}
		signPath := urlO.String()
		SignStr = mStr.Join(
			Timestamp,
			Method,
			signPath,
			string(Body),
		)
	}
	Sign := mEncrypt.Sha256(SignStr, SecretKey)
	fetch := mFetch.NewHttp(mFetch.HttpOpt{
		Origin: "https://www.okx.com",
		Path:   opt.Path,
		Data:   opt.Data,
		Event:  opt.Event,
		Header: map[string]string{
			"OK-ACCESS-KEY":        ApiKey,
			"OK-ACCESS-SIGN":       Sign,
			"OK-ACCESS-TIMESTAMP":  Timestamp,
			"OK-ACCESS-PASSPHRASE": Passphrase,
		},
	})

	if Method == "GET" {
		return fetch.Get()
	} else {
		return fetch.Post()
	}
}
