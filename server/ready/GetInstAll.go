package ready

import (
	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/utils/reqDataCenter"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mStr"
	jsoniter "github.com/json-iterator/go"
)

type ReqGetInstAllType struct {
	Code int                      `bson:"Code"` // 返回码
	Data map[string]mOKX.TypeInst `bson:"Data"` // 返回数据
	Msg  string                   `bson:"Msg"`  // 描述
}

func GetInstAll() (resList map[string]mOKX.TypeInst) {
	resList = map[string]mOKX.TypeInst{}

	resData, err := reqDataCenter.NewRest(reqDataCenter.RestOpt{
		Origin: "https://trade-api.mo7.cc",
		Path:   "/CoinMarket/public/GetInstAll",
		UserID: config.AppEnv.UserID,
		Method: "Post",
	})
	if err != nil {
		global.LogErr("ready.GetInstAll", err)
		return
	}

	var result ReqGetInstAllType
	jsoniter.Unmarshal(resData, &result)

	if result.Code < 0 {
		global.LogErr("ready.GetInstAll", "Err", result.Code, mStr.ToStr(resData))
		return
	}

	resList = result.Data

	// mFile.Write(config.Dir.JsonData+"/InstAll.json", string(mJson.ToJson(resList)))
	return
}
