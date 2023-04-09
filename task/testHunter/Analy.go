package testHunter

import (
	"fmt"

	"github.com/EasyGolang/goTools/mCount"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/EasyGolang/goTools/mTime"
)

func (_this *MockObj) Analy() {
	NowKdata := _this.TradeKdataList[len(_this.TradeKdataList)-1]
	AnalyDir := 0 // 分析的方向，默认为 0 不开仓

	if mCount.Le(NowKdata.CAP_EMA, _this.TradeKdataOpt.CAP_Max) > 0 { // 大于 CAPMax 则开多
		AnalyDir = 1
	}

	if mCount.Le(NowKdata.CAP_EMA, "-"+_this.TradeKdataOpt.CAP_Max) < 0 { // 小于 负 的 CAPMax 则开空
		AnalyDir = -1
	}

	// 更新持仓状态
	_this.CountPosition()
	// if mCount.Le(_this.NowPosition.UplRatio, _this.Billing.PositionMinRatio.Value) < 0 {
	// 	_this.Billing.PositionMinRatio.Value = _this.NowPosition.UplRatio
	// 	_this.Billing.PositionMinRatio.TimeStr = _this.NowPosition.NowTimeStr
	// }
	// if mCount.Le(_this.NowPosition.UplRatio, _this.Billing.PositionMaxRatio.Value) > 0 {
	// 	_this.Billing.PositionMaxRatio.Value = _this.NowPosition.UplRatio
	// 	_this.Billing.PositionMaxRatio.TimeStr = _this.NowPosition.NowTimeStr
	// }

	// 当前持仓与 判断方向不符合时，执行一次下单操作
	if _this.NowVirtualPosition.NowDir != AnalyDir {
		_this.OnOrder(AnalyDir)
	}

	// global.KdataLog.Println(_this.Billing.MockName, NowKdata.TimeStr, AnalyDir)
}

func (_this *MockObj) CountPosition() {
	NowKTradeData := _this.TradeKdataList[len(_this.TradeKdataList)-1]
	fmt.Println(NowKTradeData.TimeStr)

	_this.NowVirtualPosition.HunterName = _this.HunterName
	_this.NowVirtualPosition.NowTimeStr = NowKTradeData.TimeStr
	_this.NowVirtualPosition.NowTime = mTime.GetTime().TimeUnix
	_this.NowVirtualPosition.NowC = NowKTradeData.C
	_this.NowVirtualPosition.CAP_EMA = NowKTradeData.CAP_EMA
	_this.NowVirtualPosition.InstID = NowKTradeData.InstID
	_this.NowVirtualPosition.HunterConfig = NowKTradeData.Opt

	if _this.NowVirtualPosition.NowDir != 0 { // 当前为持仓状态，则计算收益率
		UplRatio := mCount.RoseCent(NowKTradeData.C, _this.NowVirtualPosition.OpenAvgPx)
		if _this.NowVirtualPosition.NowDir < 0 { // 当前为持空仓状态则翻转该值
			UplRatio = mCount.Sub("0", UplRatio)
		}
		_this.NowVirtualPosition.NowUplRatio = mCount.Mul(UplRatio, mStr.ToStr(_this.TradeKdataOpt.MaxTradeLever)) // 乘以杠杆倍数
		// 这里应该在客户端执行 金钱收益  Money * UplRatio = NowMoney
	}
}

// // 下单  参数：dir 下单方向 NowKdata : 当前市场行情
func (_this *MockObj) OnOrder(dir int) {
	// 	_this.BillingFun(NowKdata) // 下单之前 计算一次收益

	// 	if dir > 0 { // 开多
	// 		// 下订单
	// 		_this.OrderArr = append(_this.OrderArr, OrderType{
	// 			Type:    "Buy",            // 下多单
	// 			InstID:  NowKdata.InstID,  // 下单币种
	// 			AvgPx:   NowKdata.C,       // 记录下单价格
	// 			TimeStr: NowKdata.TimeStr, // 记录下单时间
	// 		})
	// 		// 更新持仓状态
	// 		_this.NowPosition = PositionType{
	// 			Dir:         1,          // 持仓多方向
	// 			OpenAvgPx:   NowKdata.C, // 持仓价格
	// 			NowTimeStr:  NowKdata.TimeStr,
	// 			UplRatio:    "0", // 当前收益率
	// 			NowC:        NowKdata.C,
	// 			OpenTimeStr: NowKdata.TimeStr, // 开仓时间
	// 			InstID:      NowKdata.InstID,  // 开仓币种
	// 			CAP_EMA:     NowKdata.CAP_EMA,
	// 		}
	// 	}

	// 	if dir < 0 { // 开空
	// 		// 下订单
	// 		_this.OrderArr = append(_this.OrderArr, OrderType{
	// 			Type:    "Sell",           // 下空单
	// 			InstID:  NowKdata.InstID,  // 下单币种
	// 			AvgPx:   NowKdata.C,       // 记录下单价格
	// 			TimeStr: NowKdata.TimeStr, // 记录下单时间
	// 		})
	// 		// 更新持仓状态
	// 		_this.NowPosition = PositionType{
	// 			Dir:         -1,         // 持仓空方向
	// 			OpenAvgPx:   NowKdata.C, // 持仓价格
	// 			NowTimeStr:  NowKdata.TimeStr,
	// 			UplRatio:    "0", // 当前收益率
	// 			NowC:        NowKdata.C,
	// 			OpenTimeStr: NowKdata.TimeStr, // 开仓时间
	// 			InstID:      NowKdata.InstID,  // 开仓币种
	// 			CAP_EMA:     NowKdata.CAP_EMA,
	// 		}
	// 	}

	// if dir == 0 { // 平仓
	//
	//		// 下订单
	//		_this.OrderArr = append(_this.OrderArr, OrderType{
	//			Type:    "Close",          // 平仓
	//			InstID:  NowKdata.InstID,  // 下单币种
	//			AvgPx:   NowKdata.C,       // 记录下单价格
	//			TimeStr: NowKdata.TimeStr, // 记录下单时间
	//		})
	//		// 更新为空仓状态
	//		_this.NowPosition = PositionType{
	//			Dir:         0,  // 持仓空方向
	//			OpenAvgPx:   "", // 持仓价格
	//			NowTimeStr:  NowKdata.TimeStr,
	//			UplRatio:    "0", // 当前收益率
	//			NowC:        NowKdata.C,
	//			OpenTimeStr: NowKdata.TimeStr, // 开仓时间
	//			InstID:      NowKdata.InstID,  // 开仓币种
	//			CAP_EMA:     NowKdata.CAP_EMA,
	//		}
	//	}
}
