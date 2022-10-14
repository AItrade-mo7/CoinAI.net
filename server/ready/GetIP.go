package ready

import (
	"fmt"

	"github.com/EasyGolang/goTools/mFetch"
	"github.com/EasyGolang/goTools/mJson"
)

func GetIP() {
	resData, err := mFetch.NewHttp(mFetch.HttpOpt{
		Origin: "https://trade-api.mo7.cc",
		Path:   "/ping",
	}).Get()

	fmt.Println(mJson.JsonFormat(resData), err)
}
