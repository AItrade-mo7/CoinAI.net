package wss

import (
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mTime"
	jsoniter "github.com/json-iterator/go"
)

type OutPut struct {
	SysName        string                    `bson:"SysName"`
	SysVersion     string                    `bson:"SysVersion"`
	Port           string                    `bson:"Port"`
	IP             string                    `bson:"IP"`
	ServeID        string                    `bson:"ServeID"`
	UserID         string                    `bson:"UserID"`
	SysTime        int64                     `bson:"SysTime"` // 系统时间
	DataSource     string                    `bson:"DataSource"`
	MaxApiKeyNum   int                       `bson:"MaxApiKeyNum"`
	ApiKeyNum      int                       `bson:"ApiKeyNum"`
	TradeKdataLast mOKX.TypeKd               `bson:"TradeKdataLast"`
	NowTicker      NowTicker                 `bson:"Type"`
	HunterData     map[string]HunterDataType `bson:"Hunter"`
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
	resData.ApiKeyNum = len(config.AppEnv.ApiKeyList)
	resData.MaxApiKeyNum = config.AppEnv.MaxApiKeyNum
	resData.NowTicker = GetNowTicker()
	// Hunter 信息
	resData.HunterData = GetHunterData()

	return
}

// ======================

type NowTicker struct {
	Version  int    `bson:"Version"`  // 当前分析版本
	Unit     string `bson:"Unit"`     // 单位
	TimeUnix int64  `bson:"TimeUnix"` // 时间
	TimeStr  string `bson:"TimeStr"`  // 时间字符串
	TimeID   string `bson:"TimeID"`   // TimeID
}

func GetNowTicker() NowTicker {
	var resData NowTicker
	jsoniter.Unmarshal(mJson.ToJson(okxInfo.NowTicker), &resData)
	return resData
}

// ======================

type NowKdataType struct {
	InstID   string `bson:"InstID"`   // 持仓币种
	TimeUnix int64  `bson:"TimeUnix"` // 毫秒时间戳
	C        string `bson:"C"`        // 收盘价格
	Dir      int    `bson:"Dir"`      // 方向 (收盘-开盘) ，1：涨 & -1：跌 & 0：横盘
}

type HunterDataType struct {
	HunterName    string
	HLPerLevel    int // 震荡等级
	MaxLen        int
	TradeInstID   string       // 交易的 InstID SWAP
	KdataInstID   string       // K线的 InstID SPOT
	NowKdata      NowKdataType // 现货的原始K线
	TradeKdataOpt okxInfo.TradeKdataOpt
}

func GetHunterData() map[string]HunterDataType {
	HunterData := make(map[string]HunterDataType)

	for key, item := range okxInfo.NowHunterData {
		var newData HunterDataType
		newData.HunterName = item.HunterName
		newData.HLPerLevel = item.HLPerLevel
		newData.MaxLen = item.MaxLen
		newData.TradeInstID = item.TradeInst.InstID
		newData.KdataInstID = item.KdataInst.InstID

		var newKdata NowKdataType
		lastKdata := item.NowKdataList[len(item.NowKdataList)-1]
		jsoniter.Unmarshal(mJson.ToJson(lastKdata), &newKdata)

		newData.NowKdata = newKdata

		newData.TradeKdataOpt = item.TradeKdataOpt

		HunterData[key] = newData
	}

	return HunterData
}
