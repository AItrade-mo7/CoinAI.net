package ready

import (
	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/okxInfo"
	"CoinAI.net/server/utils/reqDataCenter"
	"github.com/EasyGolang/goTools/mStr"
	jsoniter "github.com/json-iterator/go"
)

func GetSPOTInst() {
	resData, err := reqDataCenter.NewRest(reqDataCenter.RestOpt{
		Origin: "https://trade-api.mo7.cc",
		Path:   "/CoinMarket/public/Inst",
		UserID: config.AppEnv.UserID,
		Data: map[string]any{
			"TypeInst": "SPOT",
		},
		Method: "POST",
	})
	if err != nil {
		global.LogErr("ready.GetSPOTInst", err)
		return
	}

	var result ReqInstType
	jsoniter.Unmarshal(resData, &result)

	if result.Code < 0 {
		global.LogErr("ready.GetSPOTInst", "Err", result.Code, mStr.ToStr(resData))
		return
	}

	okxInfo.SPOT_inst = result.Data
}
