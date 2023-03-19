package wss

import (
	"CoinAI.net/server/global/config"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mTime"
	jsoniter "github.com/json-iterator/go"
)

type OutAnalyTickerType struct {
	Version  int    `bson:"Version"`
	Unit     string `bson:"Unit"`
	TimeUnix int64  `bson:"TimeUnix"`
	TimeStr  string `bson:"TimeStr"`
}

type TradeCoinType struct{}

type MainUserType struct {
	UserID    string   `bson:"UserID"`    // 用户 ID
	Email     string   `bson:"Email"`     // 用户主要的 Email
	UserEmail []string `bson:"UserEmail"` // 用户的 Email 列表
	Avatar    string   `bson:"Avatar"`    // 用户头像
	NickName  string   `bson:"NickName"`  // 用户昵称
}

type OutPut struct {
	SysName      string       `bson:"SysName"`
	SysVersion   string       `bson:"SysVersion"`
	Port         string       `bson:"Port"`
	IP           string       `bson:"IP"`
	ServeID      string       `bson:"ServeID"`
	UserID       string       `bson:"UserID"`
	SysTime      int64        `bson:"SysTime"` // 系统时间
	DataSource   string       `bson:"DataSource"`
	MaxApiKeyNum int          `bson:"MaxApiKeyNum"`
	ApiKeyNum    int          `bson:"ApiKeyNum"`
	Type         string       `bson:"Type"`
	MainUser     MainUserType `bson:"MainUser"`
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
	// 当前管理员信息
	resData.MainUser = GetMainUser()

	return
}

func GetMainUser() (resData MainUserType) {
	resData = MainUserType{}
	jsoniter.Unmarshal(mJson.ToJson(config.MainUser), &resData)
	return
}
