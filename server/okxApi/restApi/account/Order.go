package account

import (
	"fmt"
	"strings"

	"CoinAI.net/server/global"
	"github.com/EasyGolang/goTools/mOKX"
	jsoniter "github.com/json-iterator/go"
)

type OrderParam struct {
	InstID string
	OKXKey mOKX.TypeOkxKey
	OrdId  string
	Side   string // buy  sell
	Sz     string
}

// 下单接口
func Order(opt OrderParam) (resErr error) {
	resErr = nil

	if len(opt.InstID) < 3 {
		resErr = fmt.Errorf("account.Order opt.InstID 不能为空 %+v", opt.InstID)
		global.LogErr(resErr)
		return
	}
	if len(opt.OKXKey.ApiKey) < 10 {
		resErr = fmt.Errorf("account.Order opt.OKXKey.ApiKey 不能为空 %+v", opt.OKXKey.ApiKey)
		global.LogErr(resErr)
		return
	}
	if opt.Side == "buy" || opt.Side == "sell" {
	} else {
		resErr = fmt.Errorf("account.Order opt.Side 不正确 %+v", opt.Side)
		global.LogErr(resErr)
		return
	}

	tdMode := "cash" // 币币交易
	ordType := "market"
	find := strings.Contains(opt.InstID, "-SWAP")
	if find {
		tdMode = "cross" // 全仓杠杆
		ordType = "optimal_limit_ioc"
	}

	fmt.Println(opt.InstID)
	fmt.Println(opt.Side)
	fmt.Println(opt.Sz)

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
		resErr = fmt.Errorf("account.Order1 %+v", err)
		global.LogErr(resErr)
		return
	}

	var resObj mOKX.TypeReq
	jsoniter.Unmarshal(res, &resObj)
	if resObj.Code != "0" {
		resErr = fmt.Errorf("account.Order2 %+v", resObj.Data)
		global.LogErr(resErr)
		return
	}
	return
}
