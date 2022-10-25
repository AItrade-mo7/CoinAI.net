package account

import (
	"fmt"

	"CoinAI.net/server/global"
	"github.com/EasyGolang/goTools/mOKX"
	jsoniter "github.com/json-iterator/go"
)

// 设置持仓模式
func SetPositionMode(OKXKey mOKX.TypeOkxKey) (resErr error) {
	res, err := mOKX.FetchOKX(mOKX.OptFetchOKX{
		Path:   "/api/v5/account/set-position-mode",
		Method: "POST",
		OKXKey: OKXKey,
		Data: map[string]any{
			"posMode": "net_mode",
		},
	})
	if err != nil {
		resErr = fmt.Errorf("account.SetPositionMode1 %+v", err)
		global.LogErr(resErr)
		return
	}

	var resObj mOKX.TypeReq
	jsoniter.Unmarshal(res, &resObj)
	if resObj.Code != "0" {
		resErr = fmt.Errorf("account.SetPositionMode2 %+v", resObj.Data)
		global.LogErr(resErr)

		return
	}

	return
}
