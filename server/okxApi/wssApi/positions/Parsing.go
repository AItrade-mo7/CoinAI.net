package positions

import (
	"fmt"

	"github.com/EasyGolang/goTools/mCount"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/EasyGolang/goTools/mTime"
	jsoniter "github.com/json-iterator/go"
)

type OkxWssRes struct {
	Op   string       `json:"op"`
	Args []LoginWrite `json:"args"`
}

type LoginWrite struct {
	APIKey     string `json:"apiKey"`
	Passphrase string `json:"passphrase"`
	Timestamp  string `json:"timestamp"`
	Sign       string `json:"sign"`
}

func Write_LoginInfo(cont any) (resData LoginWrite, resErr error) {
	var result OkxWssRes
	resErr = jsoniter.Unmarshal(cont.([]byte), &result)

	loginInfo := result.Args[0]
	if len(loginInfo.APIKey) < 20 || len(loginInfo.Sign) < 20 {
		resErr = fmt.Errorf("positions.Write_LoginInfo 长度不足")
		return
	}
	loginInfo.APIKey = mStr.Join(
		loginInfo.APIKey[0:4], "****", loginInfo.APIKey[10:14],
	)
	loginInfo.Sign = mStr.Join(
		loginInfo.Sign[0:4], "****", loginInfo.Sign[10:14],
	)
	loginInfo.Passphrase = mStr.Join(
		"******",
	)

	loginInfo.Timestamp = mTime.UnixFormat(mCount.Mul(loginInfo.Timestamp, "1000"))

	resData = loginInfo
	return
}

type LoginRes struct {
	Event string `json:"event"`
	Msg   string `json:"msg"`
	Code  string `json:"code"`
}

func Read_Login(cont []byte) (resData bool) {
	var result LoginRes
	err := jsoniter.Unmarshal(cont, &result)
	if err != nil {
		resData = false
	}

	if result.Code != "0" {
		return false
	}

	return true
}
