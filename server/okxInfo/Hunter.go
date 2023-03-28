package okxInfo

import "github.com/EasyGolang/goTools/mOKX"

var HunterTicking = make(chan string, 2) // 计算频率

var KdataInst mOKX.TypeInst // 这里一定为现货

var TradeInst mOKX.TypeInst // 这里一定为合约

var MaxLen = 900

var NowKdataList = []mOKX.TypeKd{}

var HLPerLeVel = 3 // 涨跌幅等级  按照涨跌幅度排名，  数字越大越不稳定

type TradeKdType struct {
	mOKX.TypeKd
	EMA          string // EMA 值
	MA           string // MA 值
	RSI          string // RSI 的值
	RSI_EMA      string // RSI 的值
	CAP_EMA      string // 基于 EMA 的 平滑点数 0-100 的浮点类型
	CAP_MA       string // 基于 EMA 的 平滑点数 0-100 的浮点类型
	CAPIdx       int
	RsiEmaRegion int // 整型 Rsi 的震荡区域  -3 -2 -1 0 1 2 3
	Opt          TradeKdataOpt
}

type TradeKdataOpt struct {
	MA_Period      int // 108
	RSI_Period     int // 18
	RSI_EMA_Period int // 14
	CAP_Period     int // 3
}

var TradeKdataList []TradeKdType
