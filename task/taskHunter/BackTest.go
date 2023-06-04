package taskHunter

import (
	"fmt"
	"github.com/tomcraven/goga"
	fo "github.com/tomcraven/goga/function_optimizer"
	"time"

	//"sync"
	"sync/atomic"

	"go.uber.org/zap"
	"strconv"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/hunter/testHunter"
	"CoinAI.net/server/utils/taskPush"
	"github.com/EasyGolang/goTools/mCount"
	"github.com/EasyGolang/goTools/mFile"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mPath"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/EasyGolang/goTools/mTime"
	"github.com/lukechampine/randmap"
)

type BackOpt struct {
	InstID       string
	StartTime    int64
	EndTime      int64
	GetConfigOpt testHunter.GetConfigOpt
	OutPutDir    string // 输出目录
}
type TaskConf struct {
	backObj  *testHunter.TestObj
	mockConf testHunter.NewMockOpt
	opt      BackOpt
	counter  atomic.Int64

	candidate         map[goga.Genome]map[string]float64
	removedCandidates map[string]float64
	ensuredCandidates []goga.Genome
	beginCount        int
}

func tmToStr(t int64) string {
	return time.Unix(t/1000, 0).Format(mTime.Lay_ss)
}

func (tc *TaskConf) runTask(params []float64) float64 {
	// 新建数据
	tc.mockConf.TradeKdataOpt.EMA_Period = int(params[0])
	tc.mockConf.TradeKdataOpt.CAP_Period = int(params[1])
	tc.mockConf.TradeKdataOpt.CAP_Max = fmt.Sprintf("%.1f", params[2])
	tc.mockConf.TradeKdataOpt.CAP_Min = fmt.Sprintf("%.1f", params[3])
	tc.mockConf.TradeKdataOpt.MaxTradeLever = int(params[4])
	MockObj := tc.backObj.NewMock(tc.mockConf)
	Billing := MockObj.MockRun()
	money, _ := strconv.ParseFloat(Billing.ResultMoney, 64)
	tc.counter.Add(1)
	fmt.Printf("run task: %d [%v] done money: %f stop reason: %s start time: %s, end time: %s\n",
		tc.counter.Load(), params, money, Billing.StopReason, Billing.StartTime, Billing.EndTime)
	if Billing.StopReason != "norm" {
		st := mTime.TimeParse(mTime.Lay_ss, Billing.StartTime)
		if st == 0 {
			st = tc.backObj.StartTime
		}
		et := mTime.TimeParse(mTime.Lay_ss, Billing.EndTime)
		if et == 0 {
			et = tc.backObj.EndTime
		}
		money = float64(et-st) / float64(tc.backObj.EndTime-tc.backObj.StartTime) * 500
	}
	return money
}

func (tc *TaskConf) RefreshBackTestTimeRange(startTime, endTime int64, instID string) {
	// 新建数据
	tc.opt.InstID = instID
	tc.opt.StartTime = startTime
	tc.opt.EndTime = endTime
	backObj := testHunter.NewDataBase(testHunter.TestOpt{
		StartTime: tc.opt.StartTime,
		EndTime:   tc.opt.EndTime,
		InstID:    tc.opt.InstID,
	})
	tc.backObj = backObj
	err := backObj.StuffDBKdata()
	if err != nil {
		panic(fmt.Errorf("出错: %+v", err))
	}
	err = backObj.CheckKdataList() // 检查数据是否出错
	if err != nil {
		panic(fmt.Errorf("出错: %+v", err))
	}
	fmt.Printf("RefreshBackTestTimeRange to start time: %s end time: %s\n", tmToStr(startTime), tmToStr(endTime))
}

func (tc *TaskConf) OnEnd(gs []goga.Genome) {
	for _, g := range gs {
		if _, ok := tc.removedCandidates[g.Key()]; ok {
			continue
		}
		if g.GetFitness() > 1100 {
			if m, ok := tc.candidate[g]; ok {
				m[tmToStr(tc.opt.StartTime)+" "+tmToStr(tc.opt.EndTime)] = g.GetFitness()
			} else {
				tc.candidate[g] = make(map[string]float64)
				tc.candidate[g][tmToStr(tc.opt.StartTime)+" "+tmToStr(tc.opt.EndTime)] = g.GetFitness()
			}
		} else if g.GetFitness() < 1000 {
			delete(tc.candidate, g)
			tc.removedCandidates[g.Key()] = g.GetFitness()
		}
	}
}

func (tc *TaskConf) OnBegin() []goga.Genome {
	l := make([]goga.Genome, len(tc.ensuredCandidates))
	i := 0
	for _, e := range tc.ensuredCandidates {
		l[i] = e
		i += 1
	}
	k := goga.NewGenome(goga.Bitset{})
	m := make(map[string]float64)
	it := randmap.FastIter(tc.candidate, &k, &m)
	for it.Next() {
		l = append(l, k)
		if len(l) > tc.beginCount {
			break
		}
	}
	return l
}

func getGenome(vs []float64) goga.Genome {
	b := &goga.Bitset{}
	b.Create(len(vs) * 8)
	for i, v := range vs {
		byteArr := goga.Float64ToByte(v)
		for idx := 0; idx < 8; idx++ {
			b.Set(i*8+idx, int(byteArr[idx]))
		}
	}
	return goga.NewGenome(*b)
}

func BackTestWithGeneticAlgo(startTime, endTime int64, instID, outPutDir string) {
	if !mPath.Exists(outPutDir) {
		err := fmt.Errorf("目录不存在 %+v", outPutDir)
		panic(err)
	}
	var timeRange = (endTime - startTime) / 4
	var timeMove = timeRange / 3

	tc := &TaskConf{
		mockConf: testHunter.NewMockOpt{
			InitMoney: "1000", // 初始金钱  1000
			ChargeUpl: "0.05", // 手续费率  0.05
		},
		opt: BackOpt{
			OutPutDir: outPutDir,
			InstID:    instID,
		},
		candidate:         make(map[goga.Genome]map[string]float64),
		removedCandidates: make(map[string]float64),
		ensuredCandidates: []goga.Genome{getGenome([]float64{294, 6, 3, -0.5, 4})},
		beginCount:        10,
	}

	initStartTime := endTime - timeRange
	tc.RefreshBackTestTimeRange(initStartTime, initStartTime+timeRange, instID)
	paramSize := 5
	requirement := goga.Float64Requirement{
		Specific: map[int]struct {
			Precision float64
			MaxValue  float64
			MinValue  float64
		}{
			0: {
				Precision: 1,
				MaxValue:  30,
				MinValue:  500,
			},
			1: {
				Precision: 1,
				MinValue:  2,
				MaxValue:  10,
			},
			2: {
				Precision: 0.5,
				MinValue:  0.5,
				MaxValue:  3,
			},
			3: {
				Precision: 0.5,
				MinValue:  -3,
				MaxValue:  0.5,
			},
			4: {
				Precision: 1,
				MinValue:  1,
				MaxValue:  5,
			},
		},
	}
	algo := fo.NewFuncAlgo(fo.Function(tc.runTask), fo.TransFunc(func(f float64) float64 {
		if f < 0 {
			return 0
		}
		return f
	}), fo.Requirement(&requirement), fo.ParamSize(paramSize),
		fo.StableMinIter(5),
		fo.PopulationSize(100),
		fo.MaterExtraRatio(4),
		fo.LRUSize(100),
		fo.OnStable(func() {
			if initStartTime-timeMove < startTime {
				timeRange = timeRange + 30*24*3600*1000
				if timeRange > endTime-startTime {
					timeRange = endTime - startTime
				}
				timeMove = timeRange / 3
				initStartTime = endTime - timeRange
			} else {
				initStartTime -= timeMove
			}
			tc.RefreshBackTestTimeRange(initStartTime, initStartTime+timeRange, instID)
		}),
		fo.OnBegin(tc.OnBegin),
		fo.OnEnd(func(genomes []goga.Genome) {
			if initStartTime-timeMove < startTime {
				initStartTime = endTime - timeRange
			} else {
				initStartTime -= timeMove
			}
			tc.RefreshBackTestTimeRange(initStartTime, initStartTime+timeRange, instID)
			tc.OnEnd(genomes)
		}))

	algo.Simulate()
	//wg := sync.WaitGroup{}
	//cnt := 5
	//wg.Add(cnt)
	//for i := 0; i < cnt; i++ {
	//	go func() {
	//		tc.runTask([]float64{43, 2, 3, -2, 3})
	//		tc.runTask([]float64{200, 5, 3, -2, 3})
	//		wg.Done()
	//	}()
	//}
	//wg.Wait()

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

	BruteForce := func() ([]testHunter.BillingType, string) {
		// 构建参数
		// 新建回测参数 ( 按照核心数进行任务拆分 )
		configObj := testHunter.GetConfig(opt.GetConfigOpt)

		fmt.Println(opt.InstID, "任务总数:", len(configObj.ConfigArr))

		BillingArr := []testHunter.BillingType{}             // 模拟交易的结果
		TaskChan := make(chan string, len(configObj.GorMap)) // 记录线程Chan
		// 建立一个线程要运行的任务
		NewGorTask := func(GorName string, confArr []testHunter.NewMockOpt) {
			global.Run.Info("开始执行Goroutine: " + GorName)
			StartTime := mTime.GetUnix()
			for _, conf := range confArr {
				MockObj := backObj.NewMock(conf)
				Billing := MockObj.MockRun()
				BillingArr = append(BillingArr, Billing)
			}

			EndTime := mTime.GetUnix()
			DiffTime := mCount.Sub(EndTime, StartTime)
			DiffMin := mCount.Div(DiffTime, mTime.UnixTime.Minute)
			global.Run.Info("Goroutine:", zap.String(GorName, "执行结束,共计耗时:"), zap.String(DiffMin, "分钟"))
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
		global.Run.Info("BillingArr: " + BillingArr_Path)
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
		return BillingArr, BillingArr_Path
	}
	BillingArr, BillingPath := BruteForce()
	return BackReturn{
		InstID:      opt.InstID,
		StartTime:   opt.StartTime,
		EndTime:     opt.EndTime,
		BillingArr:  BillingArr,
		BillingPath: BillingPath,
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
