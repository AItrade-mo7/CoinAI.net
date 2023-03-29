package testHunter

import (
	"fmt"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/hunter"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mCount"
	"github.com/EasyGolang/goTools/mFile"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mStr"
)

/*
模拟数据流动
*/
var (
	MockName = ""

	TradeKdataOpt = okxInfo.TradeKdataOpt{}

	TradeKdataList = []okxInfo.TradeKdType{}
)

// 开仓信息记录
type PositionType struct {
	Dir         int    // 开仓方向
	OpenAvgPx   string // 开仓价格
	OpenTimeStr string // 开仓时间
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
	NowMoney    string         // 账户余额
	Level       string         // 杠杆倍数
)

type BillingOpt struct {
	MockName  string // 名字
	InitMoney string // 初始金钱
	Level     string // 杠杆倍数
}

// 模拟数据流动并执行分析交易
func (_this *TestObj) MockData(MockOpt BillingOpt, opt okxInfo.TradeKdataOpt) {
	TradeKdataOpt = opt         // 交易参数
	MockName = MockOpt.MockName // MockName

	// 交易信息清空
	NowPosition = PositionType{}   // 清空持仓
	PositionArr = []PositionType{} // 清空持仓数组
	OrderArr = []OrderType{}       // 下单数组清空
	NowMoney = MockOpt.InitMoney   // 设定当前账户资金

	// 清理 TradeKdataList
	TradeKdataList = []okxInfo.TradeKdType{}
	hunter.EMA_Arr = []string{}
	hunter.MA_Arr = []string{}
	hunter.RSI_Arr = []string{}

	global.Run.Println("新建Mock数据", MockName, mJson.Format(TradeKdataOpt))

	global.TradeLog.Println(" ============== 开始分析和交易 ============== ", MockName)
	FormatEnd := []mOKX.TypeKd{}
	for _, Kdata := range _this.KdataList {
		FormatEnd = append(FormatEnd, Kdata)
		TradeKdata := hunter.NewTradeKdata(FormatEnd, TradeKdataOpt)
		TradeKdataList = append(TradeKdataList, TradeKdata)

		fmt.Println(len(hunter.EMA_Arr))

		if len(TradeKdataList) > 100 {
			// 开始执行分析
			Analy()
		}
	}

	// 记录 整理好的数组
	TradeKdataList_Path := mStr.Join(config.Dir.JsonData, "/", MockName, "-TradeKdataList.json")
	mFile.Write(TradeKdataList_Path, string(mJson.ToJson(TradeKdataList)))
	global.Run.Println("TradeKdataList: ", TradeKdataList_Path)

	// 记录 持仓数组
	PositionArr_Path := mStr.Join(config.Dir.JsonData, "/", MockName, "-PositionArr.json")
	mFile.Write(PositionArr_Path, string(mJson.ToJson(PositionArr)))
	global.Run.Println("PositionArr: ", PositionArr_Path)

	// 记录 下单数组
	OrderArr_Path := mStr.Join(config.Dir.JsonData, "/", MockName, "-OrderArr.json")
	mFile.Write(OrderArr_Path, string(mJson.ToJson(OrderArr)))
	global.Run.Println("OrderArr: ", OrderArr_Path)
}

func Analy() {
	NowKdata := TradeKdataList[len(TradeKdataList)-1]
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
		NowPosition.UplRatio = UplRatio // 乘以杠杆倍数
	}
	PositionArr = append(PositionArr, NowPosition)

	// 记录一下持仓
}

// 下单  参数：dir 下单方向 NowKdata : 当前市场行情
func OnOrder(dir int, NowKdata okxInfo.TradeKdType) {
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
			Dir:         1,                // 持仓多方向
			OpenAvgPx:   NowKdata.C,       // 持仓价格
			UplRatio:    "0",              // 当前收益率
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
			Dir:         -1,               // 持仓空方向
			OpenAvgPx:   NowKdata.C,       // 持仓价格
			UplRatio:    "0",              // 当前收益率
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
		// 更新持仓状态
		NowPosition = PositionType{}
	}
}

// 根据下单结果进行模拟持仓
func Billing() {
}
