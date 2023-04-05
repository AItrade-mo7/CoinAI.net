package api

import (
	"fmt"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/global/dbType"
	"CoinAI.net/server/global/middle"
	"CoinAI.net/server/okxInfo"
	"CoinAI.net/server/router/result"
	"CoinAI.net/server/utils/dbUser"
	"github.com/EasyGolang/goTools/mFiber"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/gofiber/fiber/v2"
)

type SetAccountConfigParam struct {
	Name       string
	Password   string
	Hunter     string
	TradeLever int
}

func SetAccountConfig(c *fiber.Ctx) error {
	var json SetAccountConfigParam
	mFiber.Parser(c, &json)

	if json.TradeLever < 0 {
		return c.JSON(result.Fail.WithMsg(fmt.Sprintf("TradeLever不可小于 %+v", 0)))
	}

	NowHunter := okxInfo.HunterData{}
	for key, item := range okxInfo.NowHunterData {
		if json.Hunter == key {
			NowHunter = item
			break
		}
	}

	if len(json.Hunter) != 0 {
		if len(NowHunter.HunterName) < 1 {
			return c.JSON(result.Fail.WithMsg(fmt.Sprintf("必须选择一个有效策略 %+v", json.Hunter)))
		}
		if json.TradeLever > NowHunter.MaxTradeLever {
			return c.JSON(result.Fail.WithMsg(fmt.Sprintf("TradeLever不可大于 %+v", NowHunter.MaxTradeLever)))
		}
	}

	// 验证用户和密码
	UserID, err := middle.TokenAuth(c)
	if err != nil {
		return c.JSON(result.ErrToken.WithData(mStr.ToStr(err)))
	}
	UserDB, err := dbUser.NewUserDB(dbUser.NewUserOpt{
		UserID: UserID,
	})
	if err != nil {
		return c.JSON(result.ErrToken.WithData(mStr.ToStr(err)))
	}
	if err != nil {
		UserDB.DB.Close()
		return c.JSON(result.ErrDB.WithData(mStr.ToStr(err)))
	}
	defer UserDB.DB.Close()
	err = UserDB.CheckPassword(json.Password)
	if err != nil {
		return c.JSON(result.ErrDB.WithMsg(mStr.ToStr(err)))
	}
	UserDB.DB.Close()

	// 开始写入数据
	// 寻找Key
	ApiKeyList := config.AppEnv.ApiKeyList
	NewApiKey := []dbType.OkxKeyType{}

	var ListErr error
	for _, val := range ApiKeyList {
		NewKey := val
		if val.Name == json.Name {
			if val.UserID != UserID {
				ListErr = fmt.Errorf("无权操作")
				break
			}
			NewKey.TradeLever = json.TradeLever
			NewKey.Hunter = json.Hunter
		}
		NewApiKey = append(NewApiKey, NewKey)
	}

	if ListErr != nil {
		return c.JSON(result.Fail.WithMsg(ListErr))
	}

	config.AppEnv.ApiKeyList = []dbType.OkxKeyType{}
	config.AppEnv.ApiKeyList = NewApiKey

	global.WriteAppEnv()

	return c.JSON(result.Succeed.WithData("操作完成"))
}
