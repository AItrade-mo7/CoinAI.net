package ready

import (
	"CoinAI.net/server/okxInfo"
)

type ReqCoinMarketType struct {
	Code int                     `bson:"Code"` // 返回码
	Data okxInfo.AnalyTickerType `bson:"Data"` // 返回数据
	Msg  string                  `bson:"Msg"`  // 描述
}

func GetNowTickerAnaly() (resData okxInfo.AnalyTickerType) {
	resData = okxInfo.AnalyTickerType{}

	// res, err := reqDataCenter.NewRest(reqDataCenter.RestOpt{
	// 	Origin: "https://trade-api.mo7.cc",
	// 	Path:   "/CoinMarket/public/GetNowTickerAnaly",
	// 	UserID: config.AppEnv.UserID,
	// 	Method: "Post",
	// })
	// if err != nil {
	// 	global.LogErr("ready.GetCoinMarket", err)
	// 	return
	// }

	// var result ReqCoinMarketType
	// jsoniter.Unmarshal(res, &result)

	// if result.Code < 0 {
	// 	global.LogErr("ready.GetCoinMarket", "Err", result.Code, mStr.ToStr(resData))
	// 	return
	// }

	// resData = result.Data
	return
}
