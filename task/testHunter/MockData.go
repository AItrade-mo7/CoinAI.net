package testHunter

import (
	"CoinAI.net/server/global"
	"CoinAI.net/server/hunter"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mCount"
	"github.com/EasyGolang/goTools/mOKX"
)

// 模拟数据流动并执行分析交易
func (_this *MockObj) MockRun() BillingType {
	// 清理 TradeKdataList
	_this.TradeKdataList = []okxInfo.TradeKdType{}
	TradeKlineObj := hunter.NewTradeKdataObj(_this.TradeKdataOpt)

	MaxLen := 600
	FormatEnd := []mOKX.TypeKd{}
	for _, Kdata := range _this.RunKdataList {
		FormatEnd = append(FormatEnd, Kdata)

		if len(FormatEnd) < 600 {
			continue
		}

		TradeKdata := TradeKlineObj.NewTradeKdata(FormatEnd)
		_this.TradeKdataList = append(_this.TradeKdataList, TradeKdata)

		// 开始执行分析交易
		_this.Analy()

		if len(FormatEnd)-MaxLen > 0 {
			FormatEnd = FormatEnd[len(FormatEnd)-MaxLen:]
		}

		if mCount.Le(_this.NowVirtualPosition.NowUplRatio, "-45") < 0 {
			global.Log.Println("爆仓！", _this.Billing.MockName, Kdata.TimeStr)
			break
		}
	}

	if len(_this.TradeKdataList) > 0 {
		_this.Billing.EndTime = _this.TradeKdataList[len(_this.TradeKdataList)-1].TimeStr
	}

	// 搜集和整理结果
	global.TradeLog.Println(" ===== 分析交易结束 ===== ", _this.Billing.MockName)
	// _this.ResultCollect()

	// 在这里抛出结果
	return _this.Billing
}
