package account

import (
	"fmt"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/dbType"
	"github.com/EasyGolang/goTools/mCount"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mStr"
	jsoniter "github.com/json-iterator/go"
)

type ClosePositionParam struct {
	TradeInst mOKX.TypeInst // 交易币种信息
	OKXKey    dbType.OkxKeyType
}

// 平仓接口
func ClosePosition(opt ClosePositionParam) (resErr error) {
	resErr = nil

	if len(opt.TradeInst.InstID) < 3 {
		resErr = fmt.Errorf("account.Order opt.InstID 不能为空 %+v Name:%+v", opt.TradeInst.InstID, opt.OKXKey.Name)
		global.LogErr(resErr)
		return
	}
	if mCount.Le(opt.TradeInst.MinSz, "0") < 1 {
		resErr = fmt.Errorf("account.Order opt.TradeInst.MinSz 不能为空 %+v Name:%+v", opt.TradeInst.MinSz, opt.OKXKey.Name)
		global.LogErr(resErr)
		return
	}

	if len(opt.OKXKey.ApiKey) < 10 {
		resErr = fmt.Errorf("account.Order opt.OKXKey.ApiKey 不能为空 %+v Name:%+v", opt.OKXKey.ApiKey, opt.OKXKey.Name)
		global.LogErr(resErr)
		return
	}

	Data := map[string]any{
		"instId":  opt.TradeInst.InstID,
		"mgnMode": "cross",
	}
	res, err := mOKX.FetchOKX(mOKX.OptFetchOKX{
		Path:   "/api/v5/trade/close-position",
		Method: "POST",
		OKXKey: mOKX.TypeOkxKey{
			ApiKey:     opt.OKXKey.ApiKey,
			SecretKey:  opt.OKXKey.SecretKey,
			Passphrase: opt.OKXKey.Passphrase,
		},
		Data: Data,
	})
	if err != nil {
		resErr = fmt.Errorf("account.ClosePosition1 %+v Name:%+v", mStr.ToStr(err), opt.OKXKey.Name)
		global.LogErr(resErr)
		return
	}

	var resObj mOKX.TypeReq
	jsoniter.Unmarshal(res, &resObj)
	if resObj.Code != "0" {
		resErr = fmt.Errorf("account.ClosePosition2 %+v Name:%+v", mStr.ToStr(res), opt.OKXKey.Name)
		global.LogErr(resErr)
		return
	}

	return
}
