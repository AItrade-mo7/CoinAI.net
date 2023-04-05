package main

import (
	_ "embed"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/hunter"
	"CoinAI.net/server/okxInfo"
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

	// 数据准备
	ready.Start()

	SetHunter()

	// 启动路由
	router.Start()
}

func SetHunter() {
	// 新建策略
	MyHunter := hunter.New(hunter.HunterOpt{
		HunterName: "MyHunter",
		TradeKdataOpt: okxInfo.TradeKdataOpt{
			MA_Period:      171,
			RSI_Period:     18,
			RSI_EMA_Period: 14,
			CAP_Period:     4,
		},
	})
	err := MyHunter.SetTradeInst()
	if err != nil {
		global.LogErr(err)
		return
	}

	MyHunter.Start()
}
