package wss

import (
	"CoinAI.net/server/global/config"
	"github.com/EasyGolang/goTools/mTime"
)

type OutPut struct {
	SysTime    int64  `bson:"SysTime"`    // 系统时间
	DataSource string `bson:"DataSource"` // 数据来源
	ServeID    string `bson:"ServeID"`    //
	UserID     string `bson:"UserID"`     //
}

func GetOutPut() (resData OutPut) {
	resData = OutPut{}
	resData.SysTime = mTime.GetUnixInt64()
	resData.DataSource = "CoinAI.net"
	resData.UserID = config.AppEnv.UserID
	resData.ServeID = config.AppEnv.ServeID

	return
}
