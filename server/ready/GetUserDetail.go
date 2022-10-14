package ready

import (
	"fmt"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/utils/dbUser"
)

func GetUserInfo() {
	UserID := config.AppEnv.UserID

	UserDB, err := dbUser.NewUserDB(dbUser.NewUserOpt{
		UserID: UserID,
	})
	if err != nil {
		UserDB.DB.Close()
		errStr := fmt.Errorf("用户数据读取错误 %+v", err)
		global.LogErr(errStr)
		return
	}
}
