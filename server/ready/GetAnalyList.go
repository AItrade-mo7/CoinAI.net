package ready

import (
	"CoinAI.net/server/global"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mFetch"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mStr"
	jsoniter "github.com/json-iterator/go"
)

type PagingType struct {
	List    []any `bson:"List"`
	Total   int64 `bson:"Total"`
	Current int64 `bson:"Current"`
	Size    int64 `bson:"Size"`
}

type ReqGetAnalyListType struct {
	Code int        `bson:"Code"` // 返回码
	Data PagingType `bson:"Data"` // 返回数据
	Msg  string     `bson:"Msg"`  // 描述
}

func GetAnalyList() (resList []mOKX.TypeKd) {
	resList = []mOKX.TypeKd{}

	resData, err := mFetch.NewHttp(mFetch.HttpOpt{
		Origin: "https://trade-api.mo7.cc",
		Path:   "/CoinMarket/public/GetAnalyList",
		Data: map[string]any{
			"Size":    300,
			"Current": 0,
			"Sort": map[string]any{
				"CreateTimeUnix": -1,
			},
			"Type": "Serve",
		},
	}).Post()
	if err != nil {
		global.LogErr("ready.GetAnalyList", err)
		return
	}

	var result ReqGetAnalyListType
	jsoniter.Unmarshal(resData, &result)

	if result.Code < 0 {
		global.LogErr("ready.GetAnalyList", "Err", result.Code, mStr.ToStr(resData))
		return
	}

	var AnalyList []okxInfo.MarketTickerTable
	jsoniter.Unmarshal(mJson.ToJson(result.Data.List), &AnalyList)

	okxInfo.AnalyList = AnalyList
	return
}
