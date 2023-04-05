package hunter

import (
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mOKX"
)

type HunterOpt struct {
	HunterName    string // 默认 MyHunter
	HLPerLevel    int    // 币种的震荡等级 2
	MaxLen        int    // 900
	TradeKdataOpt okxInfo.TradeKdataOpt
}

type HunterObj struct {
	HunterName     string // 策略的名字
	HLPerLevel     int    // 震荡等级
	MaxLen         int
	TradeInst      mOKX.TypeInst         // 交易的 InstID SWAP
	KdataInst      mOKX.TypeInst         // K线的 InstID SPOT
	NowKdataList   []mOKX.TypeKd         // 现货的原始K线
	TradeKdataList []okxInfo.TradeKdType // 计算好各种指标之后的K线
	TradeKdataOpt  okxInfo.TradeKdataOpt
}

func New(opt HunterOpt) *HunterObj {
	obj := HunterObj{}
	obj.TradeInst = mOKX.TypeInst{}
	obj.KdataInst = mOKX.TypeInst{}

	obj.HunterName = opt.HunterName
	if len(obj.HunterName) < 1 {
		obj.HunterName = "MyHunter"
	}

	obj.HLPerLevel = opt.HLPerLevel
	if (obj.HLPerLevel) < 1 {
		obj.HLPerLevel = 2
	}

	obj.MaxLen = opt.MaxLen
	if (obj.MaxLen) < 900 {
		obj.MaxLen = 900
	}

	obj.TradeKdataOpt = opt.TradeKdataOpt

	if obj.TradeKdataOpt.MA_Period < 0 {
		obj.TradeKdataOpt.MA_Period = 171
	}
	if obj.TradeKdataOpt.CAP_Period < 0 {
		obj.TradeKdataOpt.CAP_Period = 4
	}

	return &obj
}
