package account

import (
	"fmt"

	"CoinAI.net/server/global"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mStr"
	jsoniter "github.com/json-iterator/go"
)

type CancelOrderParam struct {
	OKXKey mOKX.TypeOkxKey
	InstID string
	OrdId  string
}

// 取消订单
func CancelOrder(opt CancelOrderParam) (resErr error) {
	resErr = nil

	if len(opt.InstID) < 3 {
		resErr = fmt.Errorf("account.CancelOrder opt.InstID 不能为空 %+v", opt.InstID)
		global.LogErr(resErr)
		return
	}
	if len(opt.OKXKey.ApiKey) < 10 {
		resErr = fmt.Errorf("account.CancelOrder opt.OKXKey.ApiKey 不能为空 %+v", opt.OKXKey.ApiKey)
		global.LogErr(resErr)
		return
	}
	if len(opt.OrdId) < 3 {
		resErr = fmt.Errorf("account.CancelOrder opt.OrdId 不能为空 %+v", opt.OrdId)
		global.LogErr(resErr)
		return
	}

	res, err := mOKX.FetchOKX(mOKX.OptFetchOKX{
		Path:   "/api/v5/trade/cancel-order",
		Method: "POST",
		OKXKey: opt.OKXKey,
		Data: map[string]any{
			"instId": opt.InstID,
			"ordId":  opt.OrdId,
		},
	})
	if err != nil {
		resErr = err
		global.LogErr(resErr)
		return
	}

	var resObj mOKX.TypeReq
	jsoniter.Unmarshal(res, &resObj)
	if resObj.Code != "0" {
		resErr = fmt.Errorf(mStr.ToStr(resObj.Data))
		global.LogErr(resErr)
		return
	}

	return
}
