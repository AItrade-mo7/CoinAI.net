package main

import (
	_ "embed"

	"CoinAI.net/server/global"
	"CoinAI.net/task/analyConfig"
)

func main() {
	// 初始化系统参数
	global.Start()

	// taskStart.BackTest("ETH-USDT")
	// taskStart.BackTest("BTC-USDT")

	// 开始分析和提取:

	analyConfig.AnalyBillingArr("ETH-USDT")
}
