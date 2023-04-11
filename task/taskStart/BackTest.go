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
	"github.com/EasyGolang/goTools/mPath"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/EasyGolang/goTools/mTime"
)

type BackOpt struct {
	InstID       string
	StartTime    int64
	EndTime      int64
	GetConfigOpt testHunter.GetConfigOpt
	OutPutDir    string // 输出目录
}

func BackTest(opt BackOpt) BackReturn {
	if !mPath.Exists(opt.OutPutDir) {
		err := fmt.Errorf("目录不存在 %+v", opt.OutPutDir)
		panic(err)
	}

	StartTime := mTime.GetUnix()

	// 新建数据
	backObj := testHunter.NewDataBase(testHunter.TestOpt{
		StartTime: opt.StartTime,
		EndTime:   opt.EndTime,
		InstID:    opt.InstID,
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
	configObj := testHunter.GetConfig(opt.GetConfigOpt)

	fmt.Println(opt.InstID, "任务总数:", len(configObj.ConfigArr))

	BillingArr := []testHunter.BillingType{} // 模拟交易的结果

	TaskChan := make(chan string, len(configObj.GorMap)) // 记录线程Chan

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
		TaskChan <- GorName // 线程完成
	}

	goRoName := []string{}
	for key, confArr := range configObj.GorMap {
		goRoName = append(goRoName, key)
		go NewGorTask(key, confArr)
	}

	taskPush.SysEmail(taskPush.SysEmailOpt{
		From:    config.SysName,
		To:      config.NoticeEmail,
		Subject: "新建任务" + opt.InstID,
		Title:   mStr.Join("并行任务数:", len(configObj.GorMap), "任务总数:", len(configObj.ConfigArr)),
		Content: mStr.Join(
			"<br />任务视图:<br />",
			mJson.Format(configObj.GorMapView),
			"<br />线程数量:<br />",
			mJson.Format(goRoName),
		),
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

	BillingArr_Path := mStr.Join(opt.OutPutDir, "/", opt.InstID, "-BillingArr.json")
	mFile.Write(BillingArr_Path, string(mJson.ToJson(BillingArr)))
	global.Run.Println("BillingArr: ", BillingArr_Path)

	taskPush.SysEmail(taskPush.SysEmailOpt{
		From:    config.SysName,
		To:      config.NoticeEmail,
		Subject: "任务结束" + opt.InstID,
		Title:   mStr.Join("任务总数:", len(configObj.ConfigArr), "共计耗时", DiffMin, "分钟"),
		Content: mStr.Join(
			"任务视图:<br />", mJson.Format(configObj.GorMapView), "<br />",
			"线程数量:<br />"+mJson.Format(goRoName), "<br />",
			"结果:<br />", BillingArr_Path, "<br />",
		),
		Description: "回测结束通知",
	})

	return BackReturn{
		InstID:      opt.InstID,
		StartTime:   opt.StartTime,
		EndTime:     opt.EndTime,
		BillingArr:  BillingArr,
		BillingPath: BillingArr_Path,
	}
}

type BackReturn struct {
	InstID         string
	StartTime      int64
	EndTime        int64
	BillingArr     []testHunter.BillingType
	BillingPath    string
	ResultBasePath string
}
