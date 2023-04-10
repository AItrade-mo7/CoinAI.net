package hunter

import (
	"fmt"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/dbType"
	"github.com/EasyGolang/goTools/mCount"
	"github.com/EasyGolang/goTools/mFile"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mTime"
)

// // 下单  参数：dir 下单方向 NowKdata : 当前市场行情
func (_this *HunterObj) OnOrder(dir int) {
	NowKTradeData := _this.TradeKdataList[len(_this.TradeKdataList)-1]

	// 在这里计算当前的 Money
	Upl := mCount.Div(_this.NowVirtualPosition.NowUplRatio, "100")     // 格式化收益率
	ChargeUpl := mCount.Div(_this.NowVirtualPosition.ChargeUpl, "100") // 格式化手续费率

	Money := _this.NowVirtualPosition.Money // 提取 Money
	makeMoney := mCount.Mul(Money, Upl)     // 当前盈利的金钱
	Money = mCount.Add(Money, makeMoney)    // 相加得出当账户剩余资金

	nowCharge := mCount.Mul(Money, ChargeUpl) // 当前产生的手续费
	Money = mCount.Sub(Money, nowCharge)      // 减去手续费
	Money = mCount.CentRound(Money, 3)        // 四舍五入保留三位小数
	_this.NowVirtualPosition.Money = Money    // 保存结果到当前持仓

	// 在这里执行平仓, 平掉所有仓位
	_this.OrderClose()

	// 同步持仓状态, 相当于下单了
	if dir > 0 {
		// 开多
		_this.NowVirtualPosition.NowDir = 1
	}
	if dir < 0 {
		// 开空
		_this.NowVirtualPosition.NowDir = -1
	}
	// 同步下单价格
	_this.NowVirtualPosition.OpenAvgPx = NowKTradeData.C
	_this.NowVirtualPosition.OpenTimeStr = NowKTradeData.TimeStr
	_this.NowVirtualPosition.OpenTime = mTime.GetTime().TimeUnix

	// 同步平仓状态
	if dir == 0 {
		_this.NowVirtualPosition.NowDir = 0
		_this.NowVirtualPosition.OpenAvgPx = ""
		_this.NowVirtualPosition.OpenTimeStr = ""
		_this.NowVirtualPosition.OpenTime = 0
	}
	// 平仓后未实现盈亏重置为 0
	_this.NowVirtualPosition.NowUplRatio = "0"

	// 在这里执行下单
	_this.OrderOpen()
	global.TradeLog.Println(_this.HunterName, "下单一次", mJson.ToStr(_this.NowVirtualPosition))
	_this.OrderArr = append(_this.OrderArr, _this.NowVirtualPosition)
	mFile.Write(_this.OutPutDirectory+"/OrderArr.json", mJson.ToStr(_this.OrderArr))
}

func (_this *HunterObj) OrderClose() {
	// 在这里优先平掉所有仓位  在这里进行平仓结算 和 持仓状态存储
	global.Run.Println("平仓", mJson.ToStr(_this.NowVirtualPosition))
	// 数据库存储一次 平仓 通知 Message 去存储
}

func (_this *HunterObj) OrderOpen() {
	// 在这里进行 下单存储。
	global.Run.Println("下单", mJson.ToStr(_this.NowVirtualPosition))

	orderData := dbType.CoinOrderTable{}
	orderData.OrderInfo = _this.NowVirtualPosition
	fmt.Println(orderData)

	/*
		_this.NowVirtualPosition.NowDir > 0 则开多单
		_this.NowVirtualPosition.NowDir < 0 则开空单
	*/

	// 数据库存储一次 开仓 通知 Message 去存储
}
