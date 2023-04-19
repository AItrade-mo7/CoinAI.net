package config

var SysName = "CoinAI.net"

var GithubPackagePath = struct {
	Origin string
	Path   string
}{
	Origin: "https://raw.githubusercontent.com",
	Path:   "/AItrade-mo7/CoinAIPackage/main/package.json",
}

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
	SysEnv.MongoAddress = "tcy4.mo7.cc:17017"
	SysEnv.MongoPassword = "asdasd55555"
	SysEnv.MongoUserName = "mo7"
	SysEnv.MessageBaseUrl = "http://msg.mo7.cc"
}
