package dbType

/*
用来存储验证 Token
collection : VerifyToken
*/
type TokenTable struct {
	UserID     string `bson:"UserID"`     // 用户 ID
	Token      string `bson:"Token"`      // 当前登录的Token
	CreateTime int64  `bson:"CreateTime"` // 创建时间
}
