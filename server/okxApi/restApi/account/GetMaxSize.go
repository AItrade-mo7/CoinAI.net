package account

import (
	"fmt"
	"strings"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/global/dbType"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mStr"
	jsoniter "github.com/json-iterator/go"
)

type GetMaxSizeParam struct {
	InstID string
	OKXKey dbType.OkxKeyType
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
		resErr = fmt.Errorf("account.GetMaxSize opt.InstID 不能为空 %+v Name:%+v", opt.InstID, opt.OKXKey.Name)
		global.LogErr(resErr)
		return
	}
	if len(opt.OKXKey.ApiKey) < 10 {
		resErr = fmt.Errorf("account.GetMaxSize opt.OKXKey.ApiKey 不能为空 %+v Name:%+v", opt.OKXKey.ApiKey, opt.OKXKey.Name)
		global.LogErr(resErr)
		return
	}

	tdMode := "cash" // 币币交易
	find := strings.Contains(opt.InstID, config.SWAP_suffix)
	if find {
		tdMode = "cross" // 全仓杠杆
	}

	Data := map[string]any{
		"instId": opt.InstID,
		"tdMode": tdMode,
	}
	res, err := mOKX.FetchOKX(mOKX.OptFetchOKX{
		Path:   "/api/v5/account/max-size",
		Method: "GET",
		OKXKey: mOKX.TypeOkxKey{
			ApiKey:     opt.OKXKey.ApiKey,
			SecretKey:  opt.OKXKey.SecretKey,
			Passphrase: opt.OKXKey.Passphrase,
		},
		Data: Data,
	})
	// 打印接口日志
	global.OKXLogo.Println("account.GetMaxSize",
		err,
		mStr.ToStr(res),
		opt.OKXKey.Name,
		mJson.ToStr(Data),
	)

	if err != nil {
		resErr = fmt.Errorf("account.GetMaxSize1 %+v %+v Name:%+v", mStr.ToStr(err), opt.OKXKey, opt.OKXKey.Name)
		global.LogErr(resErr)
		return
	}

	var resObj mOKX.TypeReq
	jsoniter.Unmarshal(res, &resObj)
	if resObj.Code != "0" {
		resErr = fmt.Errorf("account.GetMaxSize1 %s %+v Name:%+v", mStr.ToStr(res), opt.OKXKey, opt.OKXKey.Name)
		global.LogErr(resErr)
		return
	}

	var result []MaxSizeType
	jsoniter.Unmarshal(mJson.ToJson(resObj.Data), &result)
	if len(result) > 0 {
		resData = result[0]
		return
	}

	return
}
