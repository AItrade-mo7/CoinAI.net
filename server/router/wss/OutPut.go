package wss

import (
	"CoinAI.net/server/global/config"
	"github.com/EasyGolang/goTools/mTime"
)

type OutPut struct {
	SysTime     int64  `bson:"SysTime"`     // 系统时间
	DataSource  string `bson:"DataSource"`  // 数据来源
	CoinServeID string `bson:"CoinServeID"` //
	UserID      string `bson:"UserID"`      //
}

func GetOutPut() (resData OutPut) {
	resData = OutPut{}
	resData.SysTime = mTime.GetUnixInt64()
	resData.DataSource = "CoinAI.net"
	resData.UserID = config.AppEnv.UserID

	return
}
