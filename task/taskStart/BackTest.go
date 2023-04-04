package taskStart

import (
	"fmt"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
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
		EndTime:   mTime.TimeParse(mTime.Lay_DD, "2023-04-05"),
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

	// 新建回测参数 ( 按照核心数进行任务拆分 )
	configObj := testHunter.GetConfig(testHunter.GetConfigOpt{
		EmaPArr: []int{75, 77, 79, 169, 171, 173, 543, 545, 547},
		CAPArr:  []int{3, 4},
	})

	TaskChan := make(chan string, len(configObj.GorMap)) // 记录任务完成数

	// 建立一个线程要运行的任务

	NewGorTask := func(GorName string, confArr []testHunter.NewMockOpt) {
		global.Run.Println("开始执行Goroutine:", GorName)
		StartTime := mTime.GetUnix()
		for _, config := range confArr {
			MockObj := backObj.NewMock(config)
			MockObj.MockRun()
		}
		EndTime := mTime.GetUnix()
		DiffTime := mCount.Sub(EndTime, StartTime)
		DiffMin := mCount.Div(DiffTime, mTime.UnixTime.Minute)
		global.Run.Println("Goroutine:", GorName, "执行结束,共计耗时:", DiffMin, "分钟")
		TaskChan <- GorName
	}
	for key, confArr := range configObj.GorMap {
		go NewGorTask(key, confArr)
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
		if len(taskEnd) == len(configObj.GorMap) {
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
