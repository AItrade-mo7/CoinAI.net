package okxInfo

import (
	"github.com/EasyGolang/goTools/mOKX"
)

/*  来自 CoinMarket 的数据  */
type AnalyTickerType struct {
	TickerVol   []mOKX.TypeTicker                `bson:"TickerVol"`   // 列表
	AnalyWhole  []mOKX.TypeWholeTickerAnaly      `bson:"AnalyWhole"`  // 大盘分析结果
	AnalySingle map[string][]mOKX.AnalySliceType `bson:"AnalySingle"` // 单个币种分析结果
	Unit        string                           `bson:"Unit"`
	WholeDir    int                              `bson:"WholeDir"`
	DirIndex    int                              `bson:"DirIndex"`
	TimeUnix    int64                            `bson:"TimeUnix"`
	TimeStr     string                           `bson:"TimeStr"`
	TimeID      string                           `bson:"TimeID"`
}

var NowTicker AnalyTickerType

var Inst map[string]mOKX.TypeInst
