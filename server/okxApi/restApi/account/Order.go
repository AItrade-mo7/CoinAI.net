package account

import (
	"fmt"
	"strings"

	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mStr"
	jsoniter "github.com/json-iterator/go"
)

type OrderParam struct {
	InstID string
	OKXKey mOKX.TypeOkxKey
	OrdId  string
	Side   string // buy  sell
	Sz     string
}

func Order(opt OrderParam) (resErr error) {
	resErr = nil

	if len(opt.InstID) < 3 {
		resErr = fmt.Errorf("account.Order opt.InstID 不能为空 %+v", opt.InstID)
		return
	}
	if len(opt.OKXKey.ApiKey) < 10 {
		resErr = fmt.Errorf("account.Order opt.OKXKey.ApiKey 不能为空 %+v", opt.OKXKey.ApiKey)
		return
	}
	if opt.Side == "buy" || opt.Side == "sell" {
	} else {
		resErr = fmt.Errorf("account.Order opt.Side 不正确 %+v", opt.Side)
		return
	}

	tdMode := "cash" // 币币交易
	ordType := "market"
	find := strings.Contains(opt.InstID, "-SWAP")
	if find {
		tdMode = "cross" // 全仓杠杆
		ordType = "optimal_limit_ioc"
	}

	res, err := mOKX.FetchOKX(mOKX.OptFetchOKX{
		Path:   "/api/v5/trade/order",
		Method: "POST",
		OKXKey: opt.OKXKey,
		Data: map[string]any{
			"instId":  opt.InstID,
			"tdMode":  tdMode,
			"clOrdId": opt.OrdId,
			"side":    opt.Side,
			"posSide": "net",
			"ordType": ordType,
			"sz":      opt.Sz,
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
