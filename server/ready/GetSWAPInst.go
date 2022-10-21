package ready

import (
	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/okxInfo"
	"CoinAI.net/server/utils/reqDataCenter"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mStr"
	jsoniter "github.com/json-iterator/go"
)

type ReqInstType struct {
	Code int                      `bson:"Code"` // 返回码
	Data map[string]mOKX.TypeInst `bson:"Data"` // 返回数据
	Msg  string                   `bson:"Msg"`  // 描述
}

func GetSWAPInst() {
	resData, err := reqDataCenter.NewRest(reqDataCenter.RestOpt{
		Origin: "https://trade-api.mo7.cc",
		Path:   "/CoinMarket/public/Inst",
		UserID: config.AppEnv.UserID,
		Data: map[string]any{
			"TypeInst": "SWAP",
		},
		Method: "POST",
	})
	if err != nil {
		global.LogErr("ready.GetSWAPInst", err)
		return
	}

	var result ReqInstType
	jsoniter.Unmarshal(resData, &result)

	if result.Code < 0 {
		global.LogErr("ready.GetSWAPInst", "Err", result.Code, mStr.ToStr(resData))
		return
	}

	okxInfo.SWAP_inst = result.Data
}
