package main

import (
	_ "embed"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/global/dbType"
	"CoinAI.net/server/hunter"
	"CoinAI.net/server/ready"
	"CoinAI.net/server/router"
	"CoinAI.net/server/utils/backTest"
	jsoniter "github.com/json-iterator/go"
)

//go:embed package.json
var AppPackage []byte

func main() {
	jsoniter.Unmarshal(AppPackage, &config.AppInfo)
	// 初始化系统参数
	global.Start()

	RunTest() // 回测

	// RunIng() // 正式服务
}

func RunTest() {
	// 数据回测
	ready.ReadUserInfo()

	start := dbType.ParseTime("2022-10-27")
	end := dbType.ParseTime("2022-12-1")

	tesObj := backTest.NewTest(backTest.TestOpt{
		StartTime: start,
		EndTime:   end,
		CcyName:   "ETH",
	})
	tesObj.GetDBKdata()
	err := tesObj.CheckKdataList()
	if err == nil {
		tesObj.MockData()
	}
}

func RunIng() {
	// 数据准备
	ready.Start()

	// 启动 hunter 计算
	go hunter.Start()

	// 启动路由
	router.Start()
}
