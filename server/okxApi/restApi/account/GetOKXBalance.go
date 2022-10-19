package account

import (
	"fmt"

	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mOKX"
	jsoniter "github.com/json-iterator/go"
)

func GetOKXBalance(ApiKey mOKX.TypeOkxKey) {
	resData, err := mOKX.FetchOKX(mOKX.OptFetchOKX{
		Path:   "/api/v5/account/balance",
		Method: "GET",
		OKXKey: ApiKey,
	})

	var resObj mOKX.TypeReq
	jsoniter.Unmarshal(resData, &resObj)

	mJson.Println(resObj)

	fmt.Println( err)
}
