package main

import (
	_ "embed"

	"CoinAI.net/server/global"
	"CoinAI.net/task/taskStart"
)

func main() {
	// 初始化系统参数
	global.Start()

	taskStart.BackTest("ETH-USDT")
	taskStart.BackTest("BTC-USDT")
}
