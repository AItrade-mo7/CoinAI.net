package account

import (
	"fmt"

	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mStr"
	jsoniter "github.com/json-iterator/go"
)

type GetOrdersPendingParam struct {
	OKXKey mOKX.TypeOkxKey
}

func GetOrdersPending(opt GetOrdersPendingParam) (resErr error) {
	resErr = nil

	if len(opt.OKXKey.ApiKey) < 10 {
		resErr = fmt.Errorf("account.SetLeverage opt.OKXKey.ApiKey 不能为空 %+v", opt.OKXKey.ApiKey)
		return
	}

	res, err := mOKX.FetchOKX(mOKX.OptFetchOKX{
		Path:   "/api/v5/trade/orders-pending",
		Method: "GET",
		OKXKey: opt.OKXKey,
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
