package okxInfo

import "github.com/EasyGolang/goTools/mOKX"

var HunterTicking = make(chan string, 2) // 计算频率

var KdataInst mOKX.TypeInst // 这里一定为现货

var TradeInst mOKX.TypeInst // 这里一定为合约

var MaxLen = 900

var NowKdataList = []mOKX.TypeKd{}

type TradeKdType struct {
	mOKX.TypeKd
	EMA       string // EMA 值
	RSI       string // RSI 的值
	CAP_EMA   string // 基于 EMA 的 平滑点数 0-100 的浮点类型
	RsiRegion int    // 整型 Rsi 的震荡区域  -3 -2 -1 0 1 2 3
}

var TradeKdataList []TradeKdType
