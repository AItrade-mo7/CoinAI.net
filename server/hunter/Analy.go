package hunter

import (
	"CoinAI.net/server/global"
	"github.com/EasyGolang/goTools/mJson"
)

func (_this *HunterObj) Analy() {
	if len(_this.TradeKdataList) < 100 {
		global.LogErr("hunter.Analy 数据长度错误", len(_this.TradeKdataList))
		return
	}

	Last := _this.TradeKdataList[len(_this.TradeKdataList)-1]
	LastPrint := map[string]any{
		"InstID":  Last.InstID,
		"TimeStr": Last.TimeStr,
		"AllLen":  len(_this.TradeKdataList),
		"C":       Last.C,
		"EMA":     Last.EMA,
		"MA":      Last.MA,
		"CAP_EMA": Last.CAP_EMA,
		"CAP_MA":  Last.CAP_MA,
		"Opt":     Last.Opt,
	}
	global.TradeLog.Println("hunter.Analy 开始分析并执行交易 Last", mJson.Format(LastPrint))
}
