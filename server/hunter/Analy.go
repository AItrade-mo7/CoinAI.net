package hunter

import (
	"CoinAI.net/server/global"
	"CoinAI.net/server/okxInfo"
)

func Analy() {
	if len(okxInfo.TradeKdataList) < 100 {
		global.LogErr("hunter.Analy 数据长度错误", len(okxInfo.TradeKdataList))
		return
	}
	Last := okxInfo.TradeKdataList[len(okxInfo.TradeKdataList)-1]

	global.TradeLog.Println("hunter.Analy 开始分析并执行交易", Last.TimeStr)
}
