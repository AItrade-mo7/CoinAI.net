package wss

import (
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mTime"
)

type OutPut struct {
	config.AppEnvType
	SysTime    int64         `bson:"SysTime"` // 系统时间
	DataSource string        `bson:"DataSource"`
	Leverage   string        `bson:"Leverage"`
	TradeModel string        `bson:"TradeModel"`
	OrderInst  mOKX.TypeInst `bson:"OrderInst"`
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
	resData.Leverage = okxInfo.Leverage
	resData.TradeModel = okxInfo.TradeModel
	resData.OrderInst = okxInfo.OrderInst

	return
}

func GetFuzzyApiKey() []mOKX.TypeOkxKey {
	ApiKeyList := config.AppEnv.ApiKeyList

	NewKeyList := []mOKX.TypeOkxKey{}

	for _, val := range ApiKeyList {
		NewKeyList = append(NewKeyList, mOKX.TypeOkxKey{
			Name:       val.Name,
			ApiKey:     GetKeyFuzzy(val.ApiKey),
			SecretKey:  GetKeyFuzzy(val.SecretKey),
			Passphrase: GetKeyFuzzy(val.Passphrase),
			IsTrade:    val.IsTrade,
		})
	}

	return NewKeyList
}

func GetKeyFuzzy(Ket string) string {
	return Ket[0:6] + "******" + Ket[len(Ket)-6:len(Ket)-1]
}
