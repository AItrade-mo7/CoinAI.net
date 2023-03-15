package wss

import (
	"CoinAI.net/server/global/config"
	"github.com/EasyGolang/goTools/mTime"
)

type OutAnalyTickerType struct {
	Version  int    `bson:"Version"`
	Unit     string `bson:"Unit"`
	TimeUnix int64  `bson:"TimeUnix"`
	TimeStr  string `bson:"TimeStr"`
}

type TradeCoinType struct{}

type OutPut struct {
	SysName      string `bson:"SysName"`
	SysVersion   string `bson:"SysVersion"`
	Port         string `bson:"Port"`
	IP           string `bson:"IP"`
	ServeID      string `bson:"ServeID"`
	UserID       string `bson:"UserID"`
	SysTime      int64  `bson:"SysTime"` // 系统时间
	DataSource   string `bson:"DataSource"`
	MaxApiKeyNum int    `bson:"MaxApiKeyNum"`
	ApiKeyNum    int    `bson:"ApiKeyNum"`
}

func GetOutPut() (resData OutPut) {
	resData = OutPut{}
	// 系统运行信息
	resData.SysName = config.AppEnv.SysName
	resData.SysVersion = config.AppEnv.SysVersion
	resData.Port = config.AppEnv.Port
	resData.IP = config.AppEnv.IP
	resData.ServeID = config.AppEnv.ServeID
	resData.UserID = config.AppEnv.UserID
	resData.SysTime = mTime.GetUnixInt64()
	resData.DataSource = config.SysName

	// ApiKey 信息
	resData.MaxApiKeyNum = config.AppEnv.MaxApiKeyNum
	resData.ApiKeyNum = len(config.AppEnv.ApiKeyList)

	return
}
