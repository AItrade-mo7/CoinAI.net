package taskStart

import (
	"fmt"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/utils/taskPush"
	"CoinAI.net/task/testHunter"
	"github.com/EasyGolang/goTools/mCount"
	"github.com/EasyGolang/goTools/mFile"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/EasyGolang/goTools/mTime"
)

type BackOpt struct{}

func BackTest() {
	StartTime := mTime.GetUnix()

	// 新建回测
	backObj := testHunter.New(testHunter.TestOpt{
		StartTime: mTime.TimeParse(mTime.Lay_DD, "2021-01-01"),
		EndTime:   mTime.TimeParse(mTime.Lay_DD, "2022-01-01"),
		InstID:    "BTC-USDT",
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
			60, 62, 64, 66, 68, 70, 72, 74, 76, 77, 78, 80, 82, 84, 86, 88,
			160, 162, 164, 166, 168, 170, 172, 174, 176, 178, 180, 182, 184, 186, 188, 190,
			526, 528, 530, 532, 534, 536, 538, 540, 542, 544, 546, 548, 550, 552, 554,
		},
		CAPArr: []int{
			3,
			4,
			5,
		},
		LevelArr: []int{1},
	})

	TaskChan := make(chan string, len(configObj.GorMap)) // 记录线程任务完成

	BillingArr := []testHunter.BillingType{}

	// 建立一个线程要运行的任务
	NewGorTask := func(GorName string, confArr []testHunter.NewMockOpt) {
		global.Run.Println("开始执行Goroutine:", GorName)
		StartTime := mTime.GetUnix()

		for _, conf := range confArr {
			MockObj := backObj.NewMock(conf)
			Billing := MockObj.MockRun()
			BillingArr = append(BillingArr, Billing)
		}

		EndTime := mTime.GetUnix()
		DiffTime := mCount.Sub(EndTime, StartTime)
		DiffMin := mCount.Div(DiffTime, mTime.UnixTime.Minute)
		global.Run.Println("Goroutine:", GorName, "执行结束,共计耗时:", DiffMin, "分钟")
		TaskChan <- GorName
	}

	goRoName := []string{}
	for key, confArr := range configObj.GorMap {
		goRoName = append(goRoName, key)
		go NewGorTask(key, confArr)
	}

	taskPush.SysEmail(taskPush.SysEmailOpt{
		From:        config.SysName,
		To:          config.NoticeEmail,
		Subject:     "新建任务",
		Title:       mStr.Join("Cpu核心数(默认-1):", configObj.CpuNum, "任务总数:", configObj.TaskNum),
		Content:     "任务视图:<br />" + mJson.Format(configObj.GorMapView) + "线程数量" + mJson.Format(goRoName),
		Description: "回测开始通知",
	})

	// 终止任务
	taskEnd := []string{}
	for ok := range TaskChan {
		taskEnd = append(taskEnd, ok)
		if len(taskEnd) >= len(goRoName) {
			break
		}
	}

	EndTime := mTime.GetUnix()
	DiffTime := mCount.Sub(EndTime, StartTime)
	DiffMin := mCount.Div(DiffTime, mTime.UnixTime.Minute)

	BillingArr_Path := mStr.Join(config.Dir.JsonData, "/", "BillingArr.json")
	mFile.Write(BillingArr_Path, string(mJson.ToJson(BillingArr)))
	global.Run.Println("BillingArr: ", BillingArr_Path)

	taskPush.SysEmail(taskPush.SysEmailOpt{
		From:    config.SysName,
		To:      config.NoticeEmail,
		Subject: "任务结束",
		Title:   mStr.Join("共计耗时", DiffMin, "分钟"),
		Content: mStr.Join(
			"任务视图:<br />", mJson.Format(configObj.GorMapView), "<br />",
			"结果:", mJson.Format(BillingArr),
		),
		Description: "回测结束通知",
	})
}

type ConfOpt struct {
	MA_Period  int
	CAP_Period int
}
