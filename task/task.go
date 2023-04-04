package main

import (
	_ "embed"
	"fmt"
	"os"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/hunter"
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
		StartTime: mTime.TimeParse(mTime.Lay_DD, "2020-01-01"),
		EndTime:   mTime.TimeParse(mTime.Lay_DD, "2023-03-31"),
		InstID:    "BTC-USDT",
	})
	err := backObj.StuffDBKdata()
	if err != nil {
		fmt.Println("出错", err)
	}
	err = backObj.CheckKdataList() // 检查数据是否出错
	if err != nil {
		fmt.Println("出错", err)
	}

	// configArr := testHunter.MockConfig([]int{75, 77, 79, 169, 171, 173, 543, 545, 547})

	backObj.NewMock(testHunter.NewMockOpt{
		MockName:  "123123",
		InitMoney: "1000",
		Level:     "1",
		Charge:    "0.05",
		TradeKdataOpt: hunter.TradeKdataOpt{
			MA_Period:      77,
			RSI_Period:     18,
			RSI_EMA_Period: 14,
			CAP_Period:     4,
		},
	})

	// for _, config := range configArr {
	// 	back.MockData(
	// 		config.MockOpt,
	// 		config.TradeKdataOpt,
	// 	)
	// }
	// taskPush.SysEmail(taskPush.SysEmailOpt{
	// 	From:        config.SysName,
	// 	To:          config.NoticeEmail,
	// 	Subject:     "参数跑完了",
	// 	Title:       "第一批参数组合跑完了",
	// 	Content:     "参数值:" + mJson.Format(configArr),
	// 	Description: "回测结束通知",
	// })
}
