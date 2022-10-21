package ready

import (
	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/utils/reqDataCenter"
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

	resData, err := reqDataCenter.NewRest(reqDataCenter.RestOpt{
		Origin: "https://trade-api.mo7.cc",
		Path:   "/CoinMarket/public/GetNowKdata",
		UserID: config.AppEnv.UserID,
		Method: "Post",
		Data: map[string]any{
			"InstID": InstID,
		},
	})
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
