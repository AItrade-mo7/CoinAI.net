package ready

import (
	"fmt"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/dbData"
)

func Start() {
	GetUserInfo()
	GetOkxKey()

	if len(dbData.CoinServe.OkxKeyID) < 5 {
		errStr := fmt.Errorf("读取 dbData.CoinServe 失败 %+v", dbData.CoinServe)
		global.LogErr(errStr)
		panic(errStr)
	}

	if len(dbData.UserInfo.UserID) < 5 {
		errStr := fmt.Errorf("读取 dbData.UserInfo 失败 %+v", dbData.UserInfo)
		global.LogErr(errStr)
		panic(errStr)
	}
}
