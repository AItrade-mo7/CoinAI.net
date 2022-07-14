package dbType

// 秘钥的 数据模板 ============ 对外展示  ===============
type OkxKeyTable struct {
	OkxKeyID   string `bson:"OkxKeyID"`   // OkxKeyID
	ApiKey     string `bson:"ApiKey"`     // ApiKey
	SecretKey  string `bson:"SecretKey"`  // SecretKey  密钥
	Passphrase string `bson:"Passphrase"` // Passphrase  密码
	IP         string `bson:"IP"`         // IP地址
	Name       string `bson:"Name"`       // 备注名,同一个账号下必须唯一
	Note       string `bson:"Note"`       // 备注
	CreateTime int64  `bson:"CreateTime"` // 创建时间
}

// 账户的 表结构  ========== Account ==============
type AccountTable struct {
	/* type UserInfo struct */
	Email        string `bson:"Email"`        // 用户 Email
	UserID       string `bson:"UserID"`       // 用户 ID
	Avatar       string `bson:"Avatar"`       // 用户头像
	NickName     string `bson:"NickName"`     // 用户昵称
	CreateTime   int64  `bson:"CreateTime"`   // 创建时间
	UpdateTime   int64  `bson:"UpdateTime"`   // 更新时间
	SecurityCode string `bson:"SecurityCode"` // 防伪标识
	/* type UserInfo struct */

	Password   string        `bson:"Password"` // 用户密码
	OkxKeyList []OkxKeyTable `bson:"OkxKeyList"`
}
