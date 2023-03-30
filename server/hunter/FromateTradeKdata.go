package hunter

import (
	"CoinAI.net/server/global"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mTalib"
	jsoniter "github.com/json-iterator/go"
)

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

var TradeKdataList []TradeKdType

type TradeKdataOpt struct {
	MA_Period      int // 108
	RSI_Period     int // 18
	RSI_EMA_Period int // 14
	CAP_Period     int // 3
}

var (
	EMA_Arr = []string{}
	MA_Arr  = []string{}
	RSI_Arr = []string{}
)

func FormatTradeKdata(TradeKdataOpt TradeKdataOpt) {
	if len(okxInfo.NowKdataList) < TradeKdataOpt.MA_Period {
		global.LogErr("hunter.FormatTradeKdata 数据不足")
		return
	}

	// 清理 TradeKdataList
	TradeKdataList = []TradeKdType{}

	EMA_Arr = []string{}
	MA_Arr = []string{}
	RSI_Arr = []string{}

	FormatEnd := []mOKX.TypeKd{}
	for _, Kdata := range okxInfo.NowKdataList {
		FormatEnd = append(FormatEnd, Kdata)
		TradeKdata := NewTradeKdata(FormatEnd, TradeKdataOpt)
		TradeKdataList = append(TradeKdataList, TradeKdata)
	}
}

func NewTradeKdata(TradeKdataList []mOKX.TypeKd, opt TradeKdataOpt) (TradeKdata TradeKdType) {
	TradeKdata = TradeKdType{}
	jsonByte := mJson.ToJson(TradeKdataList[len(TradeKdataList)-1])
	jsoniter.Unmarshal(jsonByte, &TradeKdata)

	TradeKdata.Opt = opt

	if TradeKdata.Opt.MA_Period == 0 ||
		TradeKdata.Opt.RSI_Period == 0 ||
		TradeKdata.Opt.RSI_EMA_Period == 0 ||
		TradeKdata.Opt.CAP_Period == 0 {
		return
	}

	// EMA
	TradeKdata.EMA = mTalib.ClistNew(mTalib.ClistOpt{
		KDList: TradeKdataList,
		Period: opt.MA_Period,
	}).EMA().ToStr()
	EMA_Arr = append(EMA_Arr, TradeKdata.EMA)

	// MA
	TradeKdata.MA = mTalib.ClistNew(mTalib.ClistOpt{
		KDList: TradeKdataList,
		Period: opt.MA_Period,
	}).MA().ToStr()
	MA_Arr = append(MA_Arr, TradeKdata.MA)

	// RSI
	TradeKdata.RSI = mTalib.ClistNew(mTalib.ClistOpt{
		KDList: TradeKdataList,
		Period: opt.RSI_Period,
	}).RSI().ToStr()
	RSI_Arr = append(RSI_Arr, TradeKdata.RSI)

	// RSI_EMA
	TradeKdata.RSI_EMA = mTalib.ClistNew(mTalib.ClistOpt{
		CList:  RSI_Arr,
		Period: opt.RSI_EMA_Period,
	}).EMA().ToStr()

	// CAP_EMA
	TradeKdata.CAP_EMA = mTalib.ClistNew(mTalib.ClistOpt{
		CList:  EMA_Arr,
		Period: opt.CAP_Period,
	}).CAP().ToStr()
	// CAP_MA
	TradeKdata.CAP_MA = mTalib.ClistNew(mTalib.ClistOpt{
		CList:  MA_Arr,
		Period: opt.CAP_Period,
	}).CAP().ToStr()

	// CAPIdx 计算
	TradeKdata.CAPIdx = GetCAPIdx(TradeKdata)

	// 区域计算
	TradeKdata.RsiEmaRegion = GetRsiRegion(TradeKdata)

	return
}
