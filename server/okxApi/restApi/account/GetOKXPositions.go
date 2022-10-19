package account

import (
	"fmt"

	"CoinAI.net/server/global/config"
	"github.com/EasyGolang/goTools/mFile"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mStr"
	jsoniter "github.com/json-iterator/go"
)

func GetOKXPositions(ApiKey mOKX.TypeOkxKey) {
	resData, err := mOKX.FetchOKX(mOKX.OptFetchOKX{
		Path:   "/api/v5/account/positions",
		Method: "GET",
		OKXKey: ApiKey,
	})
	var resObj mOKX.TypeReq
	jsoniter.Unmarshal(resData, &resObj)

	Positions_file := mStr.Join(config.Dir.JsonData, "/Positions.json")
	mFile.Write(Positions_file, string(resData))

	fmt.Println(err)
}
