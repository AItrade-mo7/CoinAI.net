package hunter

import (
	"CoinAI.net/server/global"
	"github.com/EasyGolang/goTools/mCount"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/EasyGolang/goTools/mTime"
)

func (_this *HunterObj) Analy() {
	if len(_this.TradeKdataList) < 100 {
		global.LogErr("hunter.Analy 数据长度错误", len(_this.TradeKdataList))
		return
	}

	NowKTradeData := _this.TradeKdataList[len(_this.TradeKdataList)-1]
	LastPrint := map[string]any{
		"InstID":  NowKTradeData.InstID,
		"TimeStr": NowKTradeData.TimeStr,
		"AllLen":  len(_this.TradeKdataList),
		"C":       NowKTradeData.C,
		"EMA":     NowKTradeData.EMA,
		"CAP_EMA": NowKTradeData.CAP_EMA,
		"Opt":     NowKTradeData.Opt,
	}
	global.TradeLog.Println("hunter.Analy 开始分析并执行交易 NowKTradeData", mJson.Format(LastPrint))
	AnalyDir := 0                                                          // 分析的方向，默认为 0 不开仓
	if mCount.Le(NowKTradeData.CAP_EMA, _this.TradeKdataOpt.CAP_Max) > 0 { // 大于 CAPMax 则开多
		AnalyDir = 1
	}

	if mCount.Le(NowKTradeData.CAP_EMA, "-"+_this.TradeKdataOpt.CAP_Max) < 0 { // 小于 负 的 CAPMax 则开空
		AnalyDir = -1
	}

	// 更新持仓状态
	_this.CountPosition()

	// 当前持仓与 判断方向不符合时，执行一次下单操作
	if _this.NowVirtualPosition.NowDir != AnalyDir {
		_this.OnOrder(AnalyDir)
	}
}

func (_this *HunterObj) CountPosition() {
	NowKTradeData := _this.TradeKdataList[len(_this.TradeKdataList)-1]
	_this.NowVirtualPosition.HunterName = _this.HunterName
	_this.NowVirtualPosition.NowTimeStr = NowKTradeData.TimeStr
	_this.NowVirtualPosition.NowTime = mTime.GetTime().TimeUnix
	_this.NowVirtualPosition.NowC = NowKTradeData.C
	_this.NowVirtualPosition.CAP_EMA = NowKTradeData.CAP_EMA
	_this.NowVirtualPosition.InstID = NowKTradeData.InstID
	_this.NowVirtualPosition.HunterConfig = NowKTradeData.Opt

	if _this.NowVirtualPosition.NowDir != 0 { // 当前为持仓状态，则计算收益率
		UplRatio := mCount.RoseCent(NowKTradeData.C, _this.NowVirtualPosition.OpenAvgPx)
		if _this.NowVirtualPosition.NowDir < 0 { // 当前为持空仓状态则翻转该值
			UplRatio = mCount.Sub("0", UplRatio)
		}
		_this.NowVirtualPosition.NowUplRatio = mCount.Mul(UplRatio, mStr.ToStr(_this.TradeKdataOpt.MaxTradeLever)) // 乘以杠杆倍数
		// 这里应该在客户端执行 金钱收益  Money * UplRatio = NowMoney
	}
}

func (_this *HunterObj) OnOrder(dir int) {
	NowKTradeData := _this.TradeKdataList[len(_this.TradeKdataList)-1]

	_this.NowVirtualPosition.OpenAvgPx = NowKTradeData.C
	_this.NowVirtualPosition.OpenTimeStr = NowKTradeData.TimeStr
	_this.NowVirtualPosition.OpenTime = mTime.GetTime().TimeUnix
	// 在这里计算当前的 Money
	Upl := mCount.Div(_this.NowVirtualPosition.NowUplRatio, "100")     // 格式化收益率
	ChargeUpl := mCount.Div(_this.NowVirtualPosition.ChargeUpl, "100") // 格式化手续费率

	Money := _this.NowVirtualPosition.Money   // 提取 Money
	makeMoney := mCount.Mul(Money, Upl)       // 当前盈利的金钱
	nowCharge := mCount.Mul(Money, ChargeUpl) // 当前产生的手续费
	Money = mCount.Add(Money, makeMoney)      // 相加得出当账户总资金量
	Money = mCount.Sub(Money, nowCharge)      // 减去手续费
	Money = mCount.CentRound(Money, 3)        // 四舍五入保留三位小数
	_this.NowVirtualPosition.Money = Money    // 保存结果到当前持仓

	if dir > 0 {
		// 开多
		_this.NowVirtualPosition.NowDir = 1
	}

	if dir < 0 {
		// 开空
		_this.NowVirtualPosition.NowDir = -1
	}

	if dir == 0 {
		// 平仓
		_this.NowVirtualPosition.NowDir = 0
	}
	// 平仓后未实现盈亏重置为 0
	_this.NowVirtualPosition.NowUplRatio = "0"
}
