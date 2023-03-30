package main

import (
	_ "embed"
	"fmt"
	"os"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/okxInfo"
	"CoinAI.net/task/testHunter"
	"github.com/EasyGolang/goTools/mTime"
	jsoniter "github.com/json-iterator/go"
)

func main() {
	AppPackage, _ := os.ReadFile("package.json")
	jsoniter.Unmarshal(AppPackage, &config.AppInfo)
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

	back.MockData(
		testHunter.BillingType{
			MockName:  "EMA_108",
			InitMoney: "1000",
			Level:     "1",
		},
		okxInfo.TradeKdataOpt{
			MA_Period:      108,
			RSI_Period:     18,
			RSI_EMA_Period: 14,
			CAP_Period:     3,
		},
	)

	// back.MockData(
	// 	testHunter.BillingType{
	// 		MockName:  "EMA_107",
	// 		InitMoney: "1000",
	// 		Level:     "1",
	// 	},
	// 	okxInfo.TradeKdataOpt{
	// 		MA_Period:      107,
	// 		RSI_Period:     18,
	// 		RSI_EMA_Period: 14,
	// 		CAP_Period:     3,
	// 	},
	// )
}
