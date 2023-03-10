package config

var SysName = "CoinAI.net"

var AppInfo struct {
	Name    string `bson:"name"`
	Version string `bson:"version"`
}

var SysEnv struct {
	MongoAddress   string
	MongoPassword  string
	MongoUserName  string
	MessageBaseUrl string
}

func DefaultSysEnv() {
	SysEnv.MongoAddress = "tcy.mo7.cc:17017"
	SysEnv.MongoPassword = "mo7"
	SysEnv.MongoUserName = "asdasd55555"
	SysEnv.MessageBaseUrl = "http://msg.mo7.cc"
}
