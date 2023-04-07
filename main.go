package main

import (
	_ "embed"
	"time"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/hunter"
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
	BTCHunter := hunter.New(hunter.HunterOpt{
		HunterName: "BTC-CoinAI",
		InstID:     "BTC-USDT",
		Describe:   "以 BTC-USDT 交易对为主执行自动交易,支持的资金量更大,更加稳定",
	})
	BTCHunter.Start()

	time.Sleep(time.Second)

	// 新建策略
	ETHHunter := hunter.New(hunter.HunterOpt{
		HunterName: "ETH-CoinAI",
		InstID:     "ETH-USDT",
		Describe:   "以 ETH-USDT 交易对为主执行自动交易,交易次数更加频发,可以收货更高收益",
	})
	ETHHunter.Start()
}
