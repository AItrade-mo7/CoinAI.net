package wss

import (
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/okxInfo"
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
	Name         string             `bson:"Name"`
	Version      string             `bson:"Version"`
	Port         string             `bson:"Port"`
	IP           string             `bson:"IP"`
	ServeID      string             `bson:"ServeID"`
	UserID       string             `bson:"UserID"`
	SysTime      int64              `bson:"SysTime"` // 系统时间
	DataSource   string             `bson:"DataSource"`
	TradeInstID  string             `bson:"TradeInst"`
	TradeLever   int                `bson:"TradeLever"`
	NowTicker    OutAnalyTickerType `bson:"NowTicker"`
	MaxApiKeyNum int                `bson:"MaxApiKeyNum"`
	ApiKeyNum    int                `bson:"ApiKeyNum"`
}

func GetOutPut() (resData OutPut) {
	resData = OutPut{}
	// 系统运行信息
	resData.Name = config.AppEnv.Name
	resData.Version = config.AppEnv.Version
	resData.Port = config.AppEnv.Port
	resData.IP = config.AppEnv.IP
	resData.ServeID = config.AppEnv.ServeID
	resData.UserID = config.AppEnv.UserID
	// 系统时间
	resData.SysTime = mTime.GetUnixInt64()
	resData.DataSource = "CoinAI.net"
	// 下单信息
	if okxInfo.IsSPOT {
		resData.TradeInstID = okxInfo.TradeInst.SPOT.InstID
	} else {
		resData.TradeInstID = okxInfo.TradeInst.SWAP.InstID
	}
	resData.TradeLever = okxInfo.TradeLever

	// Ticker 信息
	resData.NowTicker.Unit = okxInfo.NowTicker.Unit
	resData.NowTicker.TimeUnix = okxInfo.NowTicker.TimeUnix
	resData.NowTicker.TimeStr = okxInfo.NowTicker.TimeStr

	// ApiKey 信息
	resData.MaxApiKeyNum = okxInfo.MaxApiKeyNum
	resData.ApiKeyNum = len(config.AppEnv.ApiKeyList)

	return
}
