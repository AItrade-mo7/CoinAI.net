package testHunter

import (
	"CoinAI.net/server/global"
	"github.com/EasyGolang/goTools/mCount"
	"github.com/EasyGolang/goTools/mFile"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mStr"
)

// // 根据下单结果进行模拟持仓
func (_this *MockObj) BillingFun() {
	NowKTradeData := _this.TradeKdataList[len(_this.TradeKdataList)-1]

	global.KdataLog.Println(_this.Billing.MockName, "下单总结一次",
		NowKTradeData.TimeStr,
		"平仓方向", _this.NowVirtualPosition.NowDir,
		"收益率", _this.NowVirtualPosition.NowUplRatio,
		"结算", _this.NowVirtualPosition.Money,
	)

	// 记录平仓时最大亏损和单次最大盈利
	if mCount.Le(_this.NowVirtualPosition.NowUplRatio, _this.Billing.MaxRatio.Value) > 0 {
		_this.Billing.MaxRatio.Value = _this.NowVirtualPosition.NowUplRatio
		_this.Billing.MaxRatio.TimeStr = _this.NowVirtualPosition.NowTimeStr
	}
	if mCount.Le(_this.NowVirtualPosition.NowUplRatio, _this.Billing.MinRatio.Value) < 0 {
		_this.Billing.MinRatio.Value = _this.NowVirtualPosition.NowUplRatio
		_this.Billing.MinRatio.TimeStr = _this.NowVirtualPosition.NowTimeStr
	}

	// 在这里计算当前的 Money
	Upl := mCount.Div(_this.NowVirtualPosition.NowUplRatio, "100")     // 格式化收益率
	ChargeUpl := mCount.Div(_this.NowVirtualPosition.ChargeUpl, "100") // 格式化手续费率

	Money := _this.NowVirtualPosition.Money // 提取 Money
	makeMoney := mCount.Mul(Money, Upl)     // 当前盈利的金钱
	Money = mCount.Add(Money, makeMoney)    // 相加得出当账户剩余资金

	nowCharge := mCount.Mul(Money, ChargeUpl)              // 当前产生的手续费
	Money = mCount.Sub(Money, nowCharge)                   // 减去手续费
	_this.Billing.ResultMoney = mCount.CentRound(Money, 3) // 四舍五入保留三位小数

	_this.Billing.ChargeAdd = mCount.Add(_this.Billing.ChargeAdd, nowCharge) // 总手续费

	if mCount.Le(_this.Billing.ResultMoney, _this.Billing.MinMoney.Value) < 0 {
		_this.Billing.MinMoney.Value = _this.Billing.ResultMoney
		_this.Billing.MinMoney.TimeStr = _this.NowVirtualPosition.NowTimeStr
	}

	if mCount.Le(_this.Billing.ResultMoney, _this.Billing.MaxMoney.Value) > 0 {
		_this.Billing.MaxMoney.Value = _this.Billing.ResultMoney
		_this.Billing.MaxMoney.TimeStr = _this.NowVirtualPosition.NowTimeStr
	}

	// 盈利计数
	if mCount.Le(_this.NowVirtualPosition.NowUplRatio, "0") > 0 {
		_this.Billing.WinNum++                                                                                        // 盈利次数计数
		_this.Billing.WinUplRatioAdd = mCount.Add(_this.Billing.WinUplRatioAdd, _this.NowVirtualPosition.NowUplRatio) // 盈利比例相加
		_this.Billing.WinMoneyAdd = mCount.Add(_this.Billing.WinMoneyAdd, makeMoney)
	}
	// 亏损计数
	if mCount.Le(_this.NowVirtualPosition.NowUplRatio, "0") < 0 {
		_this.Billing.LoseNum++                                                                                         // 亏损次数计数
		_this.Billing.LoseUplRatioAdd = mCount.Add(_this.Billing.LoseUplRatioAdd, _this.NowVirtualPosition.NowUplRatio) // 盈利比例相加
		_this.Billing.LoseMoneyAdd = mCount.Add(_this.Billing.LoseMoneyAdd, makeMoney)
	}

	if _this.NowVirtualPosition.NowDir == 0 {
		_this.Billing.NilNum++ // 空仓计数
	} else {
		if len(_this.Billing.StartTime) == 0 {
			_this.Billing.StartTime = _this.NowVirtualPosition.OpenTimeStr // 首次开仓时间
		}
	}

	if _this.NowVirtualPosition.NowDir < 0 {
		_this.Billing.SellNum++ // 开空 计数
	}
	if _this.NowVirtualPosition.NowDir > 0 {
		_this.Billing.BuyNum++ // 开多 计数
	}
	_this.Billing.AllNum++ // 总交易计数
}

func (_this *MockObj) ResultCollect() {
	// 记录 整理好的数组
	TradeKdataList_Path := mStr.Join(_this.OutPutDirectory, "/", _this.Billing.InstID, "-TradeKdataList.json")
	mFile.Write(TradeKdataList_Path, string(mJson.ToJson(_this.TradeKdataList)))
	global.Run.Println("TradeKdataList: ", TradeKdataList_Path)

	// 记录 持仓数组
	PositionArr_Path := mStr.Join(_this.OutPutDirectory, "/", _this.Billing.InstID, "-PositionArr.json")
	mFile.Write(PositionArr_Path, string(mJson.ToJson(_this.PositionArr)))
	global.Run.Println("PositionArr: ", PositionArr_Path)

	// 记录 下单数组
	OrderArr_Path := mStr.Join(_this.OutPutDirectory, "/", _this.Billing.InstID, "-OrderArr.json")
	mFile.Write(OrderArr_Path, string(mJson.ToJson(_this.OrderArr)))
	global.Run.Println("OrderArr: ", OrderArr_Path)

	// 记录 交易结果
	Billing_Path := mStr.Join(_this.OutPutDirectory, "/", _this.Billing.InstID, "-Billing.json")
	mFile.Write(Billing_Path, string(mJson.ToJson(_this.Billing)))
	global.Run.Println("Billing: ", Billing_Path)

	Tmp := `交易结果:
参数名称: ${MockName}
InstID: ${InstID}
第一次持仓时间: ${StartTime}
数据结束时间: ${EndTime}
总天数: ${AllDay}
开仓频率: ${OrderRate} 
空仓次数: ${NilNum}
平空次数: ${SellNum}
平多次数: ${BuyNum}
下单总次数: ${AllNum}
盈利次数: ${WinNum}
亏损次数: ${LoseNum}
胜率: ${WinRatio}
盈亏比: ${PLratio}
总盈利比率: ${WinUplRatioAdd}
总盈利总金额: ${WinMoneyAdd}
总亏损比率: ${LoseUplRatioAdd}
亏损总金额: ${LoseMoneyAdd}
平仓后单笔最大盈利比率: ${MaxRatio}
平仓后单笔最小盈利比率: ${MinRatio}
总手续费: ${ChargeAdd}
平仓后历史最低余额: ${MinMoney}
平仓后历史最高余额: ${MaxMoney}
持仓过程中最低盈利比率: ${PositionMinRatio}
持仓过程中最高盈利比率: ${PositionMaxRatio}
初始金钱: ${InitMoney}
最终金钱: ${ResultMoney}
杠杆倍率: ${Level}
`

	// 计算胜率
	ValidAllNum := _this.Billing.SellNum + _this.Billing.BuyNum // 开空次数 + 开多次数
	_this.Billing.WinRatio = mCount.Div(mStr.ToStr(_this.Billing.WinNum), mStr.ToStr(ValidAllNum))

	// 盈亏比 (盈利总金额/盈利次数)  /  (亏损总金额/亏损次数)
	WinMoneyR := mCount.Div(_this.Billing.WinMoneyAdd, mStr.ToStr(_this.Billing.WinNum))
	LoseMoneyR := mCount.Div(mCount.Abs(_this.Billing.LoseMoneyAdd), mStr.ToStr(_this.Billing.LoseNum))
	_this.Billing.PLratio = mCount.Div(WinMoneyR, LoseMoneyR)

	_this.Billing.EndTime = _this.PositionArr[len(_this.PositionArr)-1].NowTimeStr

	_this.Billing.OrderRate = mCount.Div(mStr.ToStr(_this.Billing.AllDay), mStr.ToStr(_this.Billing.AllNum))

	Data := map[string]string{
		"MockName":         _this.Billing.MockName,
		"InstID":           _this.Billing.InstID,
		"AllDay":           mStr.ToStr(_this.Billing.AllDay),            // 总天数 | 结束时计算
		"OrderRate":        _this.Billing.OrderRate,                     // 开仓频率
		"StartTime":        _this.Billing.StartTime,                     // 第一次持仓时间 数组第一个 | 结束时计算
		"EndTime":          _this.Billing.EndTime,                       // 结束时间 数组组后一个
		"NilNum":           mStr.ToStr(_this.Billing.NilNum),            // 空仓次数 平仓后未开仓 NowDir = 0 | 结束时计算
		"SellNum":          mStr.ToStr(_this.Billing.SellNum),           // 开空次数 平空次数 NowDir = -1 | 结束时计算
		"BuyNum":           mStr.ToStr(_this.Billing.BuyNum),            // 开多次数 平多次数 NowDir = 1 | 结束时计算
		"AllNum":           mStr.ToStr(_this.Billing.AllNum),            // 总开仓次数 总的平仓次数 数组长度 | 结束时计算
		"WinNum":           mStr.ToStr(_this.Billing.WinNum),            // 盈利次数 NowUplRatio > 0 的次数
		"LoseNum":          mStr.ToStr(_this.Billing.LoseNum),           // 亏损次数 同 盈利次数
		"WinRatio":         _this.Billing.WinRatio,                      // 胜率 盈利次数/(平空次数+平多次数)
		"PLratio":          _this.Billing.PLratio,                       // 盈亏比
		"WinUplRatioAdd":   _this.Billing.WinUplRatioAdd,                // 总盈利比率 NowUplRatio > 0 的总和
		"WinMoneyAdd":      _this.Billing.WinMoneyAdd,                   // 盈利总金额 1000 块钱 从头计算一次 盈利部分相加
		"LoseUplRatioAdd":  _this.Billing.LoseUplRatioAdd,               // 总亏损比率 同总的盈利比率
		"LoseMoneyAdd":     _this.Billing.LoseMoneyAdd,                  // 亏损总金额 同上
		"MaxRatio":         mJson.ToStr(_this.Billing.MaxRatio),         // 平仓后单笔最大盈利比率   平仓后的记录
		"MinRatio":         mJson.ToStr(_this.Billing.MinRatio),         // 平仓后单笔最小盈利比率
		"ChargeAdd":        _this.Billing.ChargeAdd,                     // 总手续费 同上
		"MinMoney":         mJson.ToStr(_this.Billing.MinMoney),         // 平仓后历史最低余额  遍历一次就知道
		"MaxMoney":         mJson.ToStr(_this.Billing.MaxMoney),         // 平仓后历史最高余额  遍历一次就知道
		"PositionMinRatio": mJson.ToStr(_this.Billing.PositionMinRatio), // 持仓过程中最低盈利比率  // 持仓过程中才知道 结合K线才能得出
		"PositionMaxRatio": mJson.ToStr(_this.Billing.PositionMaxRatio), // 持仓过程中最高盈利比率 // 持仓过程中才知道
		"InitMoney":        _this.Billing.InitMoney,                     // 初始金钱
		"ResultMoney":      _this.Billing.ResultMoney,                   // 最终金钱
		"Level":            _this.Billing.Level,                         // 杠杆倍率
	}

	global.TradeLog.Println(mStr.Temp(Tmp, Data))
}
