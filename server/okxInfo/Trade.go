package okxInfo

import "github.com/EasyGolang/goTools/mOKX"

var (
	KdataInst mOKX.TypeInst
	TradeInst mOKX.TypeInst
)

var Ticking = make(chan string, 2)

var MaxLen = 500

var NowKdataList = []mOKX.TypeKd{}
