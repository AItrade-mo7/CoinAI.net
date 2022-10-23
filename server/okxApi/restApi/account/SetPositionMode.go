package account

import (
	"fmt"

	"github.com/EasyGolang/goTools/mOKX"
)

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

		return
	}

	fmt.Println(res, err)
	return
}
