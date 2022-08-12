package analy

import (
	"fmt"

	"CoinAI.net/server/okxInfo"
)

func MarketStart() {
	// 一旦有一个长度不对，则 Market 不合格
	if len(okxInfo.Unit) < 3 || len(okxInfo.TickerList) < 4 || len(okxInfo.AnalyWhole) < 4 || len(okxInfo.AnalySingle) < 4 {
		return
	}

	for _, item := range okxInfo.TickerList {
		fmt.Println("TickerList", item.InstID)
	}

	for _, item := range okxInfo.AnalyWhole {
		fmt.Println("AnalyWhole", item.DiffHour, item.DirIndex)
	}

	for key, item := range okxInfo.AnalySingle {
		fmt.Println("AnalySingle", key)
		for _, Slice := range item {
			fmt.Println("Slice", Slice.InstID, Slice.StartTime, Slice.EndTime)
		}
	}
}
