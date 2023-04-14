package okxInfo

import (
	"CoinAI.net/server/global/dbType"
	"github.com/EasyGolang/goTools/mOKX"
)

type RecordNodeType struct {
	Value   string // 值
	TimeStr string // 产生该值的K线时间
	Time    int64  // 产生该值的实际时间
}

type TradeKdType struct {
	mOKX.TypeKd
	EMA string // EMA 值
	// MA      string // MA 值
	CAP_EMA string // 基于 EMA 的 平滑点数 0-100 的浮点类型
	// CAP_MA  string // 基于 EMA 的 平滑点数 0-100 的浮点类型
	Opt dbType.TradeKdataOpt
}

// 对外展示策略的数据
type HunterData struct {
	HunterName         string                     // 策略的名字
	Describe           string                     // 描述
	InstID             string                     // 当前策略主打币种
	TradeInst          mOKX.TypeInst              // 交易的 InstID SWAP
	KdataInst          mOKX.TypeInst              // K线的 InstID SPOT
	NowKdataList       []mOKX.TypeKd              // 现货的原始K线
	TradeKdataList     []TradeKdType              // 计算好各种指标之后的K线
	TradeKdataOpt      dbType.TradeKdataOpt       // 当前参数
	NowVirtualPosition dbType.VirtualPositionType // 当前的虚拟持仓
}

var NowHunterData = make(map[string]HunterData)

// 最优参数
var CoinTradeConfig = make(map[string]dbType.TradeKdataOpt)

func OkxInfoInit() {
	// 加入 Hunter Auto
	NowHunterData["Auto"] = HunterData{
		HunterName: "Auto",
		Describe:   "根据市场情况为您的账户选择其中一个策略执行交易【目前此功能尚在开发中】",
	}

	// 设置最优参数
	CoinTradeConfig = map[string]dbType.TradeKdataOpt{
		"BTC-USDT": {
			EMA_Period:    272, // 参数已确定  2023-04-11 18:14
			CAP_Period:    5,
			CAP_Max:       "2.5",
			MaxTradeLever: 5,
		},
		"ETH-USDT": {
			EMA_Period:    266, // 参数确定时间 2023-4-11 20:28:37
			CAP_Period:    5,
			CAP_Max:       "3",
			MaxTradeLever: 5,
		},
	}
}
