package okxInfo

import "github.com/EasyGolang/goTools/mOKX"

var (
	KdataInst mOKX.TypeInst
	TradeInst mOKX.TypeInst
)

var Ticking = make(chan string, 2)

var MaxLen = 500

var NowKdataList = []mOKX.TypeKd{}

type TradeKdType struct {
	mOKX.TypeKd
	EMA_18    string
	MA_18     string
	RSI_18    string
	CAP_EMA   string
	CAP_MA    string
	CAPIdx    int           //  CAP_EMA 的比值 2 1 0 -1  -2
	RsiRegion int           // Rsi 的震荡区域  -3 -2 -1 0 1 2 3
	PreList   []TradeKdType // 包含当前的，往前数 5 个
}

var TradeKdataList []TradeKdType
