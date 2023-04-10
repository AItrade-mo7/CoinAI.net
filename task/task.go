package main

import (
	_ "embed"

	"CoinAI.net/server/global"
	"CoinAI.net/server/okxInfo"
	"CoinAI.net/task/taskStart"
	"CoinAI.net/task/testHunter"
	"github.com/EasyGolang/goTools/mTime"
)

var ResultBasePath = "/root/AItrade/CoinAI.net/task/analyConfig/最近8个月2"

func main() {
	// 初始化系统参数
	global.Start()

	BackAnaly()
}

var EmaPArr = []int{}

var (
	CAPArr = []int{
		2,
		3,
		4,
		5,
		6,
	}
	LevelArr = []int{1}
	CAPMax   = []string{
		"0.5",
		"1",
		"1.5",
		"2",
		"2.5",
		"3",
	}

	ConfArr = []okxInfo.TradeKdataOpt{
		{
			EMA_Period:    70,
			CAP_Period:    2,
			CAP_Max:       "0.2",
			MaxTradeLever: 1,
		},
		{
			EMA_Period:    542,
			CAP_Period:    2,
			CAP_Max:       "0.2",
			MaxTradeLever: 1,
		},
	}
)

func BackAnaly() {
	// StarNum := 70
	// for i := 0; i < 520; i += 2 {
	// 	StarNum = 60 + i
	// 	EmaPArr = append(EmaPArr, StarNum)
	// }

	// fmt.Println(mJson.ToStr(EmaPArr))

	EndTime := mTime.GetUnixInt64()
	StartTime := EndTime - (mTime.UnixTimeInt64.Day * 260)

	InstID := "BTC-USDT"
	BTCResult := taskStart.BackTest(taskStart.BackOpt{
		StartTime: StartTime,
		EndTime:   EndTime,
		InstID:    InstID,
		GetConfigOpt: testHunter.GetConfigOpt{
			EmaPArr:  EmaPArr,
			CAPArr:   CAPArr,
			LevelArr: LevelArr,
			CAPMax:   CAPMax,
			ConfArr:  ConfArr,
		},
	})
	BTCResult.ResultBasePath = ResultBasePath
	// analyConfig.GetWinArr(
	// 	BTCResult,
	// 	// taskStart.BackReturn{
	// 	// 	InstID:         InstID,
	// 	// 	BillingPath:    mStr.Join(ResultBasePath, "/", InstID, "-BillingArr.json"),
	// 	// 	ResultBasePath: ResultBasePath,
	// 	// },
	// )

	InstID = "ETH-USDT"
	ETHResult := taskStart.BackTest(taskStart.BackOpt{
		StartTime: StartTime,
		EndTime:   EndTime,
		InstID:    InstID,
		GetConfigOpt: testHunter.GetConfigOpt{
			EmaPArr:  EmaPArr,
			CAPArr:   CAPArr,
			LevelArr: LevelArr,
			CAPMax:   CAPMax,
			ConfArr:  ConfArr,
		},
	})
	ETHResult.ResultBasePath = ResultBasePath
	// analyConfig.GetWinArr(
	// 	ETHResult,
	// 	// taskStart.BackReturn{
	// 	// 	InstID:         InstID,
	// 	// 	BillingPath:    mStr.Join(ResultBasePath, "/", InstID, "-BillingArr.json"),
	// 	// 	ResultBasePath: ResultBasePath,
	// 	// },
	// )
}
