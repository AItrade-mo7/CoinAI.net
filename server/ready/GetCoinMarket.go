package ready

import (
	"CoinAI.net/server/global"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mFetch"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mStr"
	jsoniter "github.com/json-iterator/go"
)

type TickerResType struct {
	List        []mOKX.TypeTicker                `json:"List"`        // 列表
	AnalyWhole  []mOKX.TypeWholeTickerAnaly      `json:"AnalyWhole"`  // 大盘分析结果
	AnalySingle map[string][]mOKX.AnalySliceType `json:"AnalySingle"` // 单个币种分析结果
	Unit        string                           `json:"Unit"`
	WholeDir    int
}

type ReqCoinMarketType struct {
	Code int           `json:"Code"` // 返回码
	Data TickerResType `json:"Data"` // 返回数据
	Msg  string        `json:"Msg"`  // 描述
}

func GetCoinMarket() {
	resData, err := mFetch.NewHttp(mFetch.HttpOpt{
		Origin: "https://trade.mo7.cc",
		Path:   "/CoinMarket/public/Tickers",
	}).Post()
	if err != nil {
		global.LogErr("ready.GetCoinMarket", err)
		return
	}

	var result ReqCoinMarketType
	jsoniter.Unmarshal(resData, &result)

	if result.Code < 0 {
		global.LogErr("ready.GetCoinMarket", "Err", result.Code, mStr.ToStr(resData))
		return
	}

	// 在这里解析出大盘方向 ，开仓币种 ， 间隔等
	okxInfo.TickerList = result.Data.List

	okxInfo.AnalyWhole = result.Data.AnalyWhole

	okxInfo.AnalySingle = result.Data.AnalySingle

	okxInfo.Unit = result.Data.Unit

	okxInfo.WholeDir = result.Data.WholeDir
}
