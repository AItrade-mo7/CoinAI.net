package wss

import (
	"CoinAI.net/server/global/config"
	"github.com/EasyGolang/goTools/mTime"
)

type OutPut struct {
	SysTime     int64  `json:"SysTime"`     // 系统时间
	DataSource  string `json:"DataSource"`  // 数据来源
	CoinServeID string `json:"CoinServeID"` //
	UserID      string `json:"UserID"`      //
}

func GetOutPut() (resData OutPut) {
	resData = OutPut{}
	resData.SysTime = mTime.GetUnixInt64()
	resData.DataSource = "CoinAI.net"
	resData.CoinServeID = config.AppEnv.CoinServeID
	resData.UserID = config.AppEnv.UserID

	return
}
