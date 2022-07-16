package wss

import "github.com/EasyGolang/goTools/mTime"

type OutPut struct {
	SysTime    int64  `json:"SysTime"`    // 系统时间
	DataSource string `json:"DataSource"` // 数据来源
}

func GetOutPut() (resData OutPut) {
	resData = OutPut{}
	resData.DataSource = "CoinAI.net"
	resData.SysTime = mTime.GetUnixInt64()

	return
}
