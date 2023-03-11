package config

// 系统通知的邮箱
var SysEmail = "meichangliang@outlook.com"

var LeverOpt = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

type OkxKeyType struct{}

type AppEnvType struct {
	SysName      string       `bson:"SysName"`      // 系统的名字  ， 自动生成项
	SysVersion   string       `bson:"SysVersion"`   // 系统的版本  ， 自动回填
	UserID       string       `bson:"UserID"`       // 用户名字 必填项  ， 禁止野生主机的存在
	Port         string       `bson:"Port"`         // 系统运行的端口 , 用户必填项
	IP           string       `bson:"IP"`           // 系统运行的 IP, 为自动获取回填
	ServeID      string       `bson:"ServeID"`      // ServeID ，  ip+端口
	MaxApiKeyNum int          `bson:"MaxApiKeyNum"` // 最大 Api 数量限制
	ApiKeyList   []OkxKeyType `bson:"ApiKeyList"`   // ApiKey 列表
}

var AppEnv AppEnvType
