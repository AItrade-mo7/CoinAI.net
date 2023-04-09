package okxInfo

import "github.com/EasyGolang/goTools/mOKX"

// Hunter 内部 数据 同步
type TradeKdataOpt struct {
	EMA_Period    int    // 171
	CAP_Period    int    // 4
	CAP_Max       string // 0.2
	MaxTradeLever int
}

type TradeKdType struct {
	mOKX.TypeKd
	EMA string // EMA 值
	// MA      string // MA 值
	CAP_EMA string // 基于 EMA 的 平滑点数 0-100 的浮点类型
	// CAP_MA  string // 基于 EMA 的 平滑点数 0-100 的浮点类型
	Opt TradeKdataOpt
}

type HunterData struct {
	HunterName     string // 策略的名字
	Describe       string // 描述
	MaxLen         int
	TradeInst      mOKX.TypeInst // 交易的 InstID SWAP
	KdataInst      mOKX.TypeInst // K线的 InstID SPOT
	NowKdataList   []mOKX.TypeKd // 现货的原始K线
	TradeKdataList []TradeKdType // 计算好各种指标之后的K线
	TradeKdataOpt  TradeKdataOpt
}

var NowHunterData = make(map[string]HunterData)

// 最优参数

var CoinTradeConfig = map[string]TradeKdataOpt{
	"BTC-USDT": {
		EMA_Period:    171,
		CAP_Period:    4,
		MaxTradeLever: 2,
	},
	"ETH-USDT": {
		EMA_Period:    77,
		CAP_Period:    3,
		MaxTradeLever: 4,
	},
}
