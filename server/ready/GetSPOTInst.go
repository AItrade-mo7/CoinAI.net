package ready

import (
	"CoinAI.net/server/global"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mFetch"
	"github.com/EasyGolang/goTools/mStr"
	jsoniter "github.com/json-iterator/go"
)

func GetSPOTInst() {
	resData, err := mFetch.NewHttp(mFetch.HttpOpt{
		Origin: "https://trade-api.mo7.cc",
		Path:   "/CoinMarket/public/Inst",
		Data: map[string]any{
			"TypeInst": "SPOT",
		},
	}).Post()
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

	okxInfo.SPOT_inst = result.Data
}
