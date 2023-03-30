package testHunter

import (
	"fmt"

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
	Lose      int    // 亏损次数
	WinRatio  string // 盈利比率
	LoseRatio string // 亏损比率
	MaxWin    string //  最大盈利
	MaxLose   string //  最大亏损
	Charge    string // 手续费率
	ChargeAll string // 总手续费
	MockName  string // 名字
	InitMoney string // 初始金钱
	NowMoney  string // 账户余额
	Level     string // 杠杆倍数
}

var Billing BillingType

var BillingArr []BillingType

// 根据下单结果进行模拟持仓
func BillingFun(NowKdata hunter.TradeKdType) {
	fmt.Println("下单总结一次", NowKdata.TimeStr, "收益率", NowPosition.UplRatio)
}

// 模拟数据流动并执行分析交易
func (_this *TestObj) MockData(MockOpt BillingType, TradeKdataOpt hunter.TradeKdataOpt) {
	// 收益结算
	Billing = BillingType{}
	Billing.MockName = MockOpt.MockName
	Billing.NowMoney = MockOpt.InitMoney // 设定当前账户资金
	Billing.Level = MockOpt.Level
	BillingArr = []BillingType{}

	// 交易信息清空
	NowPosition = PositionType{}   // 清空持仓
	PositionArr = []PositionType{} // 清空持仓数组
	OrderArr = []OrderType{}       // 下单数组清空

	global.Run.Println("新建Mock数据",
		mJson.Format(map[string]any{
			"MockName": Billing.MockName,
			"NowMoney": Billing.NowMoney,
			"Level":    Billing.Level,
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
		TradeKdata := hunter.NewTradeKdata(FormatEnd, TradeKdataOpt)
		hunter.TradeKdataList = append(hunter.TradeKdataList, TradeKdata)

		if len(FormatEnd) < TradeKdataOpt.MA_Period {
			continue
		}

		// 开始执行数据整理
		hunter.FormatTradeKdata(TradeKdataOpt)
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

	// 当前持仓与 判断方向不符合时，执行一次下单操作
	if NowPosition.Dir != AnalyDir {
		OnOrder(AnalyDir, NowKdata)
	}

	// 更新持仓状态
	if NowPosition.Dir != 0 { // 当前为持持仓状态，则计算收益率
		UplRatio := mCount.RoseCent(NowKdata.C, NowPosition.OpenAvgPx)
		if NowPosition.Dir < 0 { // 当前为持 空 仓状态
			UplRatio = mCount.Sub("0", NowPosition.UplRatio)
		}
		NowPosition.UplRatio = mCount.Mul(UplRatio, Billing.Level) // 乘以杠杆倍数
	}
	NowPosition.NowTimeStr = NowKdata.TimeStr
	NowPosition.NowC = NowKdata.C
	PositionArr = append(PositionArr, NowPosition)
}

// 下单  参数：dir 下单方向 NowKdata : 当前市场行情
func OnOrder(dir int, NowKdata hunter.TradeKdType) {
	BillingFun(NowKdata)
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
