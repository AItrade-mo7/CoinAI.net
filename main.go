package main

import (
	_ "embed"
	"fmt"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/okxApi/restApi"
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

	resData, err := restApi.Fetch(restApi.FetchOpt{
		Path: "/abc/ert",
		Data: map[string]any{
			"qwe": 123,
			"abc": 456,
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

	// router.Start()
}
