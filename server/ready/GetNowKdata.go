package ready

import (
	"CoinAI.net/server/global"
	"github.com/EasyGolang/goTools/mFetch"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mStr"
	jsoniter "github.com/json-iterator/go"
)

type ReqCoinAnalyKdataType struct {
	Code int           `bson:"Code"` // 返回码
	Data []mOKX.TypeKd `bson:"Data"` // 返回数据
	Msg  string        `bson:"Msg"`  // 描述
}

func GetNowKdata(InstID string) (resList []mOKX.TypeKd) {
	resList = []mOKX.TypeKd{}

	resData, err := mFetch.NewHttp(mFetch.HttpOpt{
		Origin: "https://trade-api.mo7.cc",
		Path:   "/CoinMarket/public/GetKdata",
		Data: map[string]any{
			"InstID": InstID,
		},
	}).Post()
	if err != nil {
		global.LogErr("ready.GetCoinAnalyKdata", err)
		return
	}

	var result ReqCoinAnalyKdataType
	jsoniter.Unmarshal(resData, &result)

	if result.Code < 0 {
		global.LogErr("ready.GetCoinAnalyKdata", "Err", result.Code, InstID, mStr.ToStr(resData))
		return
	}

	resList = result.Data
	return
}