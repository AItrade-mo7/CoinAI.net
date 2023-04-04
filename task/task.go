package main

import (
	_ "embed"
	"fmt"
	"os"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
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
	backObj := testHunter.New(testHunter.TestOpt{
		StartTime: mTime.TimeParse(mTime.Lay_DD, "2023-01-01"),
		EndTime:   mTime.TimeParse(mTime.Lay_DD, "2023-04-05"),
		InstID:    "BTC-USDT",
	})
	err := backObj.StuffDBKdata()
	if err != nil {
		panic(fmt.Errorf("出错: %+v", err))
	}
	err = backObj.CheckKdataList() // 检查数据是否出错
	if err != nil {
		panic(fmt.Errorf("出错: %+v", err))
	}

	// 新建回测参数序列
	configArr := testHunter.GetConfig(testHunter.GetConfigOpt{
		EmaPArr: []int{75, 77, 79, 169, 171, 173, 543, 545, 547},
		CAPArr:  []int{3, 4},
	})

	for _, Gor := range configArr {
		for _, config := range Gor {
			fmt.Println(config)
		}
	}

	// 新建参数测试
	// MockObj := backObj.NewMock(testHunter.NewMockOpt{
	// 	MockName:  "123123",
	// 	InitMoney: "1000",
	// 	Level:     "1",
	// 	Charge:    "0.05",
	// 	TradeKdataOpt: hunter.TradeKdataOpt{
	// 		MA_Period:      77,
	// 		RSI_Period:     18,
	// 		RSI_EMA_Period: 14,
	// 		CAP_Period:     4,
	// 	},
	// })
	// MockObj.MockRun()

	select {}
}
