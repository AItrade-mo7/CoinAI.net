package config

var SysEnv = struct {
	MongoAddress  string
	MongoPassword string
	MongoUserName string
}{
	MongoAddress:  "trade.mo7.cc:17017",
	MongoPassword: "asdasd55555",
	MongoUserName: "mo7",
}

var AppEnv struct {
	Port       string `json:"Port"`
	IP         string `json:"IP"`
	UserID     string `json:"UserID"`
	ApiKey     string `json:"ApiKey"`
	SecretKey  string `json:"SecretKey"`
	Passphrase string `json:"Passphrase"`
	RunMod     int    // 0 则为正常模式 ， 1 则为数据模拟模式
}

var AppInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
