package main

import (
	_ "embed"
	"fmt"

	"CoinAI.net/server/global"
	"CoinAI.net/server/hunter/testHunter"
	"CoinAI.net/task/taskHunter"
	"github.com/EasyGolang/goTools/mFile"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/EasyGolang/goTools/mTime"
)

var ResultBasePath = "/root/AItrade/CoinAI.net/task/2022_2023"

func main() {
	// 初始化系统参数
	global.Start()

	// Step1("BTC-USDT")
	// Step2("BTC-USDT")
	Step3("BTC-USDT")
	// Step4("BTC-USDT")

	// Step1("ETH-USDT")
	// Step2("ETH-USDT")
	// Step3("ETH-USDT")
	// Step4("ETH-USDT")
}

func Step1(InstID string) {
	// 第一步： 暴力求值 （海量参数结果罗列） 需要几个小时甚至好几天
	EmaPArr := []int{}
	StarNum := 30
	for i := 0; i < 400; i += 2 {
		StarNum = 30 + i
		EmaPArr = append(EmaPArr, StarNum)
		fmt.Println(StarNum)
	}

	// LevelArr := []int{2, 3, 4, 5}
	// ConfArr := []dbType.TradeKdataOpt{}
	// for _, l := range LevelArr {
	// 	conf := dbType.TradeKdataOpt{
	// 		EMA_Period: 342,
	// 		CAP_Period: 7,
	// 		CAP_Max:    "2.5",
	// 	}
	// 	conf.MaxTradeLever = l
	// 	ConfArr = append(ConfArr, conf)
	// }

	StartTime := mTime.TimeParse(mTime.Lay_DD, "2021-01-01")
	EndTime := mTime.TimeParse(mTime.Lay_DD, "2022-01-01")
	taskHunter.BackTest(taskHunter.BackOpt{
		StartTime: StartTime,
		EndTime:   EndTime,
		InstID:    InstID,
		OutPutDir: mStr.Join(ResultBasePath),
		GetConfigOpt: testHunter.GetConfigOpt{
			EmaPArr:  EmaPArr,
			CAPArr:   []int{3},
			LevelArr: []int{1},
			CAPMax:   []string{"0.5"},
			CAPMin:   []string{"-0.5"},
			// ConfArr: ConfArr,
		},
	})
}

func Step2(InstID string) {
	// 第二步骤：根据胜率和最终营收 筛选
	taskHunter.GetWinArr(taskHunter.GetWinArrOpt{
		InstID:     InstID,
		OutPutDir:  ResultBasePath,
		MoneyRight: "3000",
		// WinRight:   "0.35",
		// Sort: "Win",
	})
}

func Step3(InstID string) {
	// 第三步：提取第二步的配置，加上杠杆得出新的参数组合 大概 几百个 然后 换个新的时间段进行新一轮测试
	confArr := taskHunter.GetWinConfig(taskHunter.GetWinConfigOpt{
		OutPutDir: ResultBasePath,
		InstID:    InstID,
	})

	// 提取 EMA 的值

	EmaArr := []int{}

	// ConfArr := []dbType.TradeKdataOpt{}
	// LevelArr := []int{2, 3, 4, 5}

	for _, conf := range confArr {
		EmaArr = append(EmaArr, conf.EMA_Period)
	}
	mFile.Write(mStr.Join(
		ResultBasePath, "/", InstID, "-EmaArr.json",
	), mJson.ToStr(EmaArr))

	// 新一轮求解，计算最优杠杆倍率 用  2022 年 8 月 的 260 天前进行回测 （此步骤会更换时间段反复进行）
	// EndTime := mTime.TimeParse(mTime.Lay_DD, "2022-10-01")
	// EndTime := mTime.TimeParse(mTime.Lay_DD, "2023-05-01")
	// StartTime := EndTime - (mTime.UnixTimeInt64.Day * 260)
	// taskHunter.BackTest(taskHunter.BackOpt{
	// 	StartTime: StartTime,
	// 	EndTime:   EndTime,
	// 	InstID:    InstID,
	// 	OutPutDir: mStr.Join(ResultBasePath, "/Step3"),
	// 	GetConfigOpt: testHunter.GetConfigOpt{
	// 		// EmaPArr:  []int{342},      // Ema 步长
	// 		// CAPArr:   []int{7},        // CAP 步长
	// 		// CAPMax:   []string{"2.5"}, // CAPMax 步长
	// 		// CAPMin:   []string{"-2.5"},
	// 		// LevelArr: []int{2},
	// 		ConfArr: []dbType.TradeKdataOpt{
	// 			{
	// 				EMA_Period:    342,
	// 				CAP_Period:    7,
	// 				CAP_Max:       "2.5",
	// 				CAP_Min:       "-2.5",
	// 				MaxTradeLever: 5,
	// 			},
	// 		},
	// 	},
	// })
}

func Step4(InstID string) {
	// 第四步： 根据第三步的结果进行筛选 (胜率和最终营收) 得出参数结果
	taskHunter.GetWinArr(taskHunter.GetWinArrOpt{
		InstID: InstID,
		// OutPutDir:  mStr.Join(ResultBasePath),
		OutPutDir: mStr.Join(ResultBasePath, "/Step3"),
		// MoneyRight: "1000",
		// WinRight:   "0.35",
		// Sort: "Win", //  Win
	})
}
