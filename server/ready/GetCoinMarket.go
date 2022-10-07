package ready

import (
	"CoinAI.net/server/global"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mFetch"
	"github.com/EasyGolang/goTools/mStr"
	jsoniter "github.com/json-iterator/go"
)

type ReqCoinMarketType struct {
	Code int                       `bson:"Code"` // 返回码
	Data okxInfo.MarketTickerTable `bson:"Data"` // 返回数据
	Msg  string                    `bson:"Msg"`  // 描述
}

func GetCoinMarket() {
	resData, err := mFetch.NewHttp(mFetch.HttpOpt{
		Origin: "https://trade-api.mo7.cc",
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

	okxInfo.MarketTicker = result.Data
}
