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

	// 清空全局的排行榜单
	okxInfo.Hour8Ticker = []mOKX.AnalySliceType{}
	okxInfo.Hour8TickerUR = []mOKX.AnalySliceType{}

	// 将 8 小时切片提取出来，做一个排名
	Hour8Ticker := []mOKX.AnalySliceType{}
	for _, item := range okxInfo.AnalySingle {
		for _, Slice := range item {
			if Slice.DiffHour == 8 {
				Hour8Ticker = append(Hour8Ticker, Slice)
			}
		}
	}
	// 按照成交量排序
	Hour8TickerVolSort := mOKX.SortAnalySlice_Volume(Hour8Ticker)
	// 前 1/2
	MaxLen := len(Hour8TickerVolSort) / 2
	Hour8TickerVol := Hour8TickerVolSort[0:MaxLen]

	// copy 一份
	Filter8Ticker := make([]mOKX.AnalySliceType, len(Hour8TickerVol))
	copy(Filter8Ticker, Hour8TickerVol)

	// 成交量排序 前 1/2
	okxInfo.Hour8Ticker = Filter8Ticker

	// 成交量排序前 1/2  的涨跌幅排序
	okxInfo.Hour8TickerUR = mOKX.SortAnalySlice_UR(Filter8Ticker)
}
