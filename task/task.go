package main

import (
	_ "embed"
	"fmt"
	"os"

	"CoinAI.net/server/global"
	"CoinAI.net/server/hunter/testHunter"
	"CoinAI.net/task/taskHunter"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/EasyGolang/goTools/mTime"
	jsoniter "github.com/json-iterator/go"
)

var ResultBasePath = "/Users/zhangxiaofei.sf/GolandProjects/CoinAI.net/task/2021_2023"

func main() {
	// 初始化系统参数
	global.Start()

	Step1("BTC-USDT")
	// Step2("BTC-USDT")
	// Step3("BTC-USDT")
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
		//fmt.Println(StarNum)
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
	taskHunter.BackTestParam(StartTime, EndTime, InstID, mStr.Join(ResultBasePath), []float64{43, 2, 3, -2, 3})
	//taskHunter.BackTestWithGeneticAlgo(
	//	StartTime, EndTime, InstID, mStr.Join(ResultBasePath))
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
	EmaFilePath := mStr.Join(ResultBasePath, "/", InstID, "-EmaArr.json")

	// confArr := taskHunter.GetWinConfig(taskHunter.GetWinConfigOpt{
	// 	OutPutDir: ResultBasePath,
	// 	InstID:    InstID,
	// })

	// // 提取 EMA 的值
	// EmaArr := []int{}
	// for _, conf := range confArr {
	// 	EmaArr = append(EmaArr, conf.EMA_Period)
	// }
	// mFile.Write(EmaFilePath, mJson.ToStr(EmaArr))

	// 新一轮求解，柔和两个时间段的值进行运算
	StartTime := mTime.TimeParse(mTime.Lay_DD, "2021-01-01")
	EndTime := mTime.TimeParse(mTime.Lay_DD, "2023-01-01")

	file, err := os.ReadFile(EmaFilePath)
	if err != nil {
		err := fmt.Errorf("读取文件出错 %+v", err)
		panic(err)
	}
	var EmaArr []int // 数据来源
	jsoniter.Unmarshal(file, &EmaArr)

	// StartTime := EndTime - (mTime.UnixTimeInt64.Day * 260)
	taskHunter.BackTest(taskHunter.BackOpt{
		StartTime: StartTime,
		EndTime:   EndTime,
		InstID:    InstID,
		OutPutDir: mStr.Join(ResultBasePath),
		GetConfigOpt: testHunter.GetConfigOpt{
			EmaPArr:  EmaArr,                                       // Ema 步长
			CAPArr:   []int{2, 3, 4, 5, 6, 7},                      // CAP 步长
			CAPMax:   []string{"0.5", "1", "1.5", "2", "2.5", "3"}, // CAPMax 步长
			CAPMin:   []string{"-0.5", "-1", "-1.5", "-2", "-2.5", "-3"},
			LevelArr: []int{1},
			// ConfArr: []dbType.TradeKdataOpt{
			// 	{
			// 		EMA_Period:    342,
			// 		CAP_Period:    7,
			// 		CAP_Max:       "2.5",
			// 		CAP_Min:       "-2.5",
			// 		MaxTradeLever: 5,
			// 	},
			// },
		},
	})
}

func Step4(InstID string) {
	// 第四步： 根据第三步的结果进行筛选 (胜率和最终营收) 得出参数结果
	taskHunter.GetWinArr(taskHunter.GetWinArrOpt{
		InstID:     InstID,
		OutPutDir:  mStr.Join(ResultBasePath),
		MoneyRight: "1000",
		// WinRight:   "0.35",
		Sort: "Win", //  Win
	})
}
