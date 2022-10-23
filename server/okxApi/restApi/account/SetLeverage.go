package account

import (
	"fmt"

	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mStr"
	jsoniter "github.com/json-iterator/go"
)

type SetLeverageParam struct {
	InstID string
	OKXKey mOKX.TypeOkxKey
	Lever  int
}

func SetLeverage(opt SetLeverageParam) (resErr error) {
	resErr = nil

	if len(opt.InstID) < 3 {
		resErr = fmt.Errorf("account.SetLeverage opt.InstID 不能为空 %+v", opt.InstID)
		return
	}
	if opt.Lever < 1 {
		resErr = fmt.Errorf("account.SetLeverage opt.Lever 不能为0 %+v", opt.Lever)
		return
	}
	if len(opt.OKXKey.ApiKey) < 10 {
		resErr = fmt.Errorf("account.SetLeverage opt.OKXKey.ApiKey 不能为空 %+v", opt.OKXKey.ApiKey)
		return
	}

	res, err := mOKX.FetchOKX(mOKX.OptFetchOKX{
		Path:   "/api/v5/account/set-leverage",
		Method: "POST",
		OKXKey: opt.OKXKey,
		Data: map[string]any{
			"instId":  opt.InstID,
			"lever":   mStr.ToStr(opt.Lever),
			"mgnMode": "cross",
		},
	})
	if err != nil {
		resErr = err
		return
	}

	var resObj mOKX.TypeReq
	jsoniter.Unmarshal(res, &resObj)
	if resObj.Code != "0" {
		resErr = fmt.Errorf(mStr.ToStr(resObj.Data))
		return
	}

	return
}
