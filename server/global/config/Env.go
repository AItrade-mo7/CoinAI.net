package config

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

type ApiKeyList struct {
	Name       string `bson:"Name"`
	ApiKey     string `bson:"ApiKey"`
	SecretKey  string `bson:"SecretKey"`
	Passphrase string `bson:"Passphrase"`
}

var AppEnv struct {
	Name       string       `bson:"Name"`
	Version    string       `bson:"Version"`
	Port       string       `bson:"Port"`
	IP         string       `bson:"IP"`
	ServeID    string       `bson:"ServeID"`
	UserID     string       `bson:"UserID"`
	ApiKeyList []ApiKeyList `bson:"ApiKeyList"`
}
