package main

import (
	_ "embed"
	"fmt"

	"CoinAI.net/server/global"
	"CoinAI.net/server/okxInfo"
	"CoinAI.net/task/testHunter"
	"github.com/EasyGolang/goTools/mTime"
)

func main() {
	// 初始化系统参数
	global.Start()

	// 新建回测
	back := testHunter.New(testHunter.TestOpt{
		StartTime: mTime.TimeParse(mTime.Lay_DD, "2023-01-01"),
		EndTime:   mTime.TimeParse(mTime.Lay_DD, "2023-02-01"),
		InstID:    "BTC-USDT",
	})
	err := back.StuffDBKdata()
	if err != nil {
		fmt.Println("出错", err)
	}
	err = back.CheckKdataList() // 检查数据是否出错
	if err != nil {
		fmt.Println("出错", err)
	}

	back.MockData("Ema108", okxInfo.TradeKdataOpt{
		MA_Period:      108,
		RSI_Period:     18,
		RSI_EMA_Period: 14,
		CAP_Period:     3,
	})
}
