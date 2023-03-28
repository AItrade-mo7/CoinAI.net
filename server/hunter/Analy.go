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
	// Last := okxInfo.TradeKdataList[len(okxInfo.TradeKdataList)-1]
	// Pre := okxInfo.TradeKdataList[len(okxInfo.TradeKdataList)-2]
	// global.RunLog.Println("hunter.Analy 开始分析", Last.TimeStr)

	// CAP_l_p_diff = mCount.Sub(Last.CAP_EMA, Pre.CAP_EMA)

	// global.TradeLog.Printf(
	// 	"%v EMA:%8v CAP_EMA:%8v %2v; MA:%8v CAP_MA:%8v %2v \n",
	// 	Last.TimeStr,                 // 1
	// 	Last.EMA_18,                  // 2
	// 	Last.CAP_EMA,                 // 3
	// 	mCount.Le(Last.CAP_EMA, "0"), // 4
	// 	Last.MA_18,                   // 5
	// 	Last.CAP_MA,                  // 6
	// 	mCount.Le(Last.CAP_MA, "0"),  // 7
	// )
}
