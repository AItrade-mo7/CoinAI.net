package ready

import (
	"CoinAI.net/server/global"
	"github.com/EasyGolang/goTools/mFetch"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mStr"
	jsoniter "github.com/json-iterator/go"
)

type ReqCoinAnalyKdataType struct {
	Code int           `json:"Code"` // 返回码
	Data []mOKX.TypeKd `json:"Data"` // 返回数据
	Msg  string        `json:"Msg"`  // 描述
}

func GetCoinAnalyKdata(InstID string) (resList []mOKX.TypeKd) {
	resList = []mOKX.TypeKd{}

	resData, err := mFetch.NewHttp(mFetch.HttpOpt{
		Origin: "https://trade.mo7.cc",
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
		global.LogErr("ready.GetCoinMarket", "Err", result.Code, mStr.ToStr(resData))
		return
	}

	resList = result.Data
	return
}
