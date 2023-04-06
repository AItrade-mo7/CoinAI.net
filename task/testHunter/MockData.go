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

// 模拟数据流动并执行分析交易
func (_this *MockObj) MockRun() {
	// 清理 TradeKdataList
	lossVal := mCount.Mul(_this.Billing.InitMoney, "0.3") // 当余额 低于 30% 时 判定为 亏完

	// 清理 TradeKdataList
	_this.TradeKdataList = []okxInfo.TradeKdType{}
	TradeKlineObj := hunter.NewTradeKdataObj(_this.TradeKdataOpt)

	FormatEnd := []mOKX.TypeKd{}
	MaxLen := 900
	for _, Kdata := range _this.RunKdataList {
		FormatEnd = append(FormatEnd, Kdata)

		if len(FormatEnd) < _this.TradeKdataOpt.EMA_Period+1 {
			continue
		}

		TradeKdata := TradeKlineObj.NewTradeKdata(FormatEnd)
		_this.TradeKdataList = append(_this.TradeKdataList, TradeKdata)

		// 开始执行分析交易
		_this.Analy()

		if len(FormatEnd)-MaxLen > 0 {
			FormatEnd = FormatEnd[len(FormatEnd)-MaxLen:]
		}

		if mCount.Le(_this.NowPosition.UplRatio, "-45") < 0 {
			global.Log.Println("爆仓！", _this.Billing.MockName, Kdata.TimeStr)
			break
		}

		if mCount.Le(_this.Billing.Money, lossVal) < 0 {
			global.Log.Println("亏完！", _this.Billing.MockName, Kdata.TimeStr)
			break
		}
	}

	if len(_this.TradeKdataList) > 0 {
		_this.Billing.EndTime = _this.TradeKdataList[len(_this.TradeKdataList)-1].TimeStr
	}

	// 搜集和整理结果
	global.TradeLog.Println(" ===== 分析交易结束 ===== ", _this.Billing.MockName)
	_this.ResultCollect()
}

func (_this *MockObj) Analy() {
	NowKdata := _this.TradeKdataList[len(_this.TradeKdataList)-1]
	AnalyDir := 0 // 分析的方向，默认为 0 不开仓

	if mCount.Le(NowKdata.CAP_EMA, "0") > 0 { // 大于 0 则开多
		AnalyDir = 1
	}

	if mCount.Le(NowKdata.CAP_EMA, "0") < 0 { // 小于 0 则开空
		AnalyDir = -1
	}

	// 更新持仓状态
	if _this.NowPosition.Dir != 0 { // 当前为持持仓状态，则计算收益率
		UplRatio := mCount.RoseCent(NowKdata.C, _this.NowPosition.OpenAvgPx)
		if _this.NowPosition.Dir < 0 { // 当前为持 空 仓状态 则翻转该值
			UplRatio = mCount.Sub("0", UplRatio)
		}
		_this.NowPosition.UplRatio = mCount.Mul(UplRatio, _this.Billing.Level) // 乘以杠杆倍数
	}
	_this.NowPosition.NowTimeStr = NowKdata.TimeStr
	_this.NowPosition.NowC = NowKdata.C
	_this.PositionArr = append(_this.PositionArr, _this.NowPosition)

	if mCount.Le(_this.NowPosition.UplRatio, _this.Billing.PositionMinRatio.Value) < 0 {
		_this.Billing.PositionMinRatio.Value = _this.NowPosition.UplRatio
		_this.Billing.PositionMinRatio.TimeStr = _this.NowPosition.NowTimeStr
	}
	if mCount.Le(_this.NowPosition.UplRatio, _this.Billing.PositionMaxRatio.Value) > 0 {
		_this.Billing.PositionMaxRatio.Value = _this.NowPosition.UplRatio
		_this.Billing.PositionMaxRatio.TimeStr = _this.NowPosition.NowTimeStr
	}

	// 当前持仓与 判断方向不符合时，执行一次下单操作
	if _this.NowPosition.Dir != AnalyDir {
		_this.OnOrder(AnalyDir, NowKdata)
	}

	global.KdataLog.Println(_this.Billing.MockName, NowKdata.TimeStr, AnalyDir)
}

// 根据下单结果进行模拟持仓
func (_this *MockObj) BillingFun(NowKdata okxInfo.TradeKdType) {
	fmt.Println(_this.Billing.MockName, "下单总结一次",
		NowKdata.TimeStr,
		"持仓方向", _this.NowPosition.Dir,
		"收益率", _this.NowPosition.UplRatio,
		"结算", _this.Billing.Money,
	)

	// 记录单次最大亏损和单次最大盈利
	if mCount.Le(_this.NowPosition.UplRatio, _this.Billing.MaxRatio.Value) > 0 {
		_this.Billing.MaxRatio.Value = _this.NowPosition.UplRatio
		_this.Billing.MaxRatio.TimeStr = _this.NowPosition.NowTimeStr
	}
	if mCount.Le(_this.NowPosition.UplRatio, _this.Billing.MinRatio.Value) < 0 {
		_this.Billing.MinRatio.Value = _this.NowPosition.UplRatio
		_this.Billing.MinRatio.TimeStr = _this.NowPosition.NowTimeStr
	}

	Upl := mCount.Div(_this.NowPosition.UplRatio, "100") // 格式化收益率
	ChargeUpl := mCount.Div(_this.Billing.Charge, "100") // 格式化手续费率

	makeMoney := mCount.Mul(_this.Billing.Money, Upl)                // 当前盈利的金钱
	_this.Billing.Money = mCount.Add(_this.Billing.Money, makeMoney) // 相加得出当账户总资金量

	nowCharge := mCount.Mul(_this.Billing.Money, ChargeUpl)                  // 当前产生的手续费
	_this.Billing.Money = mCount.Sub(_this.Billing.Money, nowCharge)         // 减去手续费
	_this.Billing.ChargeAll = mCount.Add(_this.Billing.ChargeAll, nowCharge) // 记录一下手续费

	_this.Billing.Money = mCount.CentRound(_this.Billing.Money, 3)         // 四舍五入保留两位小数
	_this.Billing.ChargeAll = mCount.CentRound(_this.Billing.ChargeAll, 3) // 四舍五入保留两位小数

	if mCount.Le(_this.Billing.Money, _this.Billing.MinMoney.Value) < 0 {
		_this.Billing.MinMoney.Value = _this.Billing.Money
		_this.Billing.MinMoney.TimeStr = _this.NowPosition.NowTimeStr
	}

	if mCount.Le(_this.Billing.Money, _this.Billing.MaxMoney.Value) > 0 {
		_this.Billing.MaxMoney.Value = _this.Billing.Money
		_this.Billing.MaxMoney.TimeStr = _this.NowPosition.NowTimeStr
	}

	// 盈利计数
	if mCount.Le(_this.NowPosition.UplRatio, "0") > 0 {
		_this.Billing.Win++                                                                     // 盈利次数计数
		_this.Billing.WinRatio = mCount.Add(_this.NowPosition.UplRatio, _this.Billing.WinRatio) // 盈利比例相加
		_this.Billing.WinMoney = mCount.Add(_this.Billing.WinMoney, makeMoney)
	}
	// 亏损计数
	if mCount.Le(_this.NowPosition.UplRatio, "0") < 0 {
		_this.Billing.Lose++                                                                      // 亏损次数计数
		_this.Billing.LoseRatio = mCount.Add(_this.NowPosition.UplRatio, _this.Billing.LoseRatio) // 盈利比例相加
		_this.Billing.LoseMoney = mCount.Add(_this.Billing.LoseMoney, makeMoney)
	}

	if _this.NowPosition.Dir == 0 {
		_this.Billing.NilNum++ // 空仓计数
	} else {
		if len(_this.Billing.StartTime) == 0 {
			_this.Billing.StartTime = _this.NowPosition.OpenTimeStr // 首次开仓时间
		}
	}

	if _this.NowPosition.Dir < 0 {
		_this.Billing.SellNum++ // 开空 计数
	}
	if _this.NowPosition.Dir > 0 {
		_this.Billing.BuyNum++ // 开多 计数
	}
	_this.Billing.AllNum++ // 总交易计数
}

// 下单  参数：dir 下单方向 NowKdata : 当前市场行情
func (_this *MockObj) OnOrder(dir int, NowKdata okxInfo.TradeKdType) {
	_this.BillingFun(NowKdata) // 下单之前 计算一次收益

	if dir > 0 { // 开多
		// 下订单
		_this.OrderArr = append(_this.OrderArr, OrderType{
			Type:    "Buy",            // 下多单
			InstID:  NowKdata.InstID,  // 下单币种
			AvgPx:   NowKdata.C,       // 记录下单价格
			TimeStr: NowKdata.TimeStr, // 记录下单时间
		})
		// 更新持仓状态
		_this.NowPosition = PositionType{
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
		_this.OrderArr = append(_this.OrderArr, OrderType{
			Type:    "Sell",           // 下空单
			InstID:  NowKdata.InstID,  // 下单币种
			AvgPx:   NowKdata.C,       // 记录下单价格
			TimeStr: NowKdata.TimeStr, // 记录下单时间
		})
		// 更新持仓状态
		_this.NowPosition = PositionType{
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
		_this.OrderArr = append(_this.OrderArr, OrderType{
			Type:    "Close",          // 平仓
			InstID:  NowKdata.InstID,  // 下单币种
			AvgPx:   NowKdata.C,       // 记录下单价格
			TimeStr: NowKdata.TimeStr, // 记录下单时间
		})
		// 更新为空仓状态
		_this.NowPosition = PositionType{
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

func (_this *MockObj) ResultCollect() {
	// 记录 整理好的数组
	TradeKdataList_Path := mStr.Join(config.Dir.JsonData, "/", _this.Billing.MockName, "-TradeKdataList.json")
	mFile.Write(TradeKdataList_Path, string(mJson.ToJson(_this.TradeKdataList)))
	global.Run.Println("TradeKdataList: ", TradeKdataList_Path)

	// 记录 持仓数组
	PositionArr_Path := mStr.Join(config.Dir.JsonData, "/", _this.Billing.MockName, "-PositionArr.json")
	mFile.Write(PositionArr_Path, string(mJson.ToJson(_this.PositionArr)))
	global.Run.Println("PositionArr: ", PositionArr_Path)

	// 记录 下单数组
	OrderArr_Path := mStr.Join(config.Dir.JsonData, "/", _this.Billing.MockName, "-OrderArr.json")
	mFile.Write(OrderArr_Path, string(mJson.ToJson(_this.OrderArr)))
	global.Run.Println("OrderArr: ", OrderArr_Path)

	// 记录 交易结果
	Billing_Path := mStr.Join(config.Dir.JsonData, "/", _this.Billing.MockName, "-Billing.json")
	mFile.Write(Billing_Path, string(mJson.ToJson(_this.Billing)))
	global.Run.Println("Billing: ", Billing_Path)

	Tmp := `交易结果:
InstID: ${InstID}
第一次持仓时间: ${StartTime}
数据结束时间: ${EndTime}
总天数: ${Days}
空仓次数: ${NilNum}
开空次数: ${SellNum}
开多次数: ${BuyNum}
总开仓次数: ${AllNum}
盈利次数: ${Win}
总盈利比率: ${WinRatio}
亏损次数: ${Lose}
总亏损比率: ${LoseRatio}
平仓后单笔最大盈利比率: ${MaxRatio}
平仓后单笔最小盈利比率: ${MinRatio}
手续费率: ${Charge}
总手续费: ${ChargeAll}
参数名称: ${MockName}
初始金钱: ${InitMoney}
账户当前余额: ${Money}
平仓后历史最低余额: ${MinMoney}
平仓后历史最高余额: ${MaxMoney}
持仓过程中最低盈利比率: ${PositionMinRatio}
持仓过程中最高盈利比率: ${PositionMaxRatio}
盈利总金额: ${WinMoney}
亏损总金额: ${LoseMoney}
杠杆倍数: ${Level}
胜率: ${WinRatioAll}
平均盈利利率: ${AveWinRatio}
平均亏损利率: ${AveLoseRatio}
盈亏比: ${PLratio}
`

	WinRatioAll := mCount.Div(mStr.ToStr(_this.Billing.Win), mStr.ToStr(_this.Billing.AllNum))

	LoseMoneyAbs := mCount.Abs(mStr.ToStr(_this.Billing.LoseMoney))

	AveWinRatio := mCount.Div(mStr.ToStr(_this.Billing.WinMoney), mStr.ToStr(_this.Billing.Win))
	AveLoseRatio := mCount.Div(LoseMoneyAbs, mStr.ToStr(_this.Billing.Lose))
	PLratio := mCount.Div(AveWinRatio, AveLoseRatio)

	Data := map[string]string{
		"InstID":           _this.Billing.InstID,
		"StartTime":        _this.Billing.StartTime, // 开始时间
		"EndTime":          _this.Billing.EndTime,   // 结束时间
		"Days":             mStr.ToStr(_this.Billing.Days),
		"NilNum":           mStr.ToStr(_this.Billing.NilNum),
		"SellNum":          mStr.ToStr(_this.Billing.SellNum),
		"BuyNum":           mStr.ToStr(_this.Billing.BuyNum),
		"AllNum":           mStr.ToStr(_this.Billing.AllNum),
		"Win":              mStr.ToStr(_this.Billing.Win),
		"WinRatio":         _this.Billing.WinRatio,
		"Lose":             mStr.ToStr(_this.Billing.Lose),
		"LoseRatio":        _this.Billing.LoseRatio,
		"MaxRatio":         mJson.ToStr(_this.Billing.MaxRatio),
		"MinRatio":         mJson.ToStr(_this.Billing.MinRatio),
		"Charge":           _this.Billing.Charge,
		"ChargeAll":        _this.Billing.ChargeAll,
		"MockName":         _this.Billing.MockName,
		"InitMoney":        _this.Billing.InitMoney,
		"Money":            _this.Billing.Money,
		"MinMoney":         mJson.ToStr(_this.Billing.MinMoney),
		"MaxMoney":         mJson.ToStr(_this.Billing.MaxMoney),
		"PositionMinRatio": mJson.ToStr(_this.Billing.PositionMinRatio),
		"PositionMaxRatio": mJson.ToStr(_this.Billing.PositionMaxRatio),
		"WinMoney":         _this.Billing.WinMoney,
		"LoseMoney":        _this.Billing.LoseMoney,
		"Level":            _this.Billing.Level,
		"WinRatioAll":      WinRatioAll,
		"AveWinRatio":      AveWinRatio,
		"AveLoseRatio":     AveLoseRatio,
		"PLratio":          PLratio,
	}

	global.TradeLog.Println(mStr.Temp(Tmp, Data))
}
