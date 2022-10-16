package order

import (
	"fmt"

	"CoinAI.net/server/global"
	"github.com/EasyGolang/goTools/mCount"
	"github.com/EasyGolang/goTools/mEncrypt"
	"github.com/EasyGolang/goTools/mOKX"
)

type OrderParam struct {
	Side string // buy：买， sell：卖
	Sz   string // 委托数量
	Px   string // 委托价格
}

func Order(opt OrderParam) {
	// Side 参数判断
	if opt.Side == "buy" || opt.Side == "sell" {
	} else {
		global.LogErr("order.Order opt.Side 值不正确", opt.Side)
		return
	}
	// 委托数量判断
	Sz := mCount.Add(opt.Sz, "0")
	if mCount.Le(Sz, "0") < 1 {
		global.LogErr("order.Order opt.Sz 值不正确", opt.Sz)
		return
	}
	// 委托价格判断
	Px := mCount.Add(opt.Px, "0")
	if mCount.Le(Px, "0") < 1 {
		global.LogErr("order.Order opt.Px 值不正确", opt.Px)
		return
	}

	// 生成订单ID
	clOrdId := mEncrypt.TimeID()

	// 在这里要强制设置为单向持仓模式

	resData, err := mOKX.FetchOKX(mOKX.OptFetchOKX{
		Path: "/api/v5/trade/order",
		Data: map[string]any{
			"tdMode":  "isolated", // 逐仓杠杆
			"clOrdId": clOrdId,    // 根据时间生成自定义ID
			"side":    opt.Side,   // 方向
			"ordType": "limit",    // 固定限价单
			"sz":      Sz,         // 委托数量（张）
			"px":      Px,         // 委托价格
		},
		Method: "POST",
	})

	fmt.Println(resData, err)
}
