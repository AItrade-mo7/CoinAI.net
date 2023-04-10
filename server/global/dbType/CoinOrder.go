package dbType

import "CoinAI.net/server/okxInfo"

type CoinOrderTable struct {
	okxInfo.VirtualPositionType

	CreateTime int64  `bson:"CreateTime"` // 创建时间
	Type       string `bson:"Type"`       // 当前订单类型  Close:平仓  Buy:买多  Sell:卖空
	ServeID    string `bson:"ServeID"`    // ServeID ，  ip+端口
	TimeID     string `bson:"TimeID"`     // 精确到分钟
	OrderID    string `bson:"OrderID"`    // UUID
}
