package hunter

import (
	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mFile"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mTalib"
	jsoniter "github.com/json-iterator/go"
)

var (
	EMA_Arr = []string{}
	MA_Arr  = []string{}
	RSI_Arr = []string{}
)

func FormatTradeKdata() {
	if len(okxInfo.NowKdataList) < 200 {
		global.LogErr("hunter.FormatTradeKdata 数据不足")
		return
	}
	// 清理 TradeKdataList
	okxInfo.TradeKdataList = []okxInfo.TradeKdType{}

	EMA_Arr = []string{}
	MA_Arr = []string{}
	RSI_Arr = []string{}

	FormatEnd := []mOKX.TypeKd{}
	for _, Kdata := range okxInfo.NowKdataList {
		FormatEnd = append(FormatEnd, Kdata)
		TradeKdata := NewTradeKdata(FormatEnd, okxInfo.TradeKdataOpt{
			MA_Period:      108,
			RSI_Period:     18,
			RSI_EMA_Period: 14,
			CAP_Period:     3,
		})
		okxInfo.TradeKdataList = append(okxInfo.TradeKdataList, TradeKdata)
	}

	Last := okxInfo.TradeKdataList[len(okxInfo.TradeKdataList)-1]
	LastPrint := map[string]any{
		"InstID":       Last.InstID,
		"TimeStr":      Last.TimeStr,
		"AllLen":       len(okxInfo.TradeKdataList),
		"C":            Last.C,
		"EMA":          Last.EMA,
		"MA":           Last.MA,
		"RSI":          Last.RSI,
		"CAP_EMA":      Last.CAP_EMA,
		"CAP_MA":       Last.CAP_MA,
		"CAPIdx":       Last.CAPIdx,
		"RsiEmaRegion": Last.RsiEmaRegion,
	}
	WriteFilePath := config.Dir.JsonData + "/TradeKdataList.json"
	global.TradeLog.Println("数据整理完毕,写入", len(okxInfo.TradeKdataList), WriteFilePath, mJson.Format(LastPrint))
	mFile.Write(WriteFilePath, string(mJson.ToJson(okxInfo.TradeKdataList)))
}

func NewTradeKdata(TradeKdataList []mOKX.TypeKd, opt okxInfo.TradeKdataOpt) (TradeKdata okxInfo.TradeKdType) {
	TradeKdata = okxInfo.TradeKdType{}
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
