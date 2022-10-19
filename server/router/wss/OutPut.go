package wss

import (
	"CoinAI.net/server/global/config"
	"github.com/EasyGolang/goTools/mTime"
)

type OutPut struct {
	SysTime    int64  `bson:"SysTime"` // 系统时间
	DataSource string `bson:"DataSource"`
	config.AppEnvType
}

func GetOutPut() (resData OutPut) {
	resData = OutPut{}
	resData.SysTime = mTime.GetUnixInt64()
	resData.DataSource = "CoinAI.net"
	resData.Name = config.AppEnv.Name
	resData.Version = config.AppEnv.Version
	resData.Port = config.AppEnv.Port
	resData.IP = config.AppEnv.IP
	resData.ServeID = config.AppEnv.ServeID
	resData.UserID = config.AppEnv.UserID
	resData.ApiKeyList = config.AppEnv.ApiKeyList

	return
}
