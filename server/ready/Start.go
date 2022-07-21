package ready

import (
	"fmt"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/dbData"
)

func Start() {
	GetUserInfo()
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
