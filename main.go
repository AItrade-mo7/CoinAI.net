package main

import (
	_ "embed"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/utils/mFetch"
	jsoniter "github.com/json-iterator/go"
)

//go:embed package.json
var AppPackage []byte

func main() {
	jsoniter.Unmarshal(AppPackage, &config.AppInfo)
	// 初始化系统参数
	global.Start()

	data := map[string]any{
		"op": "subscribe",
		"args": []string{
			"123", "345", "mei",
		},
	}

	mFetch.NewHttp(mFetch.HttpOpt{
		Origin: "http://localhost:9000",
		Path:   "/api/ping",
		Data:   data,
		Header: map[string]string{
			"Content-Type":  "appli1arset=utf-8",
			"Content-Type1": "appl2et=utf-8",
			"Content-Type2": "applicati3et=utf-8",
			"Content-Type3": "application4t=utf-8",
		},
	})

	// router.Start()
}
