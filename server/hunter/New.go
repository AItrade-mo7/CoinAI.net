package hunter

import (
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mOKX"
)

type HunterOpt struct {
	HunterName string // 默认 MyHunter
	InstID     string // 当前策略交易对
	MaxLen     int    // 900
}

type HunterObj struct {
	HunterName     string // 策略的名字
	MaxLen         int
	TradeInst      mOKX.TypeInst         // 交易的 InstID SWAP
	KdataInst      mOKX.TypeInst         // K线的 InstID SPOT
	NowKdataList   []mOKX.TypeKd         // 现货的原始K线
	TradeKdataList []okxInfo.TradeKdType // 计算好各种指标之后的K线
	TradeKdataOpt  okxInfo.TradeKdataOpt
	MaxTradeLever  int // 最优秀的交易杠杆数
}

func New(opt HunterOpt) *HunterObj {
	obj := HunterObj{}
	obj.TradeInst = mOKX.TypeInst{}
	obj.KdataInst = mOKX.TypeInst{}

	obj.HunterName = opt.HunterName
	if len(obj.HunterName) < 1 {
		obj.HunterName = "MyHunter"
	}

	obj.MaxLen = opt.MaxLen
	if (obj.MaxLen) < 900 {
		obj.MaxLen = 900
	}

	return &obj
}
