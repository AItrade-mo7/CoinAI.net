package ready

import (
	"fmt"
	"time"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/dbData"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mClock"
	"github.com/EasyGolang/goTools/mCycle"
)

func Start() {
	mCycle.New(mCycle.Opt{
		Func:      GetUserInfo,
		SleepTime: time.Hour * 4, // 每 4 时获取一次
	}).Start()

	CheckUserData()

	SetMarket()

	go mClock.New(mClock.OptType{
		Func: SetMarket,
		Spec: "0 1,16,31,46 0/1 * * ?",
	})
}

func SetMarket() {
	GetCoinMarket()
	// 一旦有一个长度不对，则 Market 不合格
	if len(okxInfo.Unit) < 3 || len(okxInfo.TickerList) < 4 || len(okxInfo.AnalyWhole) < 4 || len(okxInfo.AnalySingle) < 4 {
		okxInfo.IsMarket = false
	} else {
		okxInfo.IsMarket = true
	}
}

// 获取基本的用户数据和Key数据
func CheckUserData() {
	GetOkxKey()
	if len(dbData.CoinServe.OkxKeyID) < 10 {
		errStr := fmt.Errorf("读取 dbData.CoinServe 失败 %+v", dbData.CoinServe)
		global.LogErr(errStr)
		panic(errStr)
	}

	if len(dbData.UserInfo.OkxKeyList) < 1 {
		errStr := fmt.Errorf("读取 dbData.UserInfo 失败 %+v", dbData.UserInfo)
		global.LogErr(errStr)
		panic(errStr)
	}

	for _, val := range dbData.UserInfo.OkxKeyList {
		if dbData.CoinServe.OkxKeyID == val.OkxKeyID {
			dbData.OkxKey = val
			break
		}
	}

	if len(dbData.OkxKey.OkxKeyID) < 10 {
		errStr := fmt.Errorf("读取 dbData.OkxKey 失败 %+v", dbData.OkxKey)
		global.LogErr(errStr)
		panic(errStr)
	}
}
