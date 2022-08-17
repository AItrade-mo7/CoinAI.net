package analy

import (
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mOKX"
)

func MarketStart() {
	// 清空全局的排行榜单
	okxInfo.Hour8Ticker = []mOKX.AnalySliceType{}
	okxInfo.Hour8TickerUR = []mOKX.AnalySliceType{}

	// 一旦有一个长度不对，则 Market 不合格
	if len(okxInfo.Unit) < 3 || len(okxInfo.TickerList) < 4 || len(okxInfo.AnalyWhole) < 4 || len(okxInfo.AnalySingle) < 4 {
		return
	}

	// 设置 近 8 小时成交量的榜单
	SetHour8Ticker()
	// 判断市场大盘
	JudgeTrendDir()
}
