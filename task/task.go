package main

import (
	_ "embed"
	"os"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/task/taskStart"
	jsoniter "github.com/json-iterator/go"
)

func main() {
	AppPackage, _ := os.ReadFile("package.json")
	jsoniter.Unmarshal(AppPackage, &config.AppInfo)
	// 初始化系统参数
	global.Start()

	taskStart.BackTest()
}
