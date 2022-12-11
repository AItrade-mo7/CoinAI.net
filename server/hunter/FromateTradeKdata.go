package hunter

import (
	"CoinAI.net/server/global"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mTalib"
	jsoniter "github.com/json-iterator/go"
)

var (
	EMA_Arr = []string{}
	MA_Arr  = []string{}
)

func FormatTradeKdata() {
	if len(okxInfo.NowKdataList) < 100 {
		global.LogErr("hunter.FormatTradeKdata 数据不足")
		return
	}
	// 执行一次清理
	TradeKdataList = []TradeKdType{}
	EMA_Arr = []string{}
	MA_Arr = []string{}

	FormatEnd := []mOKX.TypeKd{}

	for _, Kdata := range okxInfo.NowKdataList {
		FormatEnd = append(FormatEnd, Kdata)
		TradeKdata := NewTradeKdata(Kdata, FormatEnd)

		TradeKdataList = append(TradeKdataList, TradeKdata)
	}

	Last := TradeKdataList[len(TradeKdataList)-1]
	LastPrint := map[string]any{
		"AllLen":  len(TradeKdataList),
		"TimeStr": Last.TimeStr,
		"C":       Last.C,
		"InstID":  Last.InstID,
		"EMA_18":  Last.EMA_18,
		"MA_18":   Last.MA_18,
		"RSI_18":  Last.RSI_18,
		"CAP_EMA": Last.CAP_EMA,
		"CAP_MA":  Last.CAP_MA,
	}
	global.RunLog.Println("数据整理完毕", mJson.Format(LastPrint))
}

func NewTradeKdata(Kdata mOKX.TypeKd, TradeKdataList []mOKX.TypeKd) (TradeKdata TradeKdType) {
	jsonByte := mJson.ToJson(Kdata)
	jsoniter.Unmarshal(jsonByte, &TradeKdata)

	// EMA
	TradeKdata.EMA_18 = mTalib.ClistNew(mTalib.ClistOpt{
		KDList: TradeKdataList,
		Period: 18,
	}).EMA().ToStr()
	EMA_Arr = append(EMA_Arr, TradeKdata.EMA_18)

	// MA
	TradeKdata.MA_18 = mTalib.ClistNew(mTalib.ClistOpt{
		KDList: TradeKdataList,
		Period: 18,
	}).MA().ToStr()
	MA_Arr = append(MA_Arr, TradeKdata.MA_18)

	// RSI_18
	TradeKdata.RSI_18 = mTalib.ClistNew(mTalib.ClistOpt{
		KDList: TradeKdataList,
		Period: 18,
	}).RSI().ToStr()

	// CAP
	TradeKdata.CAP_EMA = mTalib.ClistNew(mTalib.ClistOpt{
		CList:  EMA_Arr,
		Period: 3,
	}).CAP().ToStr()
	TradeKdata.CAP_MA = mTalib.ClistNew(mTalib.ClistOpt{
		CList:  MA_Arr,
		Period: 3,
	}).CAP().ToStr()

	global.Log.Println("数据整理", EMA_Arr[len(EMA_Arr)-1], MA_Arr[len(MA_Arr)-1], mJson.JsonFormat((mJson.ToJson(TradeKdata))))

	return
}
