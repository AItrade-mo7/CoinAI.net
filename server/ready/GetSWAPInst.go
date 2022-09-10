package ready

import (
	"CoinAI.net/server/global"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mFetch"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mStr"
	jsoniter "github.com/json-iterator/go"
)

type ReqSWAPInstType struct {
	Code int                      `json:"Code"` // 返回码
	Data map[string]mOKX.TypeInst `json:"Data"` // 返回数据
	Msg  string                   `json:"Msg"`  // 描述
}

func GetSWAPInst() {
	resData, err := mFetch.NewHttp(mFetch.HttpOpt{
		Origin: "https://trade.mo7.cc",
		Path:   "/CoinMarket/public/Inst",
		Data: map[string]any{
			"TypeInst": "SWAP",
		},
	}).Post()
	if err != nil {
		global.LogErr("ready.GetSWAPInst", err)
		return
	}

	var result ReqSWAPInstType
	jsoniter.Unmarshal(resData, &result)

	if result.Code < 0 {
		global.LogErr("ready.GetSWAPInst", "Err", result.Code, mStr.ToStr(resData))
		return
	}

	okxInfo.SWAP_inst = result.Data
}
