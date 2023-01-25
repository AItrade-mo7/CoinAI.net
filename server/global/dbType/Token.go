package dbType

// Token 表结构  ========== Account ==============
type TokenTable struct {
	UserID     string `bson:"UserID"`     // 用户 ID
	Token      string `bson:"Token"`      // 当前登录的Token
	CreateTime int64  `bson:"CreateTime"` // 创建时间
}
