package restApi

import (
	"github.com/EasyGolang/goTools/mFetch"
	"github.com/EasyGolang/goTools/mJson"
	jsoniter "github.com/json-iterator/go"
)

type FetchOpt struct {
	Path  string
	Data  any
	Event func(string, any)
}

func Fetch(opt FetchOpt) (_this *mFetch.Http) {
	mFetchOpt := mFetch.HttpOpt{}
	mFetchOpt.Origin = "https://www.okx.com"
	mFetchOpt.Path = opt.Path
	mFetchOpt.Event = opt.Event
	if mFetchOpt.Event == nil {
		mFetchOpt.Event = func(s1 string, s2 any) {}
	}
	// 处理data
	jsonByte := mJson.ToJson(opt.Data)
	jsoniter.Unmarshal(jsonByte, &mFetchOpt.Data)
	// 处理 Header
	mFetchOpt.Header = map[string]string{
		// "OK-ACCESS-KEY":        info.APIKey,
		// "OK-ACCESS-SIGN":       info.Sign,
		// "OK-ACCESS-TIMESTAMP":  info.Timestamp,
		// "OK-ACCESS-PASSPHRASE": info.Passphrase,
	}

	_this = mFetch.NewHttp(mFetchOpt)
	return
}
