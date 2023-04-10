package okxInfo

import "github.com/EasyGolang/goTools/mOKX"

// 交易K线需要的参数
type TradeKdataOpt struct {
	EMA_Period    int    // EMA 步长 171
	CAP_Period    int    // CAP 步长 4
	CAP_Max       string // CAP 判断的边界值 0.2
	MaxTradeLever int
}

type RecordNodeType struct {
	Value   string // 值
	TimeStr string // 产生该值的K线时间
	Time    int64  // 产生该值的实际时间
}

// 模拟持仓的数据
type VirtualPositionType struct {
	InstID       string        `bson:"InstID"`       // 下单币种 | 运行中设置
	HunterName   string        `bson:"HunterName"`   // 策略名称 | 运行中设置
	NowTimeStr   string        `bson:"NowTimeStr"`   // 当前K线时间 | 运行中设置
	NowTime      int64         `bson:"NowTime"`      // 当前实际时间戳 | 运行中设置
	NowC         string        `bson:"NowC"`         // 当前收盘价 | 运行中设置
	CAP_EMA      string        `bson:"CAP_EMA"`      // 当前的 CAP 值 | 运行中设置
	EMA          string        `bson:"EMA"`          // 当前的 EMA 值 | 运行中设置
	HunterConfig TradeKdataOpt `bson:"HunterConfig"` // 当前的交易K线参数  | 运行中设置
	// 下单时设置
	OpenAvgPx   string `bson:"OpenAvgPx"`   // 开仓价格 | 下单时设置
	OpenTimeStr string `bson:"OpenTimeStr"` // 开仓K线时间 | 下单时设置
	OpenTime    int64  `bson:"OpenTime"`    // 开仓实际时间戳  | 下单时设置
	NowDir      int    `bson:"NowDir"`      // 当前持仓状态 没持仓0  持多仓 1  持空仓 -1 | 初始化设置为0，下单时更新

	// 通过原始数据计算得出
	InitMoney   string `bson:"InitMoney"`   // 初始金钱 | 固定值初始值设置
	ChargeUpl   string `bson:"ChargeUpl"`   // 当前手续费率 | 固定值初始值设置
	NowUplRatio string `bson:"NowUplRatio"` // 当前未实现收益率(计算得出) | 运行中设置
	Money       string `bson:"Money"`       // 账户当前余额 | 如果没有初始值设置一次，下单时计算
}

type TradeKdType struct {
	mOKX.TypeKd
	EMA string // EMA 值
	// MA      string // MA 值
	CAP_EMA string // 基于 EMA 的 平滑点数 0-100 的浮点类型
	// CAP_MA  string // 基于 EMA 的 平滑点数 0-100 的浮点类型
	Opt TradeKdataOpt
}

// 对外展示策略的数据
type HunterData struct {
	HunterName         string              // 策略的名字
	Describe           string              // 描述
	InstID             string              // 当前策略主打币种
	TradeInst          mOKX.TypeInst       // 交易的 InstID SWAP
	KdataInst          mOKX.TypeInst       // K线的 InstID SPOT
	NowKdataList       []mOKX.TypeKd       // 现货的原始K线
	TradeKdataList     []TradeKdType       // 计算好各种指标之后的K线
	TradeKdataOpt      TradeKdataOpt       // 当前参数
	NowVirtualPosition VirtualPositionType // 当前的虚拟持仓
}

var NowHunterData = make(map[string]HunterData)

// 最优参数
var CoinTradeConfig = make(map[string]TradeKdataOpt)

func OkxInfoInit() {
	// 加入 Hunter Auto
	NowHunterData["Auto"] = HunterData{
		HunterName: "Auto",
		Describe:   "根据市场情况为您的账户选择其中一个策略执行交易【目前此功能尚在开发中】",
	}

	// 设置最优参数
	CoinTradeConfig = map[string]TradeKdataOpt{
		"BTC-USDT": {
			EMA_Period:    86,
			CAP_Period:    2,
			CAP_Max:       "0.5",
			MaxTradeLever: 2,
		},
		"ETH-USDT": {
			EMA_Period:    78,
			CAP_Period:    2,
			CAP_Max:       "0.5",
			MaxTradeLever: 4,
		},
	}
}
