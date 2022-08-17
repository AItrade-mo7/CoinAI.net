package ready

import (
	"fmt"
	"time"

	"CoinAI.net/server/analy"
	"CoinAI.net/server/global"
	"CoinAI.net/server/global/dbData"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mClock"
	"github.com/EasyGolang/goTools/mCycle"
	"github.com/EasyGolang/goTools/mOKX"
)

func Start() {
	// 用户和 OkxKey 基本信息
	mCycle.New(mCycle.Opt{
		Func:      SetUserInfo,
		SleepTime: time.Hour * 4, // 每 4 时获取一次
	}).Start()

	SetMarket()
	go mClock.New(mClock.OptType{
		Func: SetMarket,
		Spec: "0 1,16,31,46 0/1 * * ?",
	})
}

func SetUserInfo() {
	GetUserInfo()
	GetOkxKey()
}

func CheckOkx() (resErr error) {
	resErr = nil
	if len(dbData.CoinServe.OkxKeyID) < 10 {
		resErr = fmt.Errorf("读取 dbData.CoinServe 失败 %+v", dbData.CoinServe)
		global.LogErr(resErr)
		return
	}

	if len(dbData.UserInfo.OkxKeyList) < 1 {
		resErr = fmt.Errorf("读取 dbData.UserInfo 失败 %+v", dbData.UserInfo)
		global.LogErr(resErr)
		return
	}

	for _, val := range dbData.UserInfo.OkxKeyList {
		if dbData.CoinServe.OkxKeyID == val.OkxKeyID {
			dbData.OkxKey = val
			break
		}
	}

	if len(dbData.OkxKey.OkxKeyID) < 10 {
		resErr = fmt.Errorf("读取 dbData.OkxKey 失败 %+v", dbData.OkxKey)
		global.LogErr(resErr)
		return
	}

	return
}

func SetMarket() {
	err := CheckOkx()
	if err != nil {
		return
	}

	// 获取市场行情
	GetCoinMarket()

	// 筛选最近币种的信息
	RecentTickerList := analy.RecentTicker()

	//
	okxInfo.AnalyKdata = make(map[string][]mOKX.TypeKd)
	AnalyKdata := make(map[string][]mOKX.TypeKd)
	if len(RecentTickerList) > 3 {
		for _, item := range RecentTickerList {
			list := GetCoinAnalyKdata(item.InstID)
			if len(list) == 300 {
				AnalyKdata[item.InstID] = list
			}
		}
	}

	okxInfo.AnalyKdata = AnalyKdata

	analy.SelectInst(AnalyKdata)
}
