package main

import (
	_ "embed"

	"CoinAI.net/server/global"
	"CoinAI.net/server/okxInfo"
	"CoinAI.net/task/analyConfig"
	"github.com/EasyGolang/goTools/mTime"
)

var ResultBasePath = "/root/AItrade/CoinAI.net/task/analyConfig/最近8个月2"

func main() {
	// 初始化系统参数
	global.Start()

	// Step1("BTC-USDT")
	Step2("BTC-USDT")
}

func Step1(InstID string) {
	// 第一步： 暴力求值 （海量参数结果罗列）
	EmaPArr := []int{}
	StarNum := 70
	for i := 0; i < 520; i += 2 {
		StarNum = 60 + i
		EmaPArr = append(EmaPArr, StarNum)
	}
	EndTime := mTime.GetUnixInt64()
	StartTime := EndTime - (mTime.UnixTimeInt64.Day * 260)
	analyConfig.Violence(analyConfig.ViolenceOpt{
		StartTime: StartTime,
		EndTime:   EndTime,
		InstID:    InstID,
		EmaPArr:   EmaPArr,
		CAPArr:    []int{2, 3, 4, 5, 6},
		LevelArr:  []int{1},
		CAPMax:    []string{"0.5", "1", "1.5", "2", "2.5", "3"},
		ConfArr:   []okxInfo.TradeKdataOpt{},
		OutPutDir: ResultBasePath,
	})
}

func Step2(InstID string) {
	// 第二步骤：根据胜率和最高营收 进行筛选
	analyConfig.GetWinArr(analyConfig.GetWinArrOpt{
		InstID:     InstID,
		OutPutDir:  ResultBasePath,
		MoneyRight: "1700",
		WinRight:   "0.3",
	})
}
