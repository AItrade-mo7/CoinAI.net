package account

import (
	"fmt"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/dbType"
	"github.com/EasyGolang/goTools/mOKX"
	jsoniter "github.com/json-iterator/go"
)

type CancelOrderParam struct {
	OKXKey dbType.OkxKeyType
	Order  PendingOrderType
}

// 取消订单
func CancelOrder(opt CancelOrderParam) (resErr error) {
	resErr = nil

	if len(opt.OKXKey.ApiKey) < 10 {
		resErr = fmt.Errorf("account.CancelOrder opt.OKXKey.ApiKey 不能为空 Name:%+v", opt.OKXKey.Name)
		global.LogErr(resErr)
		return
	}
	if len(opt.Order.InstID) < 3 {
		resErr = fmt.Errorf("account.CancelOrder opt.Orders.InstID 不能为空:%+v Name:%+v", opt.Order.InstID, opt.OKXKey.Name)
		global.LogErr(resErr)
		return
	}

	res, err := mOKX.FetchOKX(mOKX.OptFetchOKX{
		Path:   "/api/v5/trade/cancel-order",
		Method: "POST",
		OKXKey: mOKX.TypeOkxKey{
			ApiKey:     opt.OKXKey.ApiKey,
			SecretKey:  opt.OKXKey.SecretKey,
			Passphrase: opt.OKXKey.Passphrase,
		},
		Data: map[string]any{
			"instId": opt.Order.InstID,
			"ordId":  opt.Order.OrdID,
		},
	})
	if err != nil {
		resErr = fmt.Errorf("account.CancelOrder1 %+v %+v Name:%+v", err, opt.OKXKey, opt.OKXKey.Name)
		global.LogErr(resErr)
		return
	}

	var resObj mOKX.TypeReq
	jsoniter.Unmarshal(res, &resObj)
	if resObj.Code != "0" {
		resErr = fmt.Errorf("account.CancelOrder2 %s %+v Name:%+v", res, opt.OKXKey, opt.OKXKey.Name)
		global.LogErr(resErr)
		return
	}

	return
}
