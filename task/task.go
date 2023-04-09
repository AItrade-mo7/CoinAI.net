package main

import (
	_ "embed"

	"CoinAI.net/server/global"
	"CoinAI.net/task/analyConfig"
	"CoinAI.net/task/taskStart"
	"CoinAI.net/task/testHunter"
	"github.com/EasyGolang/goTools/mTime"
)

var ResultBasePath = "/root/AItrade/CoinAI.net/task/analyConfig/data"

func main() {
	// 初始化系统参数
	global.Start()

	BackAnaly()
}

func BackAnaly() {
	EndTime := mTime.GetUnixInt64()
	StartTime := EndTime - (mTime.UnixTimeInt64.Day * 260)

	InstID := "BTC-USDT"
	BTCResult := taskStart.BackTest(taskStart.BackOpt{
		StartTime: StartTime,
		EndTime:   EndTime,
		InstID:    InstID,
		GetConfigOpt: testHunter.GetConfigOpt{
			EmaPArr:  []int{76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88},
			CAPArr:   []int{2, 3, 4, 5},
			LevelArr: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
	})
	BTCResult.ResultBasePath = ResultBasePath
	analyConfig.GetWinArr(
		BTCResult,
		// taskStart.BackReturn{
		// 	InstID:         InstID,
		// 	BillingPath:    mStr.Join(ResultBasePath, "/", InstID, "-BillingArr.json"),
		// 	ResultBasePath: ResultBasePath,
		// },
	)

	InstID = "ETH-USDT"
	ETHResult := taskStart.BackTest(taskStart.BackOpt{
		StartTime: StartTime,
		EndTime:   EndTime,
		InstID:    InstID,
		GetConfigOpt: testHunter.GetConfigOpt{
			EmaPArr:  []int{76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88},
			CAPArr:   []int{2, 3, 4, 5},
			LevelArr: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
	})
	ETHResult.ResultBasePath = ResultBasePath
	analyConfig.GetWinArr(
		ETHResult,
		// taskStart.BackReturn{
		// 	InstID:         InstID,
		// 	BillingPath:    mStr.Join(ResultBasePath, "/", InstID, "-BillingArr.json"),
		// 	ResultBasePath: ResultBasePath,
		// },
	)
}
