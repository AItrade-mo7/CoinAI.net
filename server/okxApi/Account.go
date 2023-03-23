package okxApi

import (
	"fmt"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/global/dbType"
	"CoinAI.net/server/okxApi/restApi/account"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mCount"
	"github.com/EasyGolang/goTools/mOKX"
)

type AccountParam struct {
	OkxKey dbType.OkxKeyType
}

type AccountObj struct {
	OkxKey       dbType.OkxKeyType
	TradeInst    mOKX.TypeInst // 交易币种信息
	Balance      []account.AccountBalance
	Positions    []account.PositionsData
	MaxSize      account.MaxSizeType
	PendingOrder []account.PendingOrderType
}

// 创建一个新账户
func NewAccount(opt AccountParam) (resObj *AccountObj, resErr error) {
	resObj = &AccountObj{}
	resErr = nil

	if opt.OkxKey.Status == "disable" {
		resErr = fmt.Errorf("okxApi.NewAccount 当前 Key 已被禁用 Name:" + opt.OkxKey.Name)
		return
	}
	if len(opt.OkxKey.ApiKey) < 10 {
		resErr = fmt.Errorf("okxApi.NewAccount ApiKey 不能为空 Name:" + opt.OkxKey.Name)
		return
	}

	if opt.OkxKey.TradeLever < config.LeverOpt[0] {
		opt.OkxKey.TradeLever = config.LeverOpt[0]
	}

	if opt.OkxKey.TradeLever > config.LeverOpt[len(config.LeverOpt)-1] {
		opt.OkxKey.TradeLever = config.LeverOpt[len(config.LeverOpt)-1]
	}

	resObj.OkxKey = opt.OkxKey

	// 同步当前账户的交易 币种
	resObj.TradeInst = okxInfo.TradeInst

	resObj.GetPositions()    // 获取当前持仓
	resObj.SetPositionMode() // 设置持仓模式
	return
}

// 设置持仓模式
func (_this *AccountObj) SetPositionMode() (resErr error) {
	if len(_this.Positions) < 1 {
		resErr = account.SetPositionMode(_this.OkxKey)
	}
	return
}

// 下单 买多
func (_this *AccountObj) Buy() (resErr error) {
	_this.GetMaxSize() // 获取最大开仓数量
	Sz := _this.MaxSize.MaxBuy
	resErr = account.Order(account.OrderParam{
		OKXKey:    _this.OkxKey,
		TradeInst: _this.TradeInst,
		Side:      "buy",
		Sz:        Sz,
	})
	// 如果下单数量大于最大值，则再来一次
	if mCount.Le(Sz, _this.TradeInst.MaxMktSz) > 0 {
		_this.Buy()
	}
	return
}

// 下单 买空
func (_this *AccountObj) Sell() (resErr error) {
	_this.GetMaxSize() // 获取最大开仓数量
	Sz := _this.MaxSize.MaxSell
	account.Order(account.OrderParam{
		OKXKey:    _this.OkxKey,
		TradeInst: _this.TradeInst,
		Side:      "sell",
		Sz:        Sz,
	})
	// 如果下单数量大于最大值，则再来一次
	if mCount.Le(Sz, _this.TradeInst.MaxMktSz) > 0 {
		_this.Sell()
	}
	return
}

// 设置杠杆倍数
func (_this *AccountObj) SetLeverage() (resErr error) {
	resErr = account.SetLeverage(account.SetLeverageParam{
		InstID: _this.TradeInst.InstID,
		OKXKey: _this.OkxKey,
	})
	return
}

// 获取账户余额
func (_this *AccountObj) GetBalance() (resErr error) {
	resData, resErr := account.GetOKXBalance(_this.OkxKey)
	_this.Balance = resData
	return
}

// 获取持仓信息
func (_this *AccountObj) GetPositions() (resErr error) {
	resData, resErr := account.GetOKXPositions(_this.OkxKey)
	_this.Positions = resData
	return
}

// 获取最大可开仓数量
func (_this *AccountObj) GetMaxSize() (resErr error) {
	_this.SetLeverage() // 设置杠杆倍数
	resData, resErr := account.GetMaxSize(account.GetMaxSizeParam{
		InstID: _this.TradeInst.InstID,
		OKXKey: _this.OkxKey,
	})
	_this.MaxSize = resData
	return
}

// 未成交订单信息
func (_this *AccountObj) GetOrdersPending() (resErr error) {
	resData, resErr := account.GetOrdersPending(account.GetOrdersPendingParam{
		OKXKey: _this.OkxKey,
	})
	_this.PendingOrder = resData
	return
}

// 取消所有未成交订单
func (_this *AccountObj) CancelOrder() (resErr error) {
	errArr := []error{}
	for _, val := range _this.PendingOrder {
		err := account.CancelOrder(account.CancelOrderParam{
			OKXKey: _this.OkxKey,
			Order:  val,
		})
		if err != nil {
			errArr = append(errArr, err)
		}
	}
	if len(errArr) > 0 {
		resErr = fmt.Errorf("err:%+v", errArr)
	}
	return
}

// 下单 平仓,平掉当前所有仓位
func (_this *AccountObj) Close() (resErr error) {
	_this.GetOrdersPending() // 获取未成交订单
	_this.CancelOrder()      // 取消所有未成交订单
	_this.GetPositions()     // 获取所有持仓

	errArr := []error{}
	isAgin := false
	for _, Position := range _this.Positions {
		TradeInst := okxInfo.Inst[Position.InstID]
		Side := ""
		Sz := "0"

		maxSize, err := account.GetMaxSize(account.GetMaxSizeParam{
			InstID: TradeInst.InstID,
			OKXKey: _this.OkxKey,
		})
		if err != nil {
			err = fmt.Errorf("平仓 获取最大数量 失败")
			global.LogErr(err)
			errArr = append(errArr, err)
		}

		if mCount.Le(Position.Pos, "0") > 0 {
			Side = "sell"
			Sz = maxSize.MaxSell
		}
		if mCount.Le(Position.Pos, "0") < 0 {
			Side = "buy"
			Sz = maxSize.MaxBuy
		}

		err = account.Order(account.OrderParam{
			OKXKey:    _this.OkxKey,
			TradeInst: TradeInst,
			Side:      Side,
			Sz:        Sz,
		})
		// 如果 Sz 大于了最大数量 , 则再来一次
		if mCount.Le(Sz, TradeInst.MaxMktSz) > 0 {
			isAgin = true
			break
		}

		if err != nil {
			err = fmt.Errorf("平仓 下单 失败")
			errArr = append(errArr, err)
		}
	}

	if len(errArr) > 0 {
		resErr = fmt.Errorf("err:%+v", errArr)
	}

	// isAgin 为真，则再来一次
	if isAgin {
		_this.Close()
	}

	return
}
