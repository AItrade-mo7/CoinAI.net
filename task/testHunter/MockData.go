package testHunter

import (
	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/hunter"
	"github.com/EasyGolang/goTools/mCount"
	"github.com/EasyGolang/goTools/mFile"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mStr"
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

var (
	NowPosition PositionType   // 当前持仓
	PositionArr []PositionType // 当前持仓
	OrderArr    []OrderType    // 下单列表

)

// 收益结算
type BillingType struct {
	NilNum    int    // 空仓次数
	SellNum   int    // 开空次数
	BuyNum    int    // 开多次数
	AllNum    int    // 总开仓次数
	Win       int    // 盈利次数
	WinRatio  string // 总盈利比率
	Lose      int    // 亏损次数
	LoseRatio string // 总亏损比率
	MaxRatio  string //  最大盈利比率
	MinRatio  string //  最小盈利比率
	Charge    string // 手续费率
	ChargeAll string // 总手续费
	MockName  string // 名字
	InitMoney string // 初始金钱
	Money     string // 账户当前余额
	Level     string // 杠杆倍数
}

var Billing BillingType

// 根据下单结果进行模拟持仓
func BillingFun(NowKdata hunter.TradeKdType) {
	// fmt.Println("下单总结一次",
	// 	NowKdata.TimeStr,
	// 	"持仓方向", NowPosition.Dir,
	// 	"收益率", NowPosition.UplRatio,
	// )

	if NowPosition.Dir == 0 {
		Billing.NilNum++ // 空仓计数
	}
	if NowPosition.Dir < 0 {
		Billing.SellNum++ // 开空 计数
	}
	if NowPosition.Dir > 0 {
		Billing.BuyNum++ // 开多 计数
	}

	if mCount.Le(NowPosition.UplRatio, "0") > 0 {
		Billing.Win++                                                         // 盈利次数计数
		Billing.WinRatio = mCount.Add(NowPosition.UplRatio, Billing.WinRatio) // 盈利比例相加
	}

	if mCount.Le(NowPosition.UplRatio, "0") < 0 {
		Billing.Lose++                                                          // 亏损次数计数
		Billing.LoseRatio = mCount.Add(NowPosition.UplRatio, Billing.LoseRatio) // 盈利比例相加
	}
	// 单次最大亏损和单次最大盈利
	if mCount.Le(NowPosition.UplRatio, Billing.MaxRatio) > 0 {
		Billing.MaxRatio = NowPosition.UplRatio
	}
	if mCount.Le(NowPosition.UplRatio, Billing.MinRatio) < 0 {
		Billing.MinRatio = NowPosition.UplRatio
	}

	Upl := mCount.Div(NowPosition.UplRatio, "100") // 格式化收益率
	ChargeUpl := mCount.Div(Billing.Charge, "100") // 格式化手续费率

	makeMoney := mCount.Mul(Billing.Money, Upl)          // 当前盈利的金钱
	Billing.Money = mCount.Add(Billing.Money, makeMoney) // 相加得出当账户总资金量

	nowCharge := mCount.Mul(Billing.Money, ChargeUpl)    // 当前产生的手续费
	Billing.Money = mCount.Sub(Billing.Money, nowCharge) // 减去手续费

	Billing.ChargeAll = mCount.Add(Billing.ChargeAll, nowCharge) // 记录一下手续费

	Billing.AllNum++ // 记录一下总交易次数
}

// 模拟数据流动并执行分析交易
func (_this *TestObj) MockData(MockOpt BillingType, TradeKdataOpt hunter.TradeKdataOpt) {
	// 收益结算
	Billing = BillingType{}
	Billing.MockName = MockOpt.MockName
	Billing.InitMoney = MockOpt.InitMoney // 设定初始资金
	Billing.Money = MockOpt.InitMoney     // 设定当前账户资金
	Billing.Level = MockOpt.Level
	Billing.Charge = MockOpt.Charge

	// 交易信息清空
	NowPosition = PositionType{}   // 清空持仓
	PositionArr = []PositionType{} // 清空持仓数组
	OrderArr = []OrderType{}       // 下单数组清空

	global.Run.Println("新建Mock数据",
		mJson.Format(map[string]any{
			"参数组名称":   Billing.MockName,
			"初始资金":    Billing.InitMoney,
			"杠杆倍率":    Billing.Level,
			"手续费率(%)": Billing.Charge,
		}),
		mJson.Format(TradeKdataOpt),
	)

	global.TradeLog.Println(" ============== 开始分析和交易 ============== ", Billing.MockName)

	// 清理 TradeKdataList
	hunter.TradeKdataList = []hunter.TradeKdType{}
	hunter.EMA_Arr = []string{}
	hunter.MA_Arr = []string{}
	hunter.RSI_Arr = []string{}

	FormatEnd := []mOKX.TypeKd{}
	for _, Kdata := range _this.KdataList {
		FormatEnd = append(FormatEnd, Kdata)

		if len(FormatEnd) < TradeKdataOpt.MA_Period {
			continue
		}

		TradeKdata := hunter.NewTradeKdata(FormatEnd, TradeKdataOpt)
		hunter.TradeKdataList = append(hunter.TradeKdataList, TradeKdata)

		// 开始执行分析交易
		Analy()

		if mCount.Le(NowPosition.UplRatio, "-45") < 0 {
			global.Log.Println("爆仓！", Billing.MockName, Kdata.TimeStr)
			break
		}
	}

	// 记录 整理好的数组
	TradeKdataList_Path := mStr.Join(config.Dir.JsonData, "/", Billing.MockName, "-TradeKdataList.json")
	mFile.Write(TradeKdataList_Path, string(mJson.ToJson(hunter.TradeKdataList)))
	global.Run.Println("TradeKdataList: ", TradeKdataList_Path)

	// 记录 持仓数组
	PositionArr_Path := mStr.Join(config.Dir.JsonData, "/", Billing.MockName, "-PositionArr.json")
	mFile.Write(PositionArr_Path, string(mJson.ToJson(PositionArr)))
	global.Run.Println("PositionArr: ", PositionArr_Path)

	// 记录 下单数组
	OrderArr_Path := mStr.Join(config.Dir.JsonData, "/", Billing.MockName, "-OrderArr.json")
	mFile.Write(OrderArr_Path, string(mJson.ToJson(OrderArr)))
	global.Run.Println("OrderArr: ", OrderArr_Path)

	global.TradeLog.Println(" 交易结果  ", mJson.Format(Billing))
}

func Analy() {
	NowKdata := hunter.TradeKdataList[len(hunter.TradeKdataList)-1]
	AnalyDir := 0 // 分析的方向，默认为 0 不开仓

	if mCount.Le(NowKdata.CAP_EMA, "0") > 0 { // 大于 0 则开多
		AnalyDir = 1
	}

	if mCount.Le(NowKdata.CAP_EMA, "0") < 0 { // 小于 0 则开空
		AnalyDir = -1
	}

	// 更新持仓状态
	if NowPosition.Dir != 0 { // 当前为持持仓状态，则计算收益率
		UplRatio := mCount.RoseCent(NowKdata.C, NowPosition.OpenAvgPx)
		if NowPosition.Dir < 0 { // 当前为持 空 仓状态 则翻转该值
			UplRatio = mCount.Sub("0", UplRatio)
		}
		NowPosition.UplRatio = mCount.Mul(UplRatio, Billing.Level) // 乘以杠杆倍数
	}
	NowPosition.NowTimeStr = NowKdata.TimeStr
	NowPosition.NowC = NowKdata.C
	PositionArr = append(PositionArr, NowPosition)

	// 当前持仓与 判断方向不符合时，执行一次下单操作
	if NowPosition.Dir != AnalyDir {
		OnOrder(AnalyDir, NowKdata)
	}
}

// 下单  参数：dir 下单方向 NowKdata : 当前市场行情
func OnOrder(dir int, NowKdata hunter.TradeKdType) {
	BillingFun(NowKdata) // 下单之前 计算一次收益

	if dir > 0 { // 开多
		// 下订单
		OrderArr = append(OrderArr, OrderType{
			Type:    "Buy",            // 下多单
			InstID:  NowKdata.InstID,  // 下单币种
			AvgPx:   NowKdata.C,       // 记录下单价格
			TimeStr: NowKdata.TimeStr, // 记录下单时间
		})
		// 更新持仓状态
		NowPosition = PositionType{
			Dir:         1,          // 持仓多方向
			OpenAvgPx:   NowKdata.C, // 持仓价格
			NowTimeStr:  NowKdata.TimeStr,
			UplRatio:    "0", // 当前收益率
			NowC:        NowKdata.C,
			OpenTimeStr: NowKdata.TimeStr, // 开仓时间
			InstID:      NowKdata.InstID,  // 开仓币种
		}
	}

	if dir < 0 { // 开空
		// 下订单
		OrderArr = append(OrderArr, OrderType{
			Type:    "Sell",           // 下空单
			InstID:  NowKdata.InstID,  // 下单币种
			AvgPx:   NowKdata.C,       // 记录下单价格
			TimeStr: NowKdata.TimeStr, // 记录下单时间
		})
		// 更新持仓状态
		NowPosition = PositionType{
			Dir:         -1,         // 持仓空方向
			OpenAvgPx:   NowKdata.C, // 持仓价格
			NowTimeStr:  NowKdata.TimeStr,
			UplRatio:    "0", // 当前收益率
			NowC:        NowKdata.C,
			OpenTimeStr: NowKdata.TimeStr, // 开仓时间
			InstID:      NowKdata.InstID,  // 开仓币种
		}
	}

	if dir == 0 { // 平仓
		// 下订单
		OrderArr = append(OrderArr, OrderType{
			Type:    "Close",          // 平仓
			InstID:  NowKdata.InstID,  // 下单币种
			AvgPx:   NowKdata.C,       // 记录下单价格
			TimeStr: NowKdata.TimeStr, // 记录下单时间
		})
		// 更新为空仓状态
		NowPosition = PositionType{
			Dir:         0,  // 持仓空方向
			OpenAvgPx:   "", // 持仓价格
			NowTimeStr:  NowKdata.TimeStr,
			UplRatio:    "0", // 当前收益率
			NowC:        NowKdata.C,
			OpenTimeStr: NowKdata.TimeStr, // 开仓时间
			InstID:      NowKdata.InstID,  // 开仓币种
		}
	}
}
