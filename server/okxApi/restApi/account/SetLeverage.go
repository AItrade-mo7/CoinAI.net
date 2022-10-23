package account

import (
	"fmt"

	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mStr"
	jsoniter "github.com/json-iterator/go"
)

// type SetLeverageParam struct {
// 	 InstID
// }

func SetLeverage(OKXKey mOKX.TypeOkxKey) (resErr error) {
	res, err := mOKX.FetchOKX(mOKX.OptFetchOKX{
		Path:   "/api/v5/account/set-leverage",
		Method: "POST",
		OKXKey: OKXKey,
		Data: map[string]any{
			"instId": "net_mode",
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
