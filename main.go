package main

import (
	_ "embed"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/ready"
	"CoinAI.net/server/router"
	jsoniter "github.com/json-iterator/go"
)

//go:embed package.json
var AppPackage []byte

func main() {
	jsoniter.Unmarshal(AppPackage, &config.AppInfo)
	// 初始化系统参数
	global.Start()

	// restApi.Fetch(restApi.FetchOpt{
	// 	Path: "/abc/ert",
	// 	Data: map[string]any{
	// 		"qwe": 123,
	// 		"abc": 456,
	// 	},
	// 	Event: func(s string, a any) {
	// 		fmt.Println(s, a)
	// 	},
	// }).Get()

	ready.Start()

	router.Start()
}
