package config

var SysEnv = struct {
	MongoAddress  string
	MongoPassword string
	MongoUserName string
}{
	MongoAddress:  "aws.mo7.cc:17017",
	MongoPassword: "asdasd55555",
	MongoUserName: "mo7",
}

var AppEnv struct {
	Port       string `bson:"Port"`
	IP         string `bson:"IP"`
	UserID     string `bson:"UserID"`
	ApiKey     string `bson:"ApiKey"`
	SecretKey  string `bson:"SecretKey"`
	Passphrase string `bson:"Passphrase"`
	RunMod     int    // 0 则为正常模式 ， 1 则为数据模拟模式
}

var AppInfo struct {
	Name    string `bson:"name"`
	Version string `bson:"version"`
}
