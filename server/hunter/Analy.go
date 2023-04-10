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
	_this.PositionArr = append(_this.PositionArr, _this.NowVirtualPosition)

	global.Run.Println(NowKTradeData.TimeStr, _this.NowVirtualPosition.NowDir, AnalyDir)

	// 当前持仓与 判断方向不符合时，执行一次下单操作
	if _this.NowVirtualPosition.NowDir != AnalyDir {
		_this.OnOrder(AnalyDir)
	}
}

func (_this *HunterObj) CountPosition() {
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
