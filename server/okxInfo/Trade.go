package okxInfo

import "github.com/EasyGolang/goTools/mOKX"

var (
	KdataInst mOKX.TypeInst
	TradeInst mOKX.TypeInst
)

var Ticking = make(chan string, 2)

type TradeKdType struct {
	mOKX.TypeKd
	EMA_18  string
	MA_18   string
	RSI_18  string
	CAP_EMA string
	CAP_MA  string
}

var MaxLen = 500

var NowKdata = []mOKX.TypeKd{}
