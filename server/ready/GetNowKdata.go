package ready

import (
	"CoinAI.net/server/global"
	"CoinAI.net/server/utils/taskPush"
	"github.com/EasyGolang/goTools/mJson"
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

	Data := map[string]any{
		"InstID": InstID,
	}

	resData, err := taskPush.Request(taskPush.RequestOpt{
		Origin: "https://trade-api.mo7.cc",
		Path:   "/CoinMarket/public/GetNowKdata",
		Data:   mJson.ToJson(Data),
	})
	if err != nil {
		global.LogErr("ready.GetNowKdata", err)
		return
	}

	var result ReqCoinAnalyKdataType
	jsoniter.Unmarshal(resData, &result)

	if result.Code < 0 {
		global.LogErr("ready.GetNowKdata", "Err", result.Code, InstID, mStr.ToStr(resData))
		return
	}

	resList = result.Data
	return
}
