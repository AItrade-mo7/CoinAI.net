package okxInfo

import (
	"github.com/EasyGolang/goTools/mOKX"
)

/*  来自 CoinMarket 的数据  */

type AnalyTickerType struct {
	TickerVol   []mOKX.TypeTicker                `bson:"TickerVol"`   // 列表
	AnalyWhole  []mOKX.TypeWholeTickerAnaly      `bson:"AnalyWhole"`  // 大盘分析结果
	AnalySingle map[string][]mOKX.AnalySliceType `bson:"AnalySingle"` // 单个币种分析结果
	MillionCoin []mOKX.AnalySliceType            `bson:"MillionCoin"` // 过亿的币种盘子
	Version     int                              `bson:"Version"`     // 当前分析版本
	Unit        string                           `bson:"Unit"`        // 单位
	TimeUnix    int64                            `bson:"TimeUnix"`    // 时间
	TimeStr     string                           `bson:"TimeStr"`     // 时间字符串
	TimeID      string                           `bson:"TimeID"`      // TimeID
}

var NowTicker AnalyTickerType

var Inst map[string]mOKX.TypeInst
