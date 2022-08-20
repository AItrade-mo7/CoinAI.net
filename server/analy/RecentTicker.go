package analy

import (
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mOKX"
)

func RecentTicker() (resList []mOKX.AnalySliceType) {
	resList = []mOKX.AnalySliceType{}

	// 一旦有一个长度不对，则 Market 不合格
	if len(okxInfo.Unit) < 3 || len(okxInfo.TickerList) < 4 || len(okxInfo.AnalyWhole) < 4 || len(okxInfo.AnalySingle) < 4 {
		return
	}

	// 将 12 小时切片提取出来，做一个排名
	Hour12Ticker := []mOKX.AnalySliceType{}
	for _, item := range okxInfo.AnalySingle {
		for _, Slice := range item {
			if Slice.DiffHour == 12 {
				Hour12Ticker = append(Hour12Ticker, Slice)
			}
		}
	}
	// 按照成交量排序
	Hour12TickerVolSort := mOKX.SortAnalySlice_Volume(Hour12Ticker)
	// 前 1/2
	MaxLen := len(Hour12TickerVolSort) - (len(Hour12TickerVolSort) / 3)
	Hour12TickerVol := Hour12TickerVolSort[0:MaxLen]

	// copy 一份
	Filter12Ticker := make([]mOKX.AnalySliceType, len(Hour12TickerVol))
	copy(Filter12Ticker, Hour12TickerVol)

	// 最近 12 小时 成交量排序 前 1/2
	resList = Filter12Ticker
	return
}
