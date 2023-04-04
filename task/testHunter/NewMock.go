package testHunter

import (
	"runtime"

	"CoinAI.net/server/global"
	"CoinAI.net/server/hunter"
	"github.com/EasyGolang/goTools/mCount"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/EasyGolang/goTools/mTime"
)

// 开仓信息记录
type PositionType struct {
	Dir         int    // 开仓方向
	OpenAvgPx   string // 开仓价格
	OpenTimeStr string // 开仓时间
	NowTimeStr  string
	NowC        string
	InstID      string // 下单币种
	UplRatio    string // 未实现收益率
}
type OrderType struct {
	Type    string // 平仓,Close  开空,Sell  开多,Buy
	AvgPx   string // 开仓价格
	InstID  string // 下单币种
	TimeStr string // 开仓时间
}

type RecordType struct {
	Value   string
	TimeStr string
}

// 收益结算
type BillingType struct {
	InstID           string     // 交易币种
	MockName         string     // 策略名称
	Days             int64      // 总天数
	StartTime        string     // 第一次持仓时间
	EndTime          string     // 结束时间
	NilNum           int        // 空仓次数
	SellNum          int        // 开空次数
	BuyNum           int        // 开多次数
	AllNum           int        // 总开仓次数
	Win              int        // 盈利次数
	WinRatio         string     // 总盈利比率
	Lose             int        // 亏损次数
	LoseRatio        string     // 总亏损比率
	MaxRatio         RecordType // 平仓后单笔最大盈利比率
	MinRatio         RecordType // 平仓后单笔最小盈利比率
	Charge           string     // 手续费率
	ChargeAll        string     // 总手续费
	InitMoney        string     // 初始金钱
	Money            string     // 账户当前余额
	MinMoney         RecordType // 平仓后历史最低余额
	MaxMoney         RecordType // 平仓后历史最高余额
	PositionMinRatio RecordType // 持仓过程中最低盈利比率
	PositionMaxRatio RecordType // 持仓过程中最高盈利比率
	Level            string     // 杠杆倍数
}

type NewMockOpt struct {
	MockName      string // 策略名字 MA_x_CAP_x
	InitMoney     string // 初始金钱  1000
	Level         string // 杠杆倍数  1
	Charge        string // 手续费  0.05
	TradeKdataOpt hunter.TradeKdataOpt
}

type MockObj struct {
	NowPosition   PositionType   // 当前持仓
	PositionArr   []PositionType // 当前持仓
	OrderArr      []OrderType    // 下单列表
	Billing       BillingType
	RunKdataList  []mOKX.TypeKd
	TradeKdataOpt hunter.TradeKdataOpt
}

func (_this *TestObj) NewMock(opt NewMockOpt) *MockObj {
	var obj MockObj

	obj.NowPosition = PositionType{}
	obj.PositionArr = []PositionType{}
	obj.OrderArr = []OrderType{}
	obj.RunKdataList = _this.KdataList
	// 开始处理参数
	obj.Billing = BillingType{}
	obj.Billing.MockName = opt.MockName
	obj.Billing.InitMoney = opt.InitMoney // 设定初始资金
	obj.Billing.Money = opt.InitMoney     // 设定当前账户资金
	obj.Billing.Level = opt.Level
	obj.Billing.Charge = opt.Charge
	obj.Billing.InstID = _this.KdataList[0].InstID
	obj.Billing.Days = (_this.EndTime - _this.StartTime) / mTime.UnixTimeInt64.Day
	obj.Billing.MinMoney.Value = opt.InitMoney
	obj.Billing.MaxMoney.Value = opt.InitMoney
	obj.TradeKdataOpt = opt.TradeKdataOpt

	global.Run.Println("新建Mock",
		mJson.Format(map[string]any{
			"参数组名称":   obj.Billing.MockName,
			"初始资金":    obj.Billing.InitMoney,
			"杠杆倍率":    obj.Billing.Level,
			"手续费率(%)": obj.Billing.Charge,
		}),
		mJson.Format(obj.TradeKdataOpt),
	)

	return &obj
}

type GetConfigOpt struct {
	EmaPArr []int
	CAPArr  []int
}

type GetConfigReturn struct {
	GorMap     map[string][]NewMockOpt
	GorMapView map[string][]string
	TaskNum    int
	CpuNum     int
}

func GetConfig(opt GetConfigOpt) GetConfigReturn {
	MockConfigArr := []NewMockOpt{}
	for _, emaP := range opt.EmaPArr {
		for _, cap := range opt.CAPArr {
			MockConfigArr = append(MockConfigArr,
				NewMockOpt{
					MockName:  "MA_" + mStr.ToStr(emaP) + "_CAP_" + mStr.ToStr(cap),
					InitMoney: "1000", // 初始资金
					Level:     "1",    // 杠杆倍数
					Charge:    "0.05", // 吃单标准手续费率 0.05%
					TradeKdataOpt: hunter.TradeKdataOpt{
						MA_Period:      emaP,
						RSI_Period:     18,
						RSI_EMA_Period: 14,
						CAP_Period:     cap,
					},
				},
			)
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
	GorMapNameArr := []string{}
	for i := 0; i < runtime.NumCPU(); i++ {
		GorName := mStr.Join("Gor_", i)
		GorMap[GorName] = []NewMockOpt{}
		GorMapNameArr = append(GorMapNameArr, GorName)
	}

	GorMapView := map[string][]string{}
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
		"当前任务总数量:", len(MockConfigArr),
		"当前CPU核心数量", CpuNumStr,
		"\n任务视图:\n", mJson.Format(GorMapView),
	)

	return GetConfigReturn{
		GorMap:     GorMap,
		GorMapView: GorMapView,
		TaskNum:    len(MockConfigArr),
		CpuNum:     CpuNum,
	}
}