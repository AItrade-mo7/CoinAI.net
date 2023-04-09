package testHunter

import (
	"os"
	"runtime"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mCount"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mPath"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/EasyGolang/goTools/mTime"
)

type OrderType struct {
	Type    string // 平仓,Close  开空,Sell  开多,Buy
	AvgPx   string // 开仓价格
	InstID  string // 下单币种
	TimeStr string // 开仓时间
}

// 收益结算
type BillingType struct {
	MockName         string
	InstID           string
	AllDay           int64                  // 总天数 | 结束时计算
	StartTime        string                 // 第一次持仓时间 数组第一个 | 结束时计算
	EndTime          string                 // 结束时间 数组组后一个
	NilNum           int                    // 空仓次数 平仓后未开仓 NowDir = 0 | 结束时计算
	SellNum          int                    // 开空次数 平空次数 NowDir = -1 | 结束时计算
	BuyNum           int                    // 开多次数 平多次数 NowDir = 1 | 结束时计算
	AllNum           int                    // 总开仓次数 总的平仓次数 数组长度 | 结束时计算
	WinNum           int                    // 盈利次数 NowUplRatio > 0 的次数
	LoseNum          int                    // 亏损次数 同 盈利次数
	WinRatio         string                 // 胜率 盈利次数/(平空次数+平多次数)
	PLratio          string                 // 盈亏比
	WinUplRatioAdd   string                 // 总盈利比率 NowUplRatio > 0 的总和
	WinMoneyAdd      string                 // 盈利总金额 1000 块钱 从头计算一次 盈利部分相加
	LoseUplRatioAdd  string                 // 总亏损比率 同总的盈利比率
	LoseMoneyAdd     string                 // 亏损总金额 同上
	MaxRatio         okxInfo.RecordNodeType // 平仓后单笔最大盈利比率   平仓后的记录
	MinRatio         okxInfo.RecordNodeType // 平仓后单笔最小盈利比率
	ChargeAll        string                 // 总手续费 同上
	MinMoney         okxInfo.RecordNodeType // 平仓后历史最低余额  遍历一次就知道
	MaxMoney         okxInfo.RecordNodeType // 平仓后历史最高余额  遍历一次就知道
	PositionMinRatio okxInfo.RecordNodeType // 持仓过程中最低盈利比率  // 持仓过程中才知道 结合K线才能得出
	PositionMaxRatio okxInfo.RecordNodeType // 持仓过程中最高盈利比率 // 持仓过程中才知道
	InitMoney        string                 // 初始金钱
	ResultMoney      string                 // 最终金钱
	Level            string                 // 杠杆倍率
}

type NewMockOpt struct {
	MockName      string // 策略名字 MA_x_CAP_x
	InitMoney     string // 初始金钱  1000
	ChargeUpl     string // 手续费率  0.05
	TradeKdataOpt okxInfo.TradeKdataOpt
}

type MockObj struct {
	HunterName         string
	NowVirtualPosition okxInfo.VirtualPositionType   // 当前持仓
	PositionArr        []okxInfo.VirtualPositionType // 当前持仓列表
	OrderArr           []okxInfo.VirtualPositionType // 平仓列表
	Billing            BillingType                   // 交易结果
	RunKdataList       []mOKX.TypeKd                 // 原始的 Kdata 数据
	TradeKdataList     []okxInfo.TradeKdType         // 计算好各种指标之后的K线
	TradeKdataOpt      okxInfo.TradeKdataOpt         // 交易指标
	OutPutDirectory    string                        // 数据读写的目录
}

/*
新建回测

接受参数 NewMockOpt
产出： 收益结果
*/
func (_this *TestObj) NewMock(opt NewMockOpt) *MockObj {
	var obj MockObj

	obj.HunterName = opt.MockName

	obj.NowVirtualPosition = okxInfo.VirtualPositionType{}
	obj.NowVirtualPosition.InitMoney = opt.InitMoney
	obj.NowVirtualPosition.Money = opt.InitMoney
	obj.NowVirtualPosition.ChargeUpl = opt.ChargeUpl

	obj.PositionArr = []okxInfo.VirtualPositionType{}
	obj.OrderArr = []okxInfo.VirtualPositionType{}
	obj.RunKdataList = _this.KdataList
	obj.TradeKdataList = []okxInfo.TradeKdType{}
	obj.OutPutDirectory = mStr.Join(config.Dir.JsonData, "/", opt.MockName)
	// 默认目录在 jsonData 下
	isOutPutDirectoryPath := mPath.Exists(obj.OutPutDirectory)
	if !isOutPutDirectoryPath {
		// 不存在则创建 logs 目录
		os.MkdirAll(obj.OutPutDirectory, 0o777)
	}

	// 汇总结果的初始化
	obj.Billing = BillingType{}
	obj.Billing.MockName = opt.MockName
	obj.Billing.InstID = _this.KdataList[len(_this.KdataList)-1].InstID
	obj.Billing.AllDay = (_this.EndTime - _this.StartTime) / mTime.UnixTimeInt64.Day
	obj.Billing.MinMoney.Value = opt.InitMoney
	obj.Billing.MaxMoney.Value = opt.InitMoney
	obj.Billing.InitMoney = opt.InitMoney
	obj.Billing.ResultMoney = opt.InitMoney
	obj.Billing.Level = mStr.ToStr(opt.TradeKdataOpt.MaxTradeLever)
	// 设置交易指标
	obj.TradeKdataOpt = opt.TradeKdataOpt

	return &obj
}

type GetConfigOpt struct {
	EmaPArr  []int                   // Ema 步长
	CAPArr   []int                   // CAP 步长
	CAPMax   []string                // CAPMax 步长
	LevelArr []int                   // 杠杆倍数
	ConfArr  []okxInfo.TradeKdataOpt // 成型的参数数组
}

type GetConfigReturn struct {
	GorMap        map[string][]NewMockOpt
	GorMapView    map[string][]string
	ConfigArr     []NewMockOpt
	GorMapNameArr []string
}

func GetConfig(opt GetConfigOpt) GetConfigReturn {
	MockConfigArr := []NewMockOpt{}

	ChargeUpl := "0.05" //  https://www.okx.com/cn/fees
	InitMoney := "1000"

	AppendConfig := func(conf okxInfo.TradeKdataOpt) {
		MockConfigArr = append(MockConfigArr,
			NewMockOpt{
				MockName:  mStr.Join("EMA_", conf.EMA_Period, "_CAP_", conf.CAP_Period, "_CAPMax_", conf.CAP_Max, "_level_", conf.MaxTradeLever),
				InitMoney: InitMoney, // 初始资金
				ChargeUpl: ChargeUpl, // 吃单标准手续费率 0.05%
				TradeKdataOpt: okxInfo.TradeKdataOpt{
					EMA_Period:    conf.EMA_Period,
					CAP_Period:    conf.CAP_Period,
					MaxTradeLever: conf.MaxTradeLever,
					CAP_Max:       conf.CAP_Max,
				},
			},
		)
	}

	for _, conf := range opt.ConfArr {
		AppendConfig(conf)
	}

	for _, emaP := range opt.EmaPArr {
		for _, cap := range opt.CAPArr {
			for _, level := range opt.LevelArr {
				for _, capMax := range opt.CAPMax {
					conf := okxInfo.TradeKdataOpt{
						EMA_Period:    emaP,
						CAP_Period:    cap,
						MaxTradeLever: level,
						CAP_Max:       capMax,
					}
					AppendConfig(conf)
				}
			}
		}
	}

	// 根据 cpu 核心数计算每个 Goroutine 的最大任务数
	CpuNum := runtime.NumCPU()
	CpuNumStr := mStr.ToStr(CpuNum)
	taskNumStr := mStr.ToStr(len(MockConfigArr))
	MaxNumStr := mCount.Div(taskNumStr, CpuNumStr)
	MaxNumInt := mCount.ToInt(MaxNumStr)
	decNum := mCount.GetDecimal(MaxNumStr)
	if decNum > 0 {
		MaxNumInt = MaxNumInt + 1
	}

	GorMap := map[string][]NewMockOpt{}
	GorMapView := map[string][]string{}
	GorMapNameArr := []string{}
	for i := 0; i < CpuNum; i++ { // 按照Cpu核心数创建线程
		GorName := mStr.Join("Gor_", i)
		GorMap[GorName] = []NewMockOpt{}
		GorMapNameArr = append(GorMapNameArr, GorName)
	}

	NowNameIdx := 0
	for _, config := range MockConfigArr {
		gorName := GorMapNameArr[NowNameIdx]
		GorMap[gorName] = append(GorMap[gorName], config)
		GorMapView[gorName] = append(GorMapView[gorName], config.MockName)
		if len(GorMap[gorName]) >= MaxNumInt {
			NowNameIdx++
		}
	}

	global.Run.Println("新建参数集合:",
		"\n 任务视图: \n", mJson.Format(GorMapView),
		"\n 线程数量: \n", mJson.Format(GorMapNameArr),
	)

	return GetConfigReturn{
		GorMap:        GorMap,
		GorMapView:    GorMapView,
		ConfigArr:     MockConfigArr,
		GorMapNameArr: GorMapNameArr,
	}
}
