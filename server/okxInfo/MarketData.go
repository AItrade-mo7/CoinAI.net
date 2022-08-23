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

// 币种历史数据
var AnalyKdata map[string][]mOKX.TypeKd

// 币种挑选策略
type AnalySelectType struct {
	InstID     string
	MaxHLPer   string
	HLPerAvg   string
	DiffMaxAvg string
}

var AnalySelect []AnalySelectType
