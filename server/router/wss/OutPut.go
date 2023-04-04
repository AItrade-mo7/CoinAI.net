package wss

import (
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mTime"
	jsoniter "github.com/json-iterator/go"
)

type OutAnalyTickerType struct {
	Version  int    `bson:"Version"`
	Unit     string `bson:"Unit"`
	TimeUnix int64  `bson:"TimeUnix"`
	TimeStr  string `bson:"TimeStr"`
}
type NowTicker struct {
	Version  int    `bson:"Version"`  // 当前分析版本
	Unit     string `bson:"Unit"`     // 单位
	TimeUnix int64  `bson:"TimeUnix"` // 时间
	TimeStr  string `bson:"TimeStr"`  // 时间字符串
	TimeID   string `bson:"TimeID"`   // TimeID
}

type TradeInstType struct {
	InstID   string `bson:"InstID"`
	QuoteCcy string `bson:"QuoteCcy"`
}

type OutPut struct {
	SysName        string        `bson:"SysName"`
	SysVersion     string        `bson:"SysVersion"`
	Port           string        `bson:"Port"`
	IP             string        `bson:"IP"`
	ServeID        string        `bson:"ServeID"`
	UserID         string        `bson:"UserID"`
	SysTime        int64         `bson:"SysTime"` // 系统时间
	DataSource     string        `bson:"DataSource"`
	MaxApiKeyNum   int           `bson:"MaxApiKeyNum"`
	ApiKeyNum      int           `bson:"ApiKeyNum"`
	Type           string        `bson:"Type"`
	TradeKdataLast mOKX.TypeKd   `bson:"TradeKdataLast"`
	NowTicker      NowTicker     `bson:"Type"`
	TradeInst      TradeInstType `bson:"TradeInst"`
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
	resData.Type = config.AppEnv.Type
	// ApiKey 信息
	resData.ApiKeyNum = len(config.AppEnv.ApiKeyList)
	resData.MaxApiKeyNum = config.AppEnv.MaxApiKeyNum

	// 最后一条数据
	// if len(okxInfo.NowKdataList) > 1 {
	// 	resData.TradeKdataLast = okxInfo.NowKdataList[len(okxInfo.NowKdataList)-1]
	// }

	resData.NowTicker = GetNowTicker()

	resData.TradeInst = GetTradeInst()

	return
}

func GetNowTicker() NowTicker {
	var resData NowTicker
	jsoniter.Unmarshal(mJson.ToJson(okxInfo.NowTicker), &resData)
	return resData
}

func GetTradeInst() TradeInstType {
	var resData TradeInstType
	// jsoniter.Unmarshal(mJson.ToJson(okxInfo.TradeInst), &resData)
	return resData
}
