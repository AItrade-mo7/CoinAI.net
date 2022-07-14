package main

import (
	_ "embed"

	"Hunter.net/server/global"
	"Hunter.net/server/global/config"
	"Hunter.net/server/router"
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
