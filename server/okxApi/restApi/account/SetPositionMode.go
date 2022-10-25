package account

import (
	"fmt"

	"CoinAI.net/server/global"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mStr"
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
		resErr = err
		global.LogErr("account.SetPositionMode1", resErr)
		return
	}

	var resObj mOKX.TypeReq
	jsoniter.Unmarshal(res, &resObj)
	if resObj.Code != "0" {
		resErr = fmt.Errorf(mStr.ToStr(resObj.Data))
		global.LogErr("account.SetPositionMode2", resErr)
		return
	}

	return
}
