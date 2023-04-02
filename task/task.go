package main

import (
	_ "embed"
	"fmt"
	"os"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/hunter"
	"CoinAI.net/server/utils/taskPush"
	"CoinAI.net/task/testHunter"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/EasyGolang/goTools/mTime"
	jsoniter "github.com/json-iterator/go"
)

type MockOptType struct {
	MockOpt       testHunter.BillingType
	TradeKdataOpt hunter.TradeKdataOpt
}

func MockConfig(EmaPArr []int) []MockOptType {
	MockConfigArr := []MockOptType{}

	for _, emaP := range EmaPArr {
		MockConfigArr = append(MockConfigArr,
			MockOptType{
				testHunter.BillingType{
					MockName:  "EMA_" + mStr.ToStr(emaP),
					InitMoney: "1000", // 初始资金
					Level:     "1",    // 杠杆倍数
					Charge:    "0.05", // 吃单标准手续费率 0.05%
				},
				hunter.TradeKdataOpt{
					MA_Period:      emaP,
					RSI_Period:     18,
					RSI_EMA_Period: 14,
					CAP_Period:     3,
				},
			},
		)
	}

	return MockConfigArr
}

func main() {
	AppPackage, _ := os.ReadFile("package.json")
	jsoniter.Unmarshal(AppPackage, &config.AppInfo)
	// 初始化系统参数
	global.Start()

	// 新建回测
	back := testHunter.New(testHunter.TestOpt{
		StartTime: mTime.TimeParse(mTime.Lay_DD, "2020-01-01"),
		EndTime:   mTime.TimeParse(mTime.Lay_DD, "2023-03-31"),
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

	configArr := MockConfig([]int{75, 77, 79})
	for _, config := range configArr {
		back.MockData(
			config.MockOpt,
			config.TradeKdataOpt,
		)
	}
	taskPush.SysEmail(taskPush.SysEmailOpt{
		From:        config.SysName,
		To:          config.NoticeEmail,
		Subject:     "参数跑完了",
		Title:       "第一批参数组合跑完了",
		Content:     "参数值:" + mJson.Format(configArr),
		Description: "回测结束通知",
	})

	// configArr := MockConfig([]int{169, 171, 173})
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
	// 	Title:       "第二批参数组合跑完了",
	// 	Content:     "参数值:" + mJson.Format(configArr),
	// 	Description: "回测结束通知",
	// })

	// configArr := MockConfig([]int{543, 545, 547})
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
	// 	Title:       "第三批参数组合跑完了",
	// 	Content:     "参数值:" + mJson.Format(configArr),
	// 	Description: "回测结束通知",
	// })

	select {}
}

/*
第三轮 综合结果  545  171  77

74-79-84 附近 偏下  77,80,83

165-170-175 附近 偏上  168,171,174

542-547-552 附近 545,548,551

*/
/*
第四轮

77 - 75,77,79

171 - 169,171,173

545 - 543,545,547


*/
