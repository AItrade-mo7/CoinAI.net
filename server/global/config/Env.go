package config

import (
	"github.com/EasyGolang/goTools/mOKX"
)

var SysName = "CoinAI.net"

var AppInfo struct {
	Name    string `bson:"name"`
	Version string `bson:"version"`
}

var SysEnv = struct {
	MongoAddress  string
	MongoPassword string
	MongoUserName string
}{
	MongoAddress:  "aws.mo7.cc:17017",
	MongoPassword: "asdasd55555",
	MongoUserName: "mo7",
}

var (
	SysEmail = "trade@mo7.cc"
	MyIP     = "3.1.68.15"
)

type AppEnvType struct {
	Name         string            `bson:"Name"`
	Version      string            `bson:"Version"`
	Port         string            `bson:"Port"`
	IP           string            `bson:"IP"`
	ServeID      string            `bson:"ServeID"`
	UserID       string            `bson:"UserID"`
	TradeLever   int               `bson:"TradeLever"`
	IsSPOT       bool              `bson:"IsSPOT"`
	MaxApiKeyNum int               `bson:"MaxApiKeyNum"`
	ApiKeyList   []mOKX.TypeOkxKey `bson:"ApiKeyList"`
}

var LeverOpt = []int{2, 3, 4, 5, 6, 7, 8, 9, 10}

var AppEnv AppEnvType

func GetOKXKey(Index int) mOKX.TypeOkxKey {
	ReturnKey := mOKX.TypeOkxKey{}

	for key, val := range AppEnv.ApiKeyList {
		if key == Index {
			ReturnKey = val
			break
		}
	}

	return ReturnKey
}

type EmailInfo struct {
	Account  string
	Password string
	To       []string
}

var Email = EmailInfo{
	Account:  SysEmail,
	Password: "Mcl931750",
	To:       []string{},
}

var PublicUserID = ""
