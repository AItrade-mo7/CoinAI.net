package analy

import (
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mOKX"
)

func MarketStart() {
	// 一旦有一个长度不对，则 Market 不合格
	if len(okxInfo.Unit) < 3 || len(okxInfo.TickerList) < 4 || len(okxInfo.AnalyWhole) < 4 || len(okxInfo.AnalySingle) < 4 {
		return
	}

	// 将 8 小时切片提取出来，做一个排名
	Hour8Ticker := []mOKX.AnalySliceType{}
	for _, item := range okxInfo.AnalySingle {
		for _, Slice := range item {
			if Slice.DiffHour == 8 {
				Hour8Ticker = append(Hour8Ticker, Slice)
			}
		}
	}
	Hour8TickerVol := mOKX.SortAnalySlice_Volume(Hour8Ticker)
	okxInfo.Hour8Ticker = Hour8TickerVol
}
