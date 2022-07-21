package main

import (
	_ "embed"
	"fmt"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/okxApi/wssApi"
	"CoinAI.net/server/ready"
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
		Type: 1,
		Event: func(s string, a any) {
			fmt.Println("Event", s, mStr.ToStr(a))
		},
	})

	wss.Read(func(msg []byte) {
		fmt.Println("读数据", mStr.ToStr(msg))
	})

	// router.Start()
}
