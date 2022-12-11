package hunter

import (
	"CoinAI.net/server/global"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mTalib"
	jsoniter "github.com/json-iterator/go"
)

func FormatTradeKdata() {
	if len(okxInfo.NowKdataList) < 100 {
		global.LogErr("hunter.FormatTradeKdata 数据不足")
		return
	}
	TradeKdataList = []TradeKdType{} // 执行一次清理

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

	TradeKdata.EMA_18 = mTalib.ClistNew(mTalib.ClistOpt{
		KDList: TradeKdataList, // 数据
		Period: 18,             // 周期
	}).EMA().ToStr()

	global.Log.Println("数据整理", mJson.JsonFormat((mJson.ToJson(TradeKdata))))

	return
}
