package hunter

import (
	"CoinAI.net/server/global"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mOKX"
)

type TradeKdType struct {
	mOKX.TypeKd
	EMA_18  string
	MA_18   string
	RSI_18  string
	CAP_EMA string
	CAP_MA  string
}

var TradeKdataList []TradeKdType

func Analy() {
	if len(TradeKdataList) < 100 || len(TradeKdataList) != len(okxInfo.NowKdataList) {
		global.LogErr("hunter.Analy 数据长度错误", len(TradeKdataList))
		return
	}
	Last := TradeKdataList[len(TradeKdataList)-1]

	global.RunLog.Println("hunter.Analy 开始分析", Last.TimeStr)
}
