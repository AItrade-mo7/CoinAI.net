package okxInfo

import "github.com/EasyGolang/goTools/mOKX"

type TradeKdataOpt struct {
	EMA_Period    int    // 171
	CAP_Period    int    // 4
	CAP_Max       string // 0.2
	MaxTradeLever int
}

type RecordType struct {
	Value   string
	TimeStr string
}

type VirtualPositionType struct {
	InstID       string        // 下单币种 | 运行中设置
	HunterName   string        // 策略名称 | 运行中设置
	Dir          int           // 没持仓0  持多仓 1  持空仓 -1 | 初始化设置，下单时设置
	OpenAvgPx    string        // 开仓价格 | 下单时设置
	OpenTimeStr  string        // 开仓K线时间 | 下单时设置
	OpenTime     int64         // 开仓实际时间戳
	NowTimeStr   string        // 当前K线时间 | 运行中设置
	NowTime      int64         // 当前实际时间戳
	NowC         string        // 当前收盘价 | 运行中设置
	UplRatio     string        // 未实现收益率 | 运行中设置
	CAP_EMA      string        // 当前的 CAP 值 | 运行中设置
	InitMoney    string        // 初始金钱 | 初始值设置
	Money        string        // 账户当前余额 | 初始值设置，平仓时计算
	Charge       string        // 当前是手续费率 | 初始值设置
	HunterConfig TradeKdataOpt // 运行中设置
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
		EMA_Period:    86,
		CAP_Period:    2,
		MaxTradeLever: 2,
		CAP_Max:       "0.5",
	},
	"ETH-USDT": {
		EMA_Period:    78,
		CAP_Period:    2,
		MaxTradeLever: 4,
		CAP_Max:       "0.5",
	},
}
