package wss

import (
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/global/dbType"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mCount"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mTime"
	jsoniter "github.com/json-iterator/go"
)

type OutPut struct {
	SysName      string                    `bson:"SysName"`
	SysVersion   string                    `bson:"SysVersion"`
	Port         string                    `bson:"Port"`
	IP           string                    `bson:"IP"`
	ServeID      string                    `bson:"ServeID"`
	UserID       string                    `bson:"UserID"`
	SysTime      int64                     `bson:"SysTime"` // 系统时间
	DataSource   string                    `bson:"DataSource"`
	MaxApiKeyNum int                       `bson:"MaxApiKeyNum"`
	ApiKeyNum    int                       `bson:"ApiKeyNum"`
	NowTicker    NowTicker                 `bson:"Type"`
	HunterData   map[string]HunterDataType `bson:"Hunter"`
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
	TimeUnix  int64    `bson:"TimeUnix"` // 时间
	TimeID    string   `bson:"TimeID"`   // TimeID
	TickerVol []string `bson:"TimeID"`
}

func GetNowTicker() NowTicker {
	resData := NowTicker{}
	TickerVol := []string{}
	for _, item := range okxInfo.NowTicker.TickerVol {
		TickerVol = append(TickerVol, item.InstID)
	}
	resData.TickerVol = TickerVol
	resData.TimeUnix = okxInfo.NowTicker.TimeUnix
	resData.TimeID = okxInfo.NowTicker.TimeID

	return resData
}

// ======================

type NowKdataType struct {
	InstID   string `bson:"InstID"`   // 持仓币种
	TimeUnix int64  `bson:"TimeUnix"` // 时间
	C        string `bson:"C"`        // 收盘价格
	Dir      int    `bson:"Dir"`      // 方向 (收盘-开盘) ，1：涨 & -1：跌 & 0：横盘
}

type PositionType struct {
	InstID     string `bson:"InstID"`     // 下单币种 | 运行中设置
	HunterName string `bson:"HunterName"` // 策略名称 | 运行中设置
	NowTimeStr string `bson:"NowTimeStr"` // 当前K线时间 | 运行中设置
	NowTime    int64  `bson:"NowTime"`    // 当前实际时间戳 | 运行中设置
	NowC       string `bson:"NowC"`       // 当前收盘价 | 运行中设置
	CAP_EMA    string `bson:"CAP_EMA"`    // 当前的 CAP 值 | 运行中设置
	EMA        string `bson:"EMA"`        // 当前的 EMA 值 | 运行中设置
	// 下单时设置
	OpenAvgPx   string `bson:"OpenAvgPx"`   // 开仓价格 | 下单时设置
	OpenTimeStr string `bson:"OpenTimeStr"` // 开仓K线时间 | 下单时设置
	OpenTime    int64  `bson:"OpenTime"`    // 开仓实际时间戳  | 下单时设置
	NowDir      int    `bson:"NowDir"`      // 当前持仓状态 没持仓0  持多仓 1  持空仓 -1 | 初始化设置为0，下单时更新
	// 通过原始数据计算得出
	InitMoney   string `bson:"InitMoney"`   // 初始金钱 | 固定值初始值设置
	ChargeUpl   string `bson:"ChargeUpl"`   // 当前手续费率 | 固定值初始值设置
	NowUplRatio string `bson:"NowUplRatio"` // 当前未实现收益率(计算得出) | 运行中设置
	Money       string `bson:"Money"`       // 账户当前余额 | 如果没有初始值设置一次，下单时计算
	MakeMoney   string `bson:"MakeMoney"`   // 本单盈利
}

type HunterDataType struct {
	HunterName         string
	Describe           string       // 描述
	TradeInstID        string       // 交易的 InstID SWAP
	KdataInstID        string       // K线的 InstID SPOT
	NowKdata           NowKdataType // 现货的原始K线
	KdataLen           int
	TradeKdataLen      int
	TradeKdataOpt      dbType.TradeKdataOpt
	NowVirtualPosition PositionType
}

func GetHunterData() map[string]HunterDataType {
	HunterData := make(map[string]HunterDataType)

	for key, item := range okxInfo.NowHunterData {
		var newData HunterDataType
		newData.HunterName = item.HunterName
		newData.KdataLen = len(item.NowKdataList)
		newData.TradeKdataLen = len(item.TradeKdataList)
		newData.TradeInstID = item.TradeInst.InstID
		newData.KdataInstID = item.KdataInst.InstID
		newData.Describe = item.Describe
		var newKdata NowKdataType
		if len(item.NowKdataList) > 1 {
			lastKdata := item.NowKdataList[len(item.NowKdataList)-1]
			jsoniter.Unmarshal(mJson.ToJson(lastKdata), &newKdata)
		}
		newData.NowKdata = newKdata
		newData.TradeKdataOpt = item.TradeKdataOpt

		// 更新当前持仓
		var Position PositionType
		jsoniter.Unmarshal(mJson.ToJson(item.NowVirtualPosition), &Position)
		newData.NowVirtualPosition = Position

		Upl := mCount.Div(item.NowVirtualPosition.NowUplRatio, "100")         // 格式化收益率
		Money := item.NowVirtualPosition.Money                                // 提取 Money
		MakeMoney := mCount.Mul(Money, Upl)                                   // 当前盈利的金钱
		Money = mCount.Add(Money, MakeMoney)                                  // 相加得出当账户应当剩余资金
		newData.NowVirtualPosition.Money = mCount.CentRound(Money, 3)         // 四舍五入
		newData.NowVirtualPosition.MakeMoney = mCount.CentRound(MakeMoney, 3) // 四舍五入

		HunterData[key] = newData
	}

	return HunterData
}
