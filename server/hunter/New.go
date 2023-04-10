package hunter

import (
	"os"

	"CoinAI.net/server/global/config"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mPath"
	"github.com/EasyGolang/goTools/mStr"
)

type HunterOpt struct {
	HunterName      string // 默认:MyHunter
	InstID          string // 当前策略交易对 默认:BTC-USDT
	Describe        string // Hunter描述  默认:空
	OutPutDirectory string // 数据读写的目录 默认 jsonData
}

type HunterObj struct {
	HunterName         string // 策略的名字
	Describe           string // 描述
	InstID             string // 当前策略主打币种
	MaxLen             int
	TradeInst          mOKX.TypeInst                 // 交易的 InstID SWAP
	KdataInst          mOKX.TypeInst                 // K线的 InstID SPOT
	NowKdataList       []mOKX.TypeKd                 // 现货的原始K线
	TradeKdataList     []okxInfo.TradeKdType         // 交易K线-计算好各种指标之后的K线
	TradeKdataOpt      okxInfo.TradeKdataOpt         // 计算交易指标的参数
	NowVirtualPosition okxInfo.VirtualPositionType   // 当前的虚拟持仓
	PositionArr        []okxInfo.VirtualPositionType // 当前持仓列表
	OrderArr           []okxInfo.VirtualPositionType // 平仓列表
	OutPutDirectory    string                        // 数据读写的目录
}

/*
新建一个 Hunter 策略

需要: HunterOpt
需要热加载策略
*/
func New(opt HunterOpt) *HunterObj {
	obj := HunterObj{}
	obj.HunterName = opt.HunterName
	obj.Describe = opt.Describe
	obj.InstID = opt.InstID
	obj.MaxLen = 900
	obj.TradeInst = mOKX.TypeInst{}
	obj.KdataInst = mOKX.TypeInst{}

	obj.NowKdataList = []mOKX.TypeKd{}
	obj.TradeKdataList = []okxInfo.TradeKdType{}

	obj.TradeKdataOpt = okxInfo.TradeKdataOpt{}
	obj.NowVirtualPosition = okxInfo.VirtualPositionType{}

	obj.OutPutDirectory = opt.OutPutDirectory

	if len(obj.HunterName) < 1 {
		obj.HunterName = "MyHunter"
	}

	if len(obj.InstID) < 1 {
		obj.InstID = "BTC-USDT"
	}

	if len(obj.Describe) < 1 {
		obj.Describe = "暂无描述"
	}

	if len(obj.OutPutDirectory) < 1 {
		obj.OutPutDirectory = mStr.Join(config.Dir.JsonData, "/", obj.HunterName)
	}

	// 默认目录在 jsonData 下
	isOutPutDirectoryPath := mPath.Exists(obj.OutPutDirectory)
	if !isOutPutDirectoryPath {
		// 不存在则创建 logs 目录
		os.MkdirAll(obj.OutPutDirectory, 0o777)
	}

	return &obj
}

func (_this *HunterObj) ReadVirtualPosition() {
	// 在这里读取虚拟持仓并设置初始值
}
