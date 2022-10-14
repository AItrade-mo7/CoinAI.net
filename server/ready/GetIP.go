package ready

import (
	"fmt"

	"CoinAI.net/server/utils/reqDataCenter"
	"github.com/EasyGolang/goTools/mJson"
)

func GetIP() {
	resData, err := reqDataCenter.NewRest(reqDataCenter.RestOpt{
		Origin: "https://trade-api.mo7.cc",
		Path:   "/ping",
		Method: "GET",
	})

	fmt.Println(mJson.JsonFormat(resData), err)
}
