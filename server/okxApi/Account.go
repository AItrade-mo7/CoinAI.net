package okxApi

import "github.com/EasyGolang/goTools/mOKX"

type AccountParam struct {
	OkxKey mOKX.TypeOkxKey
}

type AccountObj struct {
	OkxKey mOKX.TypeOkxKey
}

func NewAccount(opt AccountParam) *AccountObj {
	obj := AccountObj{}

	return &obj
}

// 设置持仓模式
func (_this *AccountObj) SetPositionMode() {
}

// 设置杠杆倍数
func (_this *AccountObj) SetLeverage() {
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
