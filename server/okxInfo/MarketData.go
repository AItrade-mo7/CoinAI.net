package okxInfo

import (
	"github.com/EasyGolang/goTools/mOKX"
)

/*  来自 CoinMarket 的数据  */
type MarketTickerTable struct {
	List           []mOKX.TypeTicker                `bson:"List"`        // 成交量排序列表
	ListU_R24      []mOKX.TypeTicker                `bson:"ListU_R24"`   // 涨跌幅排序列表
	AnalyWhole     []mOKX.TypeWholeTickerAnaly      `bson:"AnalyWhole"`  // 大盘分析结果
	AnalySingle    map[string][]mOKX.AnalySliceType `bson:"AnalySingle"` // 单个币种分析结果
	Unit           string                           `bson:"Unit"`
	WholeDir       int                              `bson:"WholeDir"`
	TimeUnix       int64                            `bson:"TimeUnix"`
	Time           string                           `bson:"Time"`
	CreateTimeUnix int64                            `bson:"CreateTimeUnix"`
	CreateTime     string                           `bson:"CreateTime"`
}

var MarketTicker MarketTickerTable

// 按照时间倒序
var AnalyList []MarketTickerTable

/* 币种历史数据 */
// 按照时间正序
var AnalyKdata_SPOT map[string][]mOKX.TypeKd
var AnalyKdata_SWAP map[string][]mOKX.TypeKd
