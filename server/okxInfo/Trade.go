package okxInfo

var Lever = "10" // 默认的杠杆倍数

// Hunter 的值
var HunterRun struct {
	InstID string `bson:"InstID"` // 正在交易的 InstID
}

func SetHunterInstID(InstID string) {
	HunterRun.InstID = "AVAX-USDT-SWAP" // 写死为默认值
}
