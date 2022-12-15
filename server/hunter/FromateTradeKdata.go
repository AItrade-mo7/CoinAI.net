package hunter

import (
	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mCount"
	"github.com/EasyGolang/goTools/mFile"
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
	okxInfo.TradeKdataList = []okxInfo.TradeKdType{}
	EMA_Arr = []string{}
	MA_Arr = []string{}

	FormatEnd := []mOKX.TypeKd{}

	for _, Kdata := range okxInfo.NowKdataList {
		FormatEnd = append(FormatEnd, Kdata)
		TradeKdata := NewTradeKdata(Kdata, FormatEnd)

		okxInfo.TradeKdataList = append(okxInfo.TradeKdataList, TradeKdata)
	}

	Last := okxInfo.TradeKdataList[len(okxInfo.TradeKdataList)-1]
	LastPrint := map[string]any{
		"AllLen":  len(okxInfo.TradeKdataList),
		"TimeStr": Last.TimeStr,
		"C":       Last.C,
		"InstID":  Last.InstID,
		"EMA_18":  Last.EMA_18,
		"MA_18":   Last.MA_18,
		"RSI_18":  Last.RSI_18,
		"CAP_EMA": Last.CAP_EMA,
		"CAP_MA":  Last.CAP_MA,
	}
	WriteFilePath := config.Dir.JsonData + "/TradeKdataList.json"
	global.RunLog.Println("数据整理完毕,写入", len(okxInfo.TradeKdataList), WriteFilePath, mJson.Format(LastPrint))
	mFile.Write(WriteFilePath, string(mJson.ToJson(okxInfo.TradeKdataList)))
}

func NewTradeKdata(Kdata mOKX.TypeKd, TradeKdataList []mOKX.TypeKd) (TradeKdata okxInfo.TradeKdType) {
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
		Period: 2,
	}).CAP().ToStr()
	TradeKdata.CAP_MA = mTalib.ClistNew(mTalib.ClistOpt{
		CList:  MA_Arr,
		Period: 2,
	}).CAP().ToStr()

	// CAPIdx 计算
	TradeKdata.CAPIdx = GetCAPIdx(TradeKdata)

	// 区域计算

	TradeKdata.RsiRegion = GetRsiRegion(TradeKdata)

	// global.Log.Println("数据整理", mJson.JsonFormat((mJson.ToJson(TradeKdata))))

	return
}

func GetCAPIdx(now okxInfo.TradeKdType) int {
	now_EMA_diff := mCount.Le(now.CAP_EMA, "0") // 1 0 -1  EMA
	now_MA_diff := mCount.Le(now.CAP_MA, "0")   // -1 0 1  MA

	nowDiff := now_EMA_diff
	if now_MA_diff == now_EMA_diff {
		nowDiff = now_MA_diff + now_EMA_diff
	}

	return nowDiff
}

/*
3   大于70
2   60-70
1   50-60
-1  40-50
-2  30-40
-3  小于 30
*/
func GetRsiRegion(now okxInfo.TradeKdType) int {
	RSI := now.RSI_18
	// 1 50-60
	if mCount.Le(RSI, "50") > 0 && mCount.Le(RSI, "60") <= 0 {
		return 1
	}

	// 2 60-70
	if mCount.Le(RSI, "60") > 0 && mCount.Le(RSI, "70") < 0 {
		return 2
	}

	// 3 大于70
	if mCount.Le(RSI, "70") >= 0 {
		return 3
	}

	if mCount.Le(RSI, "50") == 0 {
		return 0
	}

	// -1 40-50
	if mCount.Le(RSI, "40") >= 0 && mCount.Le(RSI, "50") < 0 {
		return -1
	}

	// -2 30-40
	if mCount.Le(RSI, "30") > 0 && mCount.Le(RSI, "40") < 0 {
		return -2
	}

	// -3 30-40
	if mCount.Le(RSI, "30") <= 0 {
		return -3
	}

	return 0
}
