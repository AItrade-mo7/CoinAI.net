package okxInfo

import "github.com/EasyGolang/goTools/mOKX"

// 杠杆倍数
var Leverage = "5"

func SetLeverage(lever string) {
	Leverage = lever
}

// 交易模式
var TradeModelType = map[int]string{
	1: "合约模式",
	2: "现货模式",
}
var TradeModel = "1"

func SetTradeModel(model string) {
	TradeModel = model
}

// 需要监听的下单产品的信息
var OrderInst mOKX.TypeInst

func SetOrderInst(CcyName string) {
	// if TradeModel == 1 {
	// 	OrderInst = CcyName + "-SWAP"
	// }
	// if TradeModel == 2 {
	// 	OrderInst = CcyName + "-USDT"
	// }
}
