package dbType

// 发送验证码的表结构 ========= EmailCode ============
type EmailCodeTable struct {
	Email    string `bson:"Email"`
	Code     string `bson:"Code"`
	SendTime int64  `bson:"SendTime"`
}
