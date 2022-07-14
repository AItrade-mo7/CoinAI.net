package apiType

// UserInfo 的数据结构  ======== 对外展示 =======

type UserInfo struct {
	Email        string `bson:"Email"`        // 用户 Email
	UserID       string `bson:"UserID"`       // 用户 ID
	Avatar       string `bson:"Avatar"`       // 用户头像
	NickName     string `bson:"NickName"`     // 用户昵称
	CreateTime   int64  `bson:"CreateTime"`   // 创建时间
	UpdateTime   int64  `bson:"UpdateTime"`   // 更新时间
	SecurityCode string `bson:"SecurityCode"` // 防伪标识
}

type LoginSucceedType struct {
	Email        string `bson:"Email"`        // 用户 Email
	UserID       string `bson:"UserID"`       // 用户 ID
	Avatar       string `bson:"Avatar"`       // 用户头像
	NickName     string `bson:"NickName"`     // 用户昵称
	CreateTime   int64  `bson:"CreateTime"`   // 创建时间
	UpdateTime   int64  `bson:"UpdateTime"`   // 更新时间
	SecurityCode string `bson:"SecurityCode"` // 防伪标识
	Token        string `bson:"Token"`        // Token
}
