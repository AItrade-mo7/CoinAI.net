package hunter

import (
	"CoinAI.net/server/global"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mCount"
)

func Analy() {
	if len(okxInfo.TradeKdataList) < 100 || len(okxInfo.TradeKdataList) != len(okxInfo.NowKdataList) {
		global.LogErr("hunter.Analy 数据长度错误", len(okxInfo.TradeKdataList), len(okxInfo.NowKdataList))
		return
	}
	Last := okxInfo.TradeKdataList[len(okxInfo.TradeKdataList)-1]
	// Pre := okxInfo.TradeKdataList[len(okxInfo.TradeKdataList)-2]
	global.RunLog.Println("hunter.Analy 开始分析", Last.TimeStr)

	// CAP_l_p_diff = mCount.Sub(Last.CAP_EMA, Pre.CAP_EMA)
	if mCount.Le(Last.CAP_EMA, "0") > 0 {
		global.TradeLog.Println(Last.TimeStr, 1)
	} else {
		global.TradeLog.Println(Last.TimeStr, -1)
	}
}
