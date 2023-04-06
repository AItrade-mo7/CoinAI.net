package taskStart

import (
	"fmt"

	"CoinAI.net/server/global/config"
	"CoinAI.net/server/okxInfo"
	"CoinAI.net/server/utils/taskPush"
	"CoinAI.net/task/testHunter"
	"github.com/EasyGolang/goTools/mCount"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/EasyGolang/goTools/mTime"
)

func BackTest() {
	StartTime := mTime.GetUnix()

	// 新建回测
	backObj := testHunter.New(testHunter.TestOpt{
		StartTime: mTime.TimeParse(mTime.Lay_DD, "2020-01-01"),
		EndTime:   mTime.TimeParse(mTime.Lay_DD, "2023-05-01"),
		InstID:    "ETH-USDT",
	})
	err := backObj.StuffDBKdata()
	if err != nil {
		panic(fmt.Errorf("出错: %+v", err))
	}
	err = backObj.CheckKdataList() // 检查数据是否出错
	if err != nil {
		panic(fmt.Errorf("出错: %+v", err))
	}

	// 构建参数
	// 新建回测参数 ( 按照核心数进行任务拆分 )
	configObj := testHunter.GetConfig(testHunter.GetConfigOpt{
		EmaPArr: []int{
			70, 72, 74, 76, 77, 78, 80, 82, 84,
			164, 166, 168, 170, 171, 172, 174, 176, 178,
			538, 540, 542, 544, 545, 546, 548, 550, 552,
		},
		CAPArr: []int{3, 4},
	})
	/* 	ConfigArr := GetConfigArr([]ConfOpt{
	   		{
	   			MA_Period:  77,
	   			CAP_Period: 3,
	   		},
	   		{
	   			MA_Period:  171,
	   			CAP_Period: 4,
	   		},
	   		{
	   			MA_Period:  545,
	   			CAP_Period: 3,
	   		},
	   	})
	   	configObj := testHunter.GetConfigReturn{
	   		ConfigArr: ConfigArr,
	   		TaskNum:   len(ConfigArr),
	   	}
	*/

	// 构建参数完毕

	TaskChan := make(chan string, len(configObj.GorMap)) // 记录任务完成数

	// 建立一个线程要运行的任务

	// NewGorTask := func(GorName string, confArr testHunter.NewMockOpt) {
	// 	global.Run.Println("开始执行Goroutine:", GorName)
	// 	StartTime := mTime.GetUnix()

	// 	MockObj := backObj.NewMock(confArr)
	// 	MockObj.MockRun()

	// 	EndTime := mTime.GetUnix()
	// 	DiffTime := mCount.Sub(EndTime, StartTime)
	// 	DiffMin := mCount.Div(DiffTime, mTime.UnixTime.Minute)
	// 	global.Run.Println("Goroutine:", GorName, "执行结束,共计耗时:", DiffMin, "分钟")
	// 	TaskChan <- GorName
	// }

	goRNum := 0
	for _, confArr := range configObj.ConfigArr {
		goRNum++

		fmt.Println(confArr.MockName)
		// go NewGorTask(confArr.MockName, confArr)
	}

	taskPush.SysEmail(taskPush.SysEmailOpt{
		From:        config.SysName,
		To:          config.NoticeEmail,
		Subject:     "新建任务",
		Title:       mStr.Join("Cpu核心数:", configObj.CpuNum, "任务总数:", configObj.TaskNum),
		Content:     "任务视图:<br />" + mJson.Format(configObj.GorMapView),
		Description: "回测开始通知",
	})

	// 终止任务
	taskEnd := []string{}
	for ok := range TaskChan {
		taskEnd = append(taskEnd, ok)
		if len(taskEnd) == goRNum {
			break
		}
	}

	EndTime := mTime.GetUnix()
	DiffTime := mCount.Sub(EndTime, StartTime)
	DiffMin := mCount.Div(DiffTime, mTime.UnixTime.Minute)

	taskPush.SysEmail(taskPush.SysEmailOpt{
		From:        config.SysName,
		To:          config.NoticeEmail,
		Subject:     "任务结束",
		Title:       mStr.Join("共计耗时", DiffMin, "分钟"),
		Content:     "任务视图:<br />" + mJson.Format(configObj.GorMapView),
		Description: "回测结束通知",
	})
}

type ConfOpt struct {
	MA_Period  int
	CAP_Period int
}

func GetConfigArr(confArr []ConfOpt) []testHunter.NewMockOpt {
	ConfigArr := []testHunter.NewMockOpt{}
	for _, conf := range confArr {
		emaP := conf.MA_Period
		cap := conf.CAP_Period
		ConfigArr = append(ConfigArr, testHunter.NewMockOpt{
			MockName:  mStr.Join("MA_", mStr.ToStr(emaP), "_CAP_", mStr.ToStr(cap)),
			InitMoney: "1000", // 初始资金
			Level:     "1",    // 杠杆倍数
			Charge:    "0.05", // 吃单标准手续费率 0.05%
			TradeKdataOpt: okxInfo.TradeKdataOpt{
				EMA_Period: emaP,
				CAP_Period: cap,
			},
		})
	}

	return ConfigArr
}
