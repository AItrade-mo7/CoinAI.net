package main

import (
	_ "embed"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/dbType"
	"CoinAI.net/task/analyConfig"
	"CoinAI.net/task/taskStart"
	"CoinAI.net/task/testHunter"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/EasyGolang/goTools/mTime"
)

var ResultBasePath = "/root/AItrade/CoinAI.net/task/analyConfig/最近8个月2"

func main() {
	// 初始化系统参数
	global.Start()

	// Step1("BTC-USDT")
	// Step2("BTC-USDT")
	// Step3("BTC-USDT")
	Step4("BTC-USDT")
}

func Step1(InstID string) {
	// 第一步： 暴力求值 （海量参数结果罗列） 需要几个小时甚至好几天
	EmaPArr := []int{}
	StarNum := 70
	for i := 0; i < 520; i += 2 {
		StarNum = 60 + i
		EmaPArr = append(EmaPArr, StarNum)
	}
	EndTime := mTime.GetUnixInt64()
	StartTime := EndTime - (mTime.UnixTimeInt64.Day * 260)
	taskStart.BackTest(taskStart.BackOpt{
		StartTime: StartTime,
		EndTime:   EndTime,
		InstID:    InstID,
		OutPutDir: ResultBasePath,
		GetConfigOpt: testHunter.GetConfigOpt{
			EmaPArr:  EmaPArr,
			CAPArr:   []int{2, 3, 4, 5, 6},
			LevelArr: []int{1},
			CAPMax:   []string{"0.5", "1", "1.5", "2", "2.5", "3"},
			ConfArr:  []dbType.TradeKdataOpt{},
		},
	})
}

func Step2(InstID string) {
	// 第二步骤：根据胜率和最终营收 筛选
	analyConfig.GetWinArr(analyConfig.GetWinArrOpt{
		InstID:     InstID,
		OutPutDir:  ResultBasePath,
		MoneyRight: "1700",
		WinRight:   "0.3",
	})
}

func Step3(InstID string) {
	// 第三步：提取第二步的配置，加上杠杆得出新的参数组合 大概 几百个 然后 换个新的时间进行测试
	confArr := analyConfig.GetWinConfig(analyConfig.GetWinConfigOpt{
		OutPutDir: ResultBasePath,
		InstID:    InstID,
	})

	LevelArr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	NewConfigArr := []dbType.TradeKdataOpt{}
	for _, leave := range LevelArr {
		for _, conf := range confArr {
			conf.MaxTradeLever = leave
			NewConfigArr = append(NewConfigArr, conf)
		}
	}
	// 新一轮求解，计算最优杠杆倍率 用  2022 年 8 月 的 380 天前进行回测
	EndTime := mTime.TimeParse(mTime.Lay_DD, "2022-08-01")
	StartTime := EndTime - (mTime.UnixTimeInt64.Day * 380)
	taskStart.BackTest(taskStart.BackOpt{
		StartTime: StartTime,
		EndTime:   EndTime,
		InstID:    InstID,
		OutPutDir: mStr.Join(ResultBasePath, "/三步最终结果"),
		GetConfigOpt: testHunter.GetConfigOpt{
			ConfArr: NewConfigArr,
		},
	})
}

func Step4(InstID string) {
	// 第四步： 根据第三步的结果进行筛选 (胜率和最终营收)
	analyConfig.GetWinArr(analyConfig.GetWinArrOpt{
		InstID:    InstID,
		OutPutDir: mStr.Join(ResultBasePath),
		// MoneyRight: "1000",
		// WinRight:   "0.3",
	})
}
