package taskPush

import (
	"fmt"

	"CoinAI.net/server/global/config"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mRes"
	jsoniter "github.com/json-iterator/go"
)

type CheckEmailCodeParam struct {
	Email string `bson:"Email"`
	Code  string `bson:"Code"`
}

func CheckEmailCode(opt CheckEmailCodeParam) error {
	resData, err := Request(RequestOpt{
		Origin: config.SysEnv.MessageBaseUrl,
		Path:   "/api/await/CheckEmailCode",
		Data:   mJson.ToJson(opt),
	})
	if err != nil {
		return err
	}

	var resObj mRes.ResType
	jsoniter.Unmarshal(resData, &resObj)

	if resObj.Code < 0 {
		return fmt.Errorf(resObj.Msg)
	}

	return err
}
