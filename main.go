package main

import (
	_ "embed"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/okxApi/wssApi"
	"CoinAI.net/server/ready"
	"CoinAI.net/server/router"
	"github.com/EasyGolang/goTools/mStr"
	jsoniter "github.com/json-iterator/go"
)

//go:embed package.json
var AppPackage []byte

func main() {
	jsoniter.Unmarshal(AppPackage, &config.AppInfo)
	// 初始化系统参数
	global.Start()

	ready.Start()

	wss := wssApi.New(wssApi.FetchOpt{
		Type: 0,
		Event: func(s string, a any) {
			global.WssLog.Println("Event", s, mStr.ToStr(a))
		},
	})

	go wss.Read(func(msg []byte) {
		global.WssLog.Println("读数据", mStr.ToStr(msg))
	})

	router.Start()
}
