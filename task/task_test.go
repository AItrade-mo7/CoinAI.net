package main

import (
	"CoinAI.net/server/global"
	"testing"
)

func TestStep1(t *testing.T) {
	// 初始化系统参数
	global.Start()

	Step1("BTC-USDT")
	// Step2("BTC-USDT")
	// Step3("BTC-USDT")
	// Step4("BTC-USDT")

	// Step1("ETH-USDT")
	// Step2("ETH-USDT")
	// Step3("ETH-USDT")
	// Step4("ETH-USDT")
}
