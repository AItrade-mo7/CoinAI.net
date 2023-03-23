package account

import (
	"fmt"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/global/dbType"
	"github.com/EasyGolang/goTools/mFile"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mStr"
	jsoniter "github.com/json-iterator/go"
)

type GetOrdersPendingParam struct {
	OKXKey dbType.OkxKeyType
}

type PendingOrderType struct {
	AccFillSz       string
	AvgPx           string
	CTime           string
	Category        string
	Ccy             string
	ClOrdID         string
	Fee             string
	FeeCcy          string
	FillPx          string
	FillSz          string
	FillTime        string
	InstID          string
	InstType        string
	Lever           string
	OrdID           string
	OrdType         string
	Pnl             string
	PosSide         string
	Px              string
	Rebate          string
	RebateCcy       string
	ReduceOnly      string
	Side            string
	SlOrdPx         string
	SlTriggerPx     string
	SlTriggerPxType string
	Source          string
	State           string
	Sz              string
	Tag             string
	TdMode          string
	TgtCcy          string
	TpOrdPx         string
	TpTriggerPx     string
	TpTriggerPxType string
	TradeID         string
	UTime           string
}

// 未成交订单信息
func GetOrdersPending(opt GetOrdersPendingParam) (resData []PendingOrderType, resErr error) {
	resErr = nil

	if len(opt.OKXKey.ApiKey) < 10 {
		resErr = fmt.Errorf("account.GetOrdersPending opt.OKXKey.ApiKey 不能为空 Name:%+v", opt.OKXKey.Name)
		global.LogErr(resErr)
		return
	}

	res, err := mOKX.FetchOKX(mOKX.OptFetchOKX{
		Path:   "/api/v5/trade/orders-pending",
		Method: "GET",
		OKXKey: mOKX.TypeOkxKey{
			ApiKey:     opt.OKXKey.ApiKey,
			SecretKey:  opt.OKXKey.SecretKey,
			Passphrase: opt.OKXKey.Passphrase,
		},
	})
	// 打印接口日志
	global.OKXLogo.Println("account.GetOrdersPending",
		err,
		mStr.ToStr(res),
		opt.OKXKey.Name,
	)

	mFile.Write(config.Dir.JsonData+"/OrdersPending.json", string(res))
	if err != nil {
		resErr = fmt.Errorf("account.GetOrdersPending1 %+v Name:%+v", mStr.ToStr(err), opt.OKXKey.Name)
		global.LogErr(resErr)
		return
	}

	var resObj mOKX.TypeReq
	jsoniter.Unmarshal(res, &resObj)
	if resObj.Code != "0" {
		resErr = fmt.Errorf("account.GetOrdersPending2 %s Name:%+v", mStr.ToStr(res), opt.OKXKey.Name)
		global.LogErr(resErr)
		return
	}

	var result []PendingOrderType
	jsoniter.Unmarshal(mJson.ToJson(resObj.Data), &result)

	resData = result

	return
}
