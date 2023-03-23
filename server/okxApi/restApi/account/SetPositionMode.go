package account

import (
	"fmt"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/dbType"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mStr"
	jsoniter "github.com/json-iterator/go"
)

// 设置持仓模式
func SetPositionMode(OKXKey dbType.OkxKeyType) (resErr error) {
	res, err := mOKX.FetchOKX(mOKX.OptFetchOKX{
		Path:   "/api/v5/account/set-position-mode",
		Method: "POST",
		OKXKey: mOKX.TypeOkxKey{
			ApiKey:     OKXKey.ApiKey,
			SecretKey:  OKXKey.SecretKey,
			Passphrase: OKXKey.Passphrase,
		},
		Data: map[string]any{
			"posMode": "net_mode",
		},
	})
	if err != nil {
		resErr = fmt.Errorf("account.SetPositionMode1 %+v Name:%+v", mStr.ToStr(err), OKXKey.Name)
		global.LogErr(resErr)
		return
	}

	var resObj mOKX.TypeReq
	jsoniter.Unmarshal(res, &resObj)
	if resObj.Code != "0" {
		resErr = fmt.Errorf("account.SetPositionMode2 %s Name:%+v", mStr.ToStr(res), OKXKey.Name)
		global.LogErr(resErr)

		return
	}

	return
}
