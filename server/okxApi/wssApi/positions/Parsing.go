package positions

import (
	"fmt"

	"github.com/EasyGolang/goTools/mCount"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/EasyGolang/goTools/mTime"
	jsoniter "github.com/json-iterator/go"
)

type OkxWssRes struct {
	Op   string       `bson:"op"`
	Args []LoginWrite `bson:"args"`
}

type LoginWrite struct {
	APIKey     string `bson:"apiKey"`
	Passphrase string `bson:"passphrase"`
	Timestamp  string `bson:"timestamp"`
	Sign       string `bson:"sign"`
}

func Write_LoginInfo(cont any) (resData LoginWrite, resErr error) {
	var result OkxWssRes
	resErr = jsoniter.Unmarshal(cont.([]byte), &result)

	loginInfo := result.Args[0]
	if len(loginInfo.APIKey) < 20 || len(loginInfo.Sign) < 20 {
		resErr = fmt.Errorf("positions.Write_LoginInfo 长度不足")
		return
	}
	loginInfo.APIKey = mStr.Fuzzy(loginInfo.APIKey)
	loginInfo.Sign = mStr.Fuzzy(loginInfo.Sign)
	loginInfo.Passphrase = "******"

	loginInfo.Timestamp = mTime.UnixFormat(mCount.Mul(loginInfo.Timestamp, "1000"))

	resData = loginInfo
	return
}

type LoginRes struct {
	Event string `bson:"event"`
	Msg   string `bson:"msg"`
	Code  string `bson:"code"`
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
