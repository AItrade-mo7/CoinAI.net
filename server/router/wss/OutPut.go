package wss

import (
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mTime"
)

type OutPut struct {
	config.AppEnvType
	SysTime    int64                   `bson:"SysTime"` // 系统时间
	DataSource string                  `bson:"DataSource"`
	TradeInst  mOKX.TypeInst           `bson:"TradeInst"`
	TradeLever int                     `bson:"TradeLever"`
	NowTicker  okxInfo.AnalyTickerType `bson:"NowTicker"`
	TradeCoin  mOKX.TypeTicker         `bson:"NowTicker"`
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
	resData.ApiKeyList = GetFuzzyApiKey()
	// 系统时间
	resData.SysTime = mTime.GetUnixInt64()
	resData.DataSource = "CoinAI.net"
	// 下单信息
	if okxInfo.IsSPOT {
		resData.TradeInst = okxInfo.TradeInst.SPOT
	} else {
		resData.TradeInst = okxInfo.TradeInst.SWAP
	}
	resData.TradeLever = okxInfo.TradeLever

	// Ticker 信息
	resData.NowTicker.Unit = okxInfo.NowTicker.Unit
	resData.NowTicker.WholeDir = okxInfo.NowTicker.WholeDir
	resData.NowTicker.DirIndex = okxInfo.NowTicker.DirIndex
	resData.NowTicker.TimeStr = okxInfo.NowTicker.TimeStr

	// Coin 信息
	resData.TradeCoin.Last = okxInfo.TradeInst.Last
	resData.TradeCoin.TimeUnix = okxInfo.TradeInst.TimeUnix
	resData.TradeCoin.TimeStr = okxInfo.TradeInst.TimeStr

	return
}

func GetFuzzyApiKey() []mOKX.TypeOkxKey {
	ApiKeyList := config.AppEnv.ApiKeyList

	NewKeyList := []mOKX.TypeOkxKey{}

	for _, val := range ApiKeyList {
		NewKeyList = append(NewKeyList, mOKX.TypeOkxKey{
			Name:       val.Name,
			ApiKey:     GetKeyFuzzy(val.ApiKey, 5),
			SecretKey:  GetKeyFuzzy(val.SecretKey, 5),
			Passphrase: GetKeyFuzzy(val.Passphrase, 1),
			IsTrade:    val.IsTrade,
			UserID:     val.UserID,
		})
	}

	return NewKeyList
}

func GetKeyFuzzy(Ket string, num int) string {
	return Ket[0:num] + "****" + Ket[len(Ket)-num-1:len(Ket)-1]
}
