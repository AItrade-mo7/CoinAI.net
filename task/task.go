package main

import (
	_ "embed"

	"CoinAI.net/server/global"
	"CoinAI.net/task/analyConfig"
	"CoinAI.net/task/taskStart"
	"github.com/EasyGolang/goTools/mStr"
)

var ResultBasePath = "/root/AItrade/CoinAI.net/task/analyConfig/data"

func main() {
	// 初始化系统参数
	global.Start()

	// BackBtc := taskStart.BackTest(taskStart.BackOpt{
	// 	StartTime: 3132131,
	// 	EndTime:   31313,
	// 	InstID:    "BTC-USDT",
	// 	GetConfigOpt: testHunter.GetConfigOpt{
	// 		EmaPArr:  []int{}, // Ema 步长
	// 		CAPArr:   []int{}, // CAP 步长
	// 		LevelArr: []int{}, // 杠杆倍数
	// 	},
	// })

	InstID := "BTC-USDT"
	analyConfig.GetWinArr(taskStart.BackReturn{
		InstID:         InstID,
		BillingPath:    mStr.Join(ResultBasePath, "/", InstID, "-BillingArr.json"),
		ResultBasePath: ResultBasePath,
	})
	// BackEth := taskStart.BackTest(taskStart.BackOpt{
	// 	StartTime: 3132131,
	// 	EndTime:   31313,
	// 	InstID:    "BTC-USDT",
	// 	GetConfigOpt: testHunter.GetConfigOpt{
	// 		EmaPArr:  []int{}, // Ema 步长
	// 		CAPArr:   []int{}, // CAP 步长
	// 		LevelArr: []int{}, // 杠杆倍数
	// 	},
	// })

	InstID = "ETH-USDT"
	analyConfig.GetWinArr(taskStart.BackReturn{
		InstID:         InstID,
		BillingPath:    mStr.Join(ResultBasePath, "/", InstID, "-BillingArr.json"),
		ResultBasePath: ResultBasePath,
	})
}
