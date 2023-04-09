package hunter

import (
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mOKX"
)

type HunterOpt struct {
	HunterName string // 默认 MyHunter
	InstID     string // 当前策略交易对
	MaxLen     int    // 900
	Describe   string // 描述
}

type HunterObj struct {
	HunterName         string // 策略的名字
	Describe           string // 描述
	InstID             string // 当前策略主打币种
	MaxLen             int
	TradeInst          mOKX.TypeInst         // 交易的 InstID SWAP
	KdataInst          mOKX.TypeInst         // K线的 InstID SPOT
	NowKdataList       []mOKX.TypeKd         // 现货的原始K线
	TradeKdataList     []okxInfo.TradeKdType // 计算好各种指标之后的K线
	TradeKdataOpt      okxInfo.TradeKdataOpt
	NowVirtualPosition okxInfo.VirtualPositionType // 当前虚拟持仓
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
	if obj.MaxLen < 900 {
		obj.MaxLen = 900
	}

	obj.InstID = opt.InstID

	if len(obj.InstID) < 1 {
		obj.InstID = "BTC-USDT"
	}

	obj.Describe = opt.Describe

	return &obj
}
