package hunter

import (
	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/global/dbType"
	"github.com/EasyGolang/goTools/mCount"
	"github.com/EasyGolang/goTools/mEncrypt"
	"github.com/EasyGolang/goTools/mFile"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mMongo"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mTime"
	jsoniter "github.com/json-iterator/go"
)

// // 下单  参数：dir 下单方向 NowKdata : 当前市场行情
func (_this *HunterObj) OnOrder(dir int) {
	NowKTradeData := _this.TradeKdataList[len(_this.TradeKdataList)-1]

	// 结算本期持仓
	_this.BillingFun()

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
	_this.OrderOpen() // 这里的结果 要么是 1 要么是 0  要么 是 -1 没有第三种了
	global.TradeLog.Println(_this.HunterName, "下单一次", mJson.ToStr(_this.NowVirtualPosition))
	_this.OrderArr = append(_this.OrderArr, _this.NowVirtualPosition)
	mFile.Write(_this.OutPutDirectory+"/OrderArr.json", mJson.ToStr(_this.OrderArr))
}

func (_this *HunterObj) BillingFun() {
	// 这里是下单之前的结算周期
	Upl := mCount.Div(_this.NowVirtualPosition.NowUplRatio, "100")     // 格式化收益率
	ChargeUpl := mCount.Div(_this.NowVirtualPosition.ChargeUpl, "100") // 格式化手续费率

	Money := _this.NowVirtualPosition.Money // 提取 Money
	makeMoney := mCount.Mul(Money, Upl)     // 当前盈利的金钱
	Money = mCount.Add(Money, makeMoney)    // 相加得出当账户剩余资金

	nowCharge := mCount.Mul(Money, ChargeUpl) // 当前产生的手续费
	Money = mCount.Sub(Money, nowCharge)      // 减去手续费
	Money = mCount.CentRound(Money, 3)        // 四舍五入保留三位小数
	_this.NowVirtualPosition.Money = Money    // 保存结果到当前持仓
	global.Run.Println("结算一次", mJson.ToStr(_this.NowVirtualPosition))
}

func (_this *HunterObj) OrderOpen() {
	// 在这里进行 下单存储。
	global.Run.Println("下单", mJson.ToStr(_this.NowVirtualPosition))
	if _this.NowVirtualPosition.NowDir > 0 {
		_this.SetOrderDB("Buy")
	}
	if _this.NowVirtualPosition.NowDir < 0 {
		_this.SetOrderDB("Sell")
	}

	if _this.NowVirtualPosition.NowDir == 0 {
		_this.SetOrderDB("Close")
	}
}

func (_this *HunterObj) SetOrderDB(Type string) {
	var orderData dbType.CoinOrderTable
	jsoniter.Unmarshal(mJson.ToJson(_this.NowVirtualPosition), &orderData)
	orderData.CreateTime = mTime.GetUnixInt64()
	orderData.Type = Type
	orderData.ServeID = config.AppEnv.ServeID
	orderData.TimeID = mOKX.GetTimeID(orderData.NowTime)
	orderData.OrderID = mEncrypt.GetUUID()

	db := mMongo.New(mMongo.Opt{
		UserName: config.SysEnv.MongoUserName,
		Password: config.SysEnv.MongoPassword,
		Address:  config.SysEnv.MongoAddress,
		DBName:   "AIServe",
	}).Connect().Collection("CoinOrder")
	defer db.Close()
	_, err := db.Table.InsertOne(db.Ctx, orderData)
	if err != nil {
		global.LogErr("hunter.SetOrderDB 数据存储失败", _this.HunterName, err)
	}
}
