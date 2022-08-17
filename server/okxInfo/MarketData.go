package okxInfo

import (
	"github.com/EasyGolang/goTools/mOKX"
)

/*  原始的市场数据  */
var TickerList []mOKX.TypeTicker

var AnalyWhole []mOKX.TypeWholeTickerAnaly // 大盘分析结果

var AnalySingle map[string][]mOKX.AnalySliceType // 单个币种分析合集

// 计价的锚定货币
var Unit string

var WholeDir int

/*  分析结果  */
// 最近 8 小时的市场
var Hour8Ticker []mOKX.AnalySliceType
var Hour8TickerUR []mOKX.AnalySliceType
