package main

import (
	_ "embed"

	"CoinFund.net/server/global"
	"CoinFund.net/server/global/config"
	"CoinFund.net/server/router"
	jsoniter "github.com/json-iterator/go"
)

//go:embed package.json
var AppPackage []byte

func main() {
	jsoniter.Unmarshal(AppPackage, &config.AppInfo)
	// 初始化系统参数
	global.Start()

	router.Start()
}
