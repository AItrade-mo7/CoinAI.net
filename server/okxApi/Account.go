package okxApi

import (
	"fmt"

	"CoinAI.net/server/okxApi/restApi/account"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mOKX"
)

type AccountParam struct {
	OkxKey mOKX.TypeOkxKey
}

type AccountObj struct {
	OkxKey     mOKX.TypeOkxKey
	TradeInst  mOKX.TypeTicker // 交易币种信息
	TradeLever int             // 杠杆倍数
}

// 创建一个新账户
func NewAccount(opt AccountParam) (resObj *AccountObj, resErr error) {
	obj := AccountObj{}
	resErr = nil

	if !opt.OkxKey.IsTrade {
		resErr = fmt.Errorf("当前 Key 已被禁用")
		return
	}
	if len(opt.OkxKey.ApiKey) < 10 {
		resErr = fmt.Errorf("ApiKey 不能为空 ")
		return
	}

	if len(okxInfo.TradeInst.InstID) < 3 {
		resErr = fmt.Errorf("okxInfo.TradeInst.InstID 不能为空 %+v", okxInfo.TradeInst.InstID)
		return
	}

	if (okxInfo.TradeLever) < 1 {
		resErr = fmt.Errorf("okxInfo.TradeLever 不能为空 %+v", okxInfo.TradeLever)
		return
	}

	obj.OkxKey = opt.OkxKey
	obj.TradeInst = okxInfo.TradeInst
	obj.TradeLever = okxInfo.TradeLever

	resObj = &obj
	return
}

// 设置持仓模式
func (_this *AccountObj) SetPositionMode() (resErr error) {
	resErr = account.SetPositionMode(_this.OkxKey)
	return
}

// 设置杠杆倍数
func (_this *AccountObj) SetLeverage() (resErr error) {
	resErr = account.SetLeverage(account.SetLeverageParam{
		InstID: _this.TradeInst.InstID,
		OKXKey: _this.OkxKey,
		Lever:  _this.TradeLever,
	})
	return
}

// 获取账户余额
func (_this *AccountObj) GetBalance() {
}

// 获取持仓信息
func (_this *AccountObj) GetPositions() {
}

// 获取最大可开仓数量
func (_this *AccountObj) GetMaxSize() {
}
