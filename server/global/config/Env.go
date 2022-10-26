package config

import (
	"github.com/EasyGolang/goTools/mOKX"
)

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

type AppEnvType struct {
	Name       string            `bson:"Name"`
	Version    string            `bson:"Version"`
	Port       string            `bson:"Port"`
	IP         string            `bson:"IP"`
	ServeID    string            `bson:"ServeID"`
	UserID     string            `bson:"UserID"`
	ApiKeyList []mOKX.TypeOkxKey `bson:"ApiKeyList"`
}

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
	Account:  "trade@mo7.cc",
	Password: "Mcl931750",
	To: []string{
		"meichangliang@mo7.cc",
	},
}

var PublicUserID = "468f3c7c0c684de181cad3b1fe34fab1"
