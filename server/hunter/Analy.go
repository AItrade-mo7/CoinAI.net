package hunter

import (
	"CoinAI.net/server/global"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mJson"
)

func Analy() {
	if len(okxInfo.TradeKdataList) < 100 {
		global.LogErr("hunter.Analy 数据长度错误", len(okxInfo.TradeKdataList))
		return
	}

	Last := okxInfo.TradeKdataList[len(okxInfo.TradeKdataList)-1]
	LastPrint := map[string]any{
		"InstID":       Last.InstID,
		"TimeStr":      Last.TimeStr,
		"AllLen":       len(okxInfo.TradeKdataList),
		"C":            Last.C,
		"EMA":          Last.EMA,
		"MA":           Last.MA,
		"RSI":          Last.RSI,
		"CAP_EMA":      Last.CAP_EMA,
		"CAP_MA":       Last.CAP_MA,
		"CAPIdx":       Last.CAPIdx,
		"RsiEmaRegion": Last.RsiEmaRegion,
		"Opt":          Last.Opt,
	}
	global.TradeLog.Println("hunter.Analy 开始分析并执行交易 Last", mJson.Println(LastPrint))
}
