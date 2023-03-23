package ready

import (
	"CoinAI.net/server/global"
	"CoinAI.net/server/okxInfo"
	"CoinAI.net/server/utils/taskPush"
	"github.com/EasyGolang/goTools/mStr"
	jsoniter "github.com/json-iterator/go"
)

type ReqCoinMarketType struct {
	Code int                     `bson:"Code"` // 返回码
	Data okxInfo.AnalyTickerType `bson:"Data"` // 返回数据
	Msg  string                  `bson:"Msg"`  // 描述
}

func GetNowTickerAnaly() (returnData okxInfo.AnalyTickerType) {
	returnData = okxInfo.AnalyTickerType{}

	resData, err := taskPush.Request(taskPush.RequestOpt{
		Origin: "https://trade-api.mo7.cc",
		Path:   "/CoinMarket/public/GetNowTickerAnaly",
	})
	if err != nil {
		global.LogErr("ready.GetNowTickerAnaly", err)
		return
	}

	var result ReqCoinMarketType
	jsoniter.Unmarshal(resData, &result)

	if result.Code < 0 {
		global.LogErr("ready.GetNowTickerAnaly", "Err", result.Code, mStr.ToStr(resData))
		return
	}

	returnData = result.Data
	return
}
