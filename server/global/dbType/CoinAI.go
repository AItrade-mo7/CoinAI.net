package dbType

/*
用来存储 用户信息
db: AIServe
collection : CoinAI
*/

type OkxKeyType struct {
	Name       string `bson:"Name"`
	ApiKey     string `bson:"ApiKey"` // 当做ID用
	SecretKey  string `bson:"SecretKey"`
	Passphrase string `bson:"Passphrase"`
	UserID     string `bson:"UserID"` // 用户 ID 必填项 ，禁止野生账户的存在
	Status     string `bson:"Status"` //  disable  enable del
}

type AppEnvType struct {
	SysName      string       `bson:"SysName"`      // 系统的名字  ， 自动生成项
	SysVersion   string       `bson:"SysVersion"`   // 系统的版本  ， 自动回填
	UserID       string       `bson:"UserID"`       // 用户名字 必填项  ， 禁止野生主机的存在
	Port         string       `bson:"Port"`         // 系统运行的端口 , 用户必填项
	IP           string       `bson:"IP"`           // 系统运行的 IP, 为自动获取回填
	ServeID      string       `bson:"ServeID"`      // ServeID ，  ip+端口
	MaxApiKeyNum int          `bson:"MaxApiKeyNum"` // 最大 Api 数量限制
	Type         string       `bson:"Type"`
	CreateTime   int64        `bson:"CreateTime"` // 创建时间
	UpdateTime   int64        `bson:"UpdateTime"` // 更新时间
	ApiKeyList   []OkxKeyType `bson:"ApiKeyList"` // ApiKey 列表
}
