package okxInfo

import "github.com/EasyGolang/goTools/mOKX"

var HunterTicking = make(chan string, 2) // 计算频率

var KdataInst mOKX.TypeInst // 这里一定为现货

var TradeInst mOKX.TypeInst // 这里一定为合约

var MaxLen = 900

var NowKdataList = []mOKX.TypeKd{}

type TradeKdType struct {
	mOKX.TypeKd
	EMA_18    string // 与原价格接近 string 类型
	MA_18     string // 与原价格接近 string 类型
	RSI_18    string // 0-100 的浮点类型
	RSI_EMA_9 string // 0-100 的浮点类型
	CAP_EMA   string // 0-100 的浮点类型
	CAP_MA    string // 浮点类型
	CAPIdx    int    // 整型  CAP_EMA 的比值 2 1 0 -1  -2
	RsiRegion int    // 整型 Rsi 的震荡区域  -3 -2 -1 0 1 2 3
}

var TradeKdataList []TradeKdType
