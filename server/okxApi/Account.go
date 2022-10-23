package okxApi

import (
	"fmt"

	"github.com/EasyGolang/goTools/mOKX"
)

type AccountParam struct {
	OkxKey mOKX.TypeOkxKey
}

type AccountObj struct {
	OkxKey mOKX.TypeOkxKey
}

func NewAccount(opt AccountParam) (resObj *AccountObj, resErr error) {
	obj := AccountObj{}
	resErr = nil

	if !opt.OkxKey.IsTrade {
		resErr = fmt.Errorf("当前 Key 已被禁用: %+v", opt.OkxKey)
		return
	}
	if len(opt.OkxKey.ApiKey) < 10 {
		resErr = fmt.Errorf("需要 ApiKey: %+v", opt.OkxKey)
		return
	}
	obj.OkxKey = opt.OkxKey

	resObj = &obj
	return
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
