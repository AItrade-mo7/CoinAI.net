package testHunter

import (
	"CoinAI.net/server/global"
	"github.com/EasyGolang/goTools/mCount"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/EasyGolang/goTools/mTime"
)

func (_this *MockObj) Analy() {
	NowKdata := _this.TradeKdataList[len(_this.TradeKdataList)-1]
	AnalyDir := 0 // 分析的方向，默认为 0 不开仓

	if mCount.Le(NowKdata.CAP_EMA, _this.TradeKdataOpt.CAP_Max) > 0 { // 大于 CAPMax 则开多
		AnalyDir = 1
	}

	if mCount.Le(NowKdata.CAP_EMA, _this.TradeKdataOpt.CAP_Min) < 0 { // 小于 CAPMin 则开空
		AnalyDir = -1
	}

	// 更新持仓状态
	_this.CountPosition()
	_this.PositionArr = append(_this.PositionArr, _this.NowVirtualPosition)

	// 持仓过程中最低和最高盈利比率
	if mCount.Le(_this.NowVirtualPosition.NowUplRatio, _this.Billing.PositionMinRatio.Value) < 0 {
		_this.Billing.PositionMinRatio.Value = _this.NowVirtualPosition.NowUplRatio
		_this.Billing.PositionMinRatio.TimeStr = _this.NowVirtualPosition.NowTimeStr
	}
	if mCount.Le(_this.NowVirtualPosition.NowUplRatio, _this.Billing.PositionMaxRatio.Value) > 0 {
		_this.Billing.PositionMaxRatio.Value = _this.NowVirtualPosition.NowUplRatio
		_this.Billing.PositionMaxRatio.TimeStr = _this.NowVirtualPosition.NowTimeStr
	}

	global.Run.Println(NowKdata.TimeStr, _this.NowVirtualPosition.NowDir, AnalyDir)

	// 当前持仓与 判断方向不符合时，执行一次下单操作
	if _this.NowVirtualPosition.NowDir != AnalyDir {
		_this.OnOrder(AnalyDir)
	}
}

func (_this *MockObj) CountPosition() {
	NowKTradeData := _this.TradeKdataList[len(_this.TradeKdataList)-1]

	_this.NowVirtualPosition.InstID = NowKTradeData.InstID
	_this.NowVirtualPosition.HunterName = _this.HunterName
	_this.NowVirtualPosition.NowTimeStr = NowKTradeData.TimeStr
	_this.NowVirtualPosition.NowTime = mTime.GetTime().TimeUnix
	_this.NowVirtualPosition.NowC = NowKTradeData.C
	_this.NowVirtualPosition.CAP_EMA = NowKTradeData.CAP_EMA
	_this.NowVirtualPosition.EMA = NowKTradeData.EMA
	_this.NowVirtualPosition.HunterConfig = NowKTradeData.Opt

	if _this.NowVirtualPosition.NowDir != 0 { // 当前为持仓状态，则计算收益率
		UplRatio := mCount.RoseCent(NowKTradeData.C, _this.NowVirtualPosition.OpenAvgPx)
		if _this.NowVirtualPosition.NowDir < 0 { // 当前为持空仓状态则翻转该值
			UplRatio = mCount.Sub("0", UplRatio)
		}
		Level := _this.NowVirtualPosition.HunterConfig.MaxTradeLever
		_this.NowVirtualPosition.NowUplRatio = mCount.Mul(UplRatio, mStr.ToStr(Level)) // 乘以杠杆倍数
	}
}

// // 下单  参数：dir 下单方向 NowKdata : 当前市场行情
func (_this *MockObj) OnOrder(dir int) {
	NowKTradeData := _this.TradeKdataList[len(_this.TradeKdataList)-1]

	// 在这里计算当前的 Money
	Upl := mCount.Div(_this.NowVirtualPosition.NowUplRatio, "100")     // 格式化收益率
	ChargeUpl := mCount.Div(_this.NowVirtualPosition.ChargeUpl, "100") // 格式化手续费率

	Money := _this.NowVirtualPosition.Money // 提取 Money
	makeMoney := mCount.Mul(Money, Upl)     // 当前盈利的金钱
	Money = mCount.Add(Money, makeMoney)    // 相加得出当账户剩余资金

	nowCharge := mCount.Mul(Money, ChargeUpl) // 当前产生的手续费
	Money = mCount.Sub(Money, nowCharge)      // 减去手续费
	Money = mCount.CentRound(Money, 3)        // 四舍五入保留三位小数
	_this.NowVirtualPosition.Money = Money    // 保存结果到当前持仓

	// 在这里将当前订单进行结算,相当于平仓了一次
	_this.BillingFun()

	// 同步持仓状态, 相当于下单了
	if dir > 0 {
		// 开多
		_this.NowVirtualPosition.NowDir = 1
	}
	if dir < 0 {
		// 开空
		_this.NowVirtualPosition.NowDir = -1
	}
	// 同步下单价格
	_this.NowVirtualPosition.OpenAvgPx = NowKTradeData.C
	_this.NowVirtualPosition.OpenTimeStr = NowKTradeData.TimeStr
	_this.NowVirtualPosition.OpenTime = mTime.GetTime().TimeUnix

	// 同步平仓状态
	if dir == 0 {
		_this.NowVirtualPosition.NowDir = 0
		_this.NowVirtualPosition.OpenAvgPx = ""
		_this.NowVirtualPosition.OpenTimeStr = ""
		_this.NowVirtualPosition.OpenTime = 0
	}
	// 平仓后未实现盈亏重置为 0
	_this.NowVirtualPosition.NowUplRatio = "0"

	// 下单一次
	_this.OrderArr = append(_this.OrderArr, _this.NowVirtualPosition)
}
