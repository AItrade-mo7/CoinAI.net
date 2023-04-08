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

func BackTest(instID string) {
	StartTime := mTime.GetUnix()

	// 新建回测
	backObj := testHunter.New(testHunter.TestOpt{
		StartTime: mTime.TimeParse(mTime.Lay_DD, "2021-01-01"),
		EndTime:   mTime.TimeParse(mTime.Lay_DD, "2022-01-01"),
		InstID:    instID,
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
			20, 22, 24, 26, 28, 30, 32, 34, 36, 38, 40, 42, 44, 46, 48, 50, 52, 54, 56, 58,
			60, 62, 64, 66, 68, 70, 72, 74, 76, 78, 80, 82, 84, 86, 88, 90, 92, 94, 96, 98,
			102, 104, 106, 108, 110, 112, 114, 116, 118, 120, 122, 124, 126, 128, 130, 132,
			134, 136, 138, 140, 142, 144, 146, 148, 150, 152, 154, 156, 158,
			160, 162, 164, 166, 168, 170, 172, 174, 176, 178, 180, 182, 184, 186, 188, 190,
			192, 194, 196, 198, 200, 202, 204, 206, 208, 210, 212, 214, 216, 218, 220, 222,
			224, 226, 228, 230, 232, 234, 236, 238, 240, 242, 244, 246, 248, 250, 252, 254,
			256, 258, 260, 262, 264, 266, 268, 270, 272, 274, 276, 278, 280, 282, 284,
			286, 288, 290, 292, 294, 296, 298, 300, 302, 304, 306, 308, 310, 312, 314,
			316, 318, 320, 322, 324, 326, 328, 330, 332, 334, 336, 338, 340, 342, 344,
			346, 348, 350, 352, 354, 356, 358, 360, 362, 364, 366, 368, 370, 372, 374,
			376, 378, 380, 382, 384, 386, 388, 390, 392, 394, 396, 398, 400, 402, 404,
			406, 408, 410, 412, 414, 416, 418, 420, 422, 424, 426, 428, 430, 432, 434,
			436, 438, 440, 442, 444, 446, 448, 450, 452, 454, 456, 458, 460, 462, 464, 466,
			468, 470, 472, 474, 476, 478, 480, 482, 484, 486, 488, 490, 492, 494, 496,
			498, 500, 502, 504, 506, 508, 510, 512, 514, 516, 518, 520, 522, 524, 526, 528, 530,
			532, 534, 536, 538, 540, 542, 544, 546, 548, 550, 552, 554, 556, 558, 560, 562,
			564, 566, 568, 570, 572, 574, 576, 578, 580, 582, 584, 586, 588, 590, 592,
		},
		CAPArr: []int{
			2,
			3,
			4,
			5,
			6,
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
		Subject:     "新建任务" + instID,
		Title:       mStr.Join("Cpu核心数(默认-1):", configObj.CpuNum, "任务总数:", configObj.TaskNum),
		Content:     "任务视图:<br />" + mJson.Format(configObj.GorMapView) + "线程数量:<br />" + mJson.Format(goRoName),
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

	BillingArr_Path := mStr.Join(config.Dir.JsonData, "/", instID, "-BillingArr.json")
	mFile.Write(BillingArr_Path, string(mJson.ToJson(BillingArr)))
	global.Run.Println("BillingArr: ", BillingArr_Path)

	taskPush.SysEmail(taskPush.SysEmailOpt{
		From:    config.SysName,
		To:      config.NoticeEmail,
		Subject: "任务结束",
		Title:   mStr.Join("共计耗时", DiffMin, "分钟"),
		Content: mStr.Join(
			"任务视图:<br />", mJson.Format(configObj.GorMapView), "<br />",
			"结果:<br />", mJson.Format(BillingArr), "<br />",
		),
		Description: "回测结束通知",
	})
}

type ConfOpt struct {
	MA_Period  int
	CAP_Period int
}
