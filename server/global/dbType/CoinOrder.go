package dbType

import "CoinAI.net/server/okxInfo"

type CoinOrderTable struct {
	InstID       string                `bson:"InstID"`       // 下单币种 | 运行中设置
	HunterName   string                `bson:"HunterName"`   // 策略名称 | 运行中设置
	NowTimeStr   string                `bson:"NowTimeStr"`   // 当前K线时间 | 运行中设置
	NowTime      int64                 `bson:"NowTime"`      // 当前实际时间戳 | 运行中设置
	NowC         string                `bson:"NowC"`         // 当前收盘价 | 运行中设置
	CAP_EMA      string                `bson:"CAP_EMA"`      // 当前的 CAP 值 | 运行中设置
	EMA          string                `bson:"EMA"`          // 当前的 EMA 值 | 运行中设置
	HunterConfig okxInfo.TradeKdataOpt `bson:"HunterConfig"` // 当前的交易K线参数  | 运行中设置
	OpenAvgPx    string                `bson:"OpenAvgPx"`    // 开仓价格 | 下单时设置
	OpenTimeStr  string                `bson:"OpenTimeStr"`  // 开仓K线时间 | 下单时设置
	OpenTime     int64                 `bson:"OpenTime"`     // 开仓实际时间戳  | 下单时设置
	NowDir       int                   `bson:"NowDir"`       // 当前持仓状态 没持仓0  持多仓 1  持空仓 -1 | 初始化设置为0，下单时更新
	CreateTime   int64                 `bson:"CreateTime"`   // 创建时间
	Type         string                `bson:"Type"`         // 当前订单类型  Close:平仓  Buy:买多  Sell:卖空
	ServeID      string                `bson:"ServeID"`      // ServeID ，  ip+端口
	TimeID       string                `bson:"TimeID"`       // 精确到分钟
	OrderID      string                `bson:"OrderID"`      // UUID
}
