package hunter

import "github.com/EasyGolang/goTools/mOKX"

type HunterOpt struct {
	HunterName string
	HLPerLevel int // 币种的震荡等级 2
	MaxLen     int // 900
}

type HunterObj struct {
	HunterName   string // 策略的名字
	HLPerLevel   int    // 震荡等级
	MaxLen       int
	TradeInst    mOKX.TypeInst // 交易的 InstID
	KdataInst    mOKX.TypeInst // 交易的 InstID
	NowKdataList []mOKX.TypeKd
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

	return &obj
}
