package wss

import (
	"CoinAI.net/server/router/result"
	"github.com/EasyGolang/goTools/mRes"
)

func Send() mRes.ResType {
	data := GetOutPut()
	return result.Succeed.WithData(data)
}
