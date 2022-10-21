package order

import (
	"fmt"

	"CoinAI.net/server/global"
	"github.com/EasyGolang/goTools/mOKX"
)

func Sell(OkxKey mOKX.TypeOkxKey) (resErr error) {
	resErr = nil
	if len(OkxKey.ApiKey) < 10 {
		resErr = fmt.Errorf(" ApiKey 不能为空")
		global.LogErr(resErr)
		return
	}

	return nil
}
