package account

import (
	"fmt"
	"go.uber.org/zap"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/dbType"
	"github.com/EasyGolang/goTools/mCount"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mStr"
	jsoniter "github.com/json-iterator/go"
)

type OrderParam struct {
	TradeInst mOKX.TypeInst // 交易币种信息
	OKXKey    dbType.OkxKeyType
	OrdId     string
	Side      string // buy  sell
	Sz        string
	MaxEvent  func()
}

// 下单接口
func Order(opt OrderParam) (resErr error) {
	resErr = nil

	if len(opt.TradeInst.InstID) < 3 {
		resErr = fmt.Errorf("account.Order opt.InstID 不能为空 %+v Name:%+v", opt.TradeInst.InstID, opt.OKXKey.Name)
		global.LogErr("该错误已同步至用户邮箱", resErr)
		LogErr(opt.OKXKey, resErr)
		return
	}
	if mCount.Le(opt.TradeInst.MinSz, "0") < 1 {
		resErr = fmt.Errorf("account.Order opt.TradeInst.MinSz 不能为空 %+v Name:%+v", opt.TradeInst.MinSz, opt.OKXKey.Name)
		global.LogErr("该错误已同步至用户邮箱", resErr)
		LogErr(opt.OKXKey, resErr)
		return
	}

	if len(opt.OKXKey.ApiKey) < 10 {
		resErr = fmt.Errorf("account.Order opt.OKXKey.ApiKey 不能为空 %+v Name:%+v", opt.OKXKey.ApiKey, opt.OKXKey.Name)
		global.LogErr("该错误已同步至用户邮箱", resErr)
		LogErr(opt.OKXKey, resErr)
		return
	}

	if opt.MaxEvent == nil {
		opt.MaxEvent = func() {}
	}

	if opt.Side == "buy" || opt.Side == "sell" {
	} else {
		resErr = fmt.Errorf("account.Order opt.Side 不正确 %+v Name:%+v", opt.Side, opt.OKXKey.Name)
		global.LogErr("该错误已同步至用户邮箱", resErr)
		LogErr(opt.OKXKey, resErr)
		return
	}

	tdMode := "cash" // 币币交易
	ordType := "market"
	if opt.TradeInst.InstType == "SWAP" {
		tdMode = "cross" // 全仓杠杆
		ordType = "optimal_limit_ioc"
	}

	// 默认减去一个最小开仓单位，也就是每次开仓至少为 两个最小开仓单位
	opt.Sz = mCount.Sub(opt.Sz, opt.TradeInst.MinSz)
	if mCount.Le(opt.Sz, opt.TradeInst.MinSz) < 0 {
		resErr = fmt.Errorf("交易数量太小 %+v Name:%+v", opt.Sz, opt.OKXKey.Name)
		global.LogErr("该错误已同步至用户邮箱", resErr)
		LogErr(opt.OKXKey, resErr)
		return
	}

	// 大于最大数量 则 最大数乘以 0.8
	if mCount.Le(opt.Sz, opt.TradeInst.MaxMktSz) > 0 {
		opt.Sz = mCount.Mul(opt.TradeInst.MaxMktSz, "0.8")
		resErr = fmt.Errorf("交易数量超出限制 %+v Name:%+v", opt.Sz, opt.OKXKey.Name)
		global.LogErr("该错误已同步至用户邮箱", resErr)
		LogErr(opt.OKXKey, resErr)
	}

	opt.Sz = mCount.Cent(opt.Sz, 0)

	Data := map[string]any{
		"instId":  opt.TradeInst.InstID,
		"tdMode":  tdMode,
		"clOrdId": opt.OrdId,
		"side":    opt.Side,
		"posSide": "net",
		"ordType": ordType,
		"sz":      opt.Sz,
	}
	res, err := mOKX.FetchOKX(mOKX.OptFetchOKX{
		Path:   "/api/v5/trade/order",
		Method: "POST",
		OKXKey: mOKX.TypeOkxKey{
			ApiKey:     opt.OKXKey.ApiKey,
			SecretKey:  opt.OKXKey.SecretKey,
			Passphrase: opt.OKXKey.Passphrase,
		},
		Data: Data,
	})
	// 打印接口日志
	global.OKXLogo.Info("account.Order",
		zap.Error(err),
		zap.String("res", mStr.ToStr(res)),
		zap.String(opt.OKXKey.Name, mJson.ToStr(Data)),
	)

	if err != nil {
		resErr = fmt.Errorf("account.Order1 %+v Name:%+v", mStr.ToStr(err), opt.OKXKey.Name)
		global.LogErr("该错误已同步至用户邮箱", resErr)
		LogErr(opt.OKXKey, resErr)
		return
	}

	var resObj mOKX.TypeReq
	jsoniter.Unmarshal(res, &resObj)
	if resObj.Code != "0" {
		resErr = fmt.Errorf("account.Order2 Data:%+v ; %+v Name:%+v", mStr.ToStr(Data), mStr.ToStr(res), opt.OKXKey.Name)
		global.LogErr("该错误已同步至用户邮箱", resErr)
		LogErr(opt.OKXKey, resErr)
		return
	}

	return
}
