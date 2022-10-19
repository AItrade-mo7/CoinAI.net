package wss

import (
	"github.com/EasyGolang/goTools/mTime"
)

type OutPut struct {
	SysTime int64 `bson:"SysTime"` // 系统时间
}

func GetOutPut() (resData OutPut) {
	resData = OutPut{}
	resData.SysTime = mTime.GetUnixInt64()
	// resData.DataSource = "CoinAI.net"
	// resData.UserID = config.AppEnv.UserID
	// resData.ServeID = config.AppEnv.ServeID

	return
}
