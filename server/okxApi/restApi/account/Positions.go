package account

import (
	"fmt"

	"github.com/EasyGolang/goTools/mOKX"
)

func GetOKXPositions(ApiKey mOKX.TypeOkxKey) {
	resData, err := mOKX.FetchOKX(mOKX.OptFetchOKX{
		Path:   "/api/v5/account/positions",
		Method: "GET",
		OKXKey: ApiKey,
	})

	fmt.Println(resData, err)
}
