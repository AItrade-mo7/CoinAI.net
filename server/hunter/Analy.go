package hunter

import (
	"CoinAI.net/server/global"
	"CoinAI.net/server/okxInfo"
)

func Analy() {
	if len(okxInfo.TradeKdataList) < 100 || len(okxInfo.TradeKdataList) != len(okxInfo.NowKdataList) {
		global.LogErr("hunter.Analy 数据长度错误", len(okxInfo.TradeKdataList), len(okxInfo.NowKdataList))
		return
	}
	Last := okxInfo.TradeKdataList[len(okxInfo.TradeKdataList)-1]

	global.RunLog.Println("hunter.Analy 开始分析", Last.TimeStr)
}
