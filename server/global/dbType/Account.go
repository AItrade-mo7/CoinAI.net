package dbType

/*
用来存储 用户信息
db: Account
collection : User
*/

type UserTable struct {
	UserID         string   `bson:"UserID"`         // 用户 ID
	Email          string   `bson:"Email"`          // 用户主要的 Email
	UserEmail      []string `bson:"UserEmail"`      // 用户的 Email 列表
	Avatar         string   `bson:"Avatar"`         // 用户头像
	NickName       string   `bson:"NickName"`       // 用户昵称
	CreateTime     int64    `bson:"CreateTime"`     // 创建时间
	UpdateTime     int64    `bson:"UpdateTime"`     // 更新时间
	EntrapmentCode string   `bson:"EntrapmentCode"` // 防钓鱼码
	Password       string   `bson:"Password"`       // 用户密码
}

/*
用来存储 用户的平仓订单
db: Account
collection : Order
*/

type PositionsData struct {
	AvgPx       string `bson:"AvgPx"`       // 开仓均价
	CTime       string `bson:"CTime"`       // 持仓创建时间
	Ccy         string `bson:"Ccy"`         // 币种
	InstID      string `bson:"InstID"`      // InstID
	InstType    string `bson:"InstType"`    // SWAP
	Interest    string `bson:"Interest"`    // 利息
	Last        string `bson:"Last"`        // 当前最新成交价
	Lever       string `bson:"Lever"`       // 杠杆倍数
	LiqPx       string `bson:"LiqPx"`       // 预估强平价格
	MarkPx      string `bson:"MarkPx"`      // 标记价格
	MgnRatio    string `bson:"MgnRatio"`    // 保证金率
	NotionalUsd string `bson:"NotionalUsd"` // 持仓数量
	Pos         string `bson:"Pos"`         // 持仓数量
	UTime       string `bson:"UTime"`       // 更新时间
	Upl         string `bson:"Upl"`         // 未实现收益
	UplRatio    string `bson:"UplRatio"`    // 未实现收益率
	Imr         string `bson:"Imr"`         // 初始保证金
}

type UserOrderTable struct {
	OkxPositions PositionsData `bson:"OkxPositions"`
	OkxKey       OkxKeyType    `bson:"OkxKey"`
	UserID       string        `bson:"UserID"`
	OrderID      string        `bson:"OrderID"`
	CreateTime   int64         `bson:"CreateTime"` // 创建时间
}
