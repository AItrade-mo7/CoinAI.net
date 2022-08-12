package okxInfo

import (
	"github.com/EasyGolang/goTools/mOKX"
)

var TickerList []mOKX.TypeTicker

var AnalyWhole []mOKX.TypeWholeTickerAnaly // 大盘分析结果

var AnalySingle map[string][]mOKX.AnalySliceType // 单个币种分析合集

// 计价的锚定货币
var Unit string

var IsMarket = false