package order

import (
	"fmt"

	"github.com/EasyGolang/goTools/mOKX"
)

func BuySPOT(OkxKey mOKX.TypeOkxKey)error {
	if len(OkxKey.ApiKey) < 10 {
		return fmt.Errorf("Key 不能为空")
	}

	return nil
}
