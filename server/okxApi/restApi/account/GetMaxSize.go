package account

import (
	"fmt"

	"CoinAI.net/server/global"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mStr"
	jsoniter "github.com/json-iterator/go"
)

type GetMaxSizeParam struct {
	InstID string
	OKXKey mOKX.TypeOkxKey
}

type MaxSizeType struct {
	Ccy     string
	InstID  string
	MaxBuy  string
	MaxSell string
}

// 获得最大可开仓数量
func GetMaxSize(opt GetMaxSizeParam) (resData MaxSizeType, resErr error) {
	resErr = nil
	resData = MaxSizeType{}

	if len(opt.InstID) < 3 {
		resErr = fmt.Errorf("account.GetMaxSize opt.InstID 不能为空 %+v", opt.InstID)
		global.LogErr(resErr)
		return
	}
	if len(opt.OKXKey.ApiKey) < 10 {
		resErr = fmt.Errorf("account.GetMaxSize opt.OKXKey.ApiKey 不能为空 %+v", opt.OKXKey.ApiKey)
		global.LogErr(resErr)
		return
	}

	res, err := mOKX.FetchOKX(mOKX.OptFetchOKX{
		Path:   "/api/v5/account/max-size",
		Method: "GET",
		OKXKey: opt.OKXKey,
		Data: map[string]any{
			"instId": opt.InstID,
			"tdMode": "cross",
		},
	})
	if err != nil {
		resErr = err
		global.LogErr("account.GetMaxSize1", resErr)
		return
	}

	var resObj mOKX.TypeReq
	jsoniter.Unmarshal(res, &resObj)
	if resObj.Code != "0" {
		resErr = fmt.Errorf(mStr.ToStr(resObj.Data))
		global.LogErr("account.GetMaxSize2", resErr)
		return
	}

	var result []MaxSizeType
	jsoniter.Unmarshal(mJson.ToJson(resObj.Data), &result)
	if len(result) > 0 {
		resData = result[0]
		global.LogErr("account.GetMaxSize3", resErr)
		return
	}

	return
}
