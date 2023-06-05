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

// OKX 持仓
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

// OKX 账户余额
type AccountBalance struct {
	TimeUnix int64  `bson:"TimeUnix"`
	TimeStr  string `bson:"TimeStr"`
	CcyName  string `bson:"CcyName"` // 币种
	Balance  string `bson:"Balance"` // 币种
}

// 交易K线需要的参数
type TradeKdataOpt struct {
	EMA_Period    int    // EMA 步长 171
	CAP_Period    int    // CAP 步长 4
	CAP_Max       string // CAP 最大边界值
	CAP_Min       string // CAP 最小边界值
	MaxTradeLever int    // 最大杠杆数
	FullRun       bool
}

// 模拟持仓的数据
type VirtualPositionType struct {
	InstID       string        `bson:"InstID"`       // 下单币种 | 运行中设置
	HunterName   string        `bson:"HunterName"`   // 策略名称 | 运行中设置
	NowTimeStr   string        `bson:"NowTimeStr"`   // 当前K线时间 | 运行中设置
	NowTime      int64         `bson:"NowTime"`      // 当前实际时间戳 | 运行中设置
	NowC         string        `bson:"NowC"`         // 当前收盘价 | 运行中设置
	CAP_EMA      string        `bson:"CAP_EMA"`      // 当前的 CAP 值 | 运行中设置
	EMA          string        `bson:"EMA"`          // 当前的 EMA 值 | 运行中设置
	HunterConfig TradeKdataOpt `bson:"HunterConfig"` // 当前的交易K线参数  | 运行中设置
	// 下单时设置
	OpenAvgPx   string `bson:"OpenAvgPx"`   // 开仓价格 | 下单时设置
	OpenTimeStr string `bson:"OpenTimeStr"` // 开仓K线时间 | 下单时设置
	OpenTime    int64  `bson:"OpenTime"`    // 开仓实际时间戳  | 下单时设置
	NowDir      int    `bson:"NowDir"`      // 当前持仓状态 没持仓0  持多仓 1  持空仓 -1 | 初始化设置为0，下单时更新

	// 通过原始数据计算得出
	InitMoney   string `bson:"InitMoney"`   // 初始金钱 | 固定值初始值设置
	ChargeUpl   string `bson:"ChargeUpl"`   // 当前手续费率 | 固定值初始值设置
	NowUplRatio string `bson:"NowUplRatio"` // 当前未实现收益率(计算得出) | 运行中设置
	Money       string `bson:"Money"`       // 账户当前余额 | 如果没有初始值设置一次，下单时计算
	MaxMoney    string `bson:"MaxMoney"`    // 账户历史最高余额
	FloatMoney  string `bson:"FloatMoney"`  // 账户浮动余额
}

type UserOrderTable struct {
	OkxPositions    PositionsData       `bson:"OkxPositions"`
	OKXBalance      []AccountBalance    `bson:"OKXBalance"`
	OkxKey          OkxKeyType          `bson:"OkxKey"`
	VirtualPosition VirtualPositionType `bson:"VirtualPosition"` // 当前的虚拟持仓 数据库 OrderArr 最后一位
	UserID          string              `bson:"UserID"`
	OrderID         string              `bson:"OrderID"`
	CreateTime      int64               `bson:"CreateTime"` // 创建时间
}
