package order

import (
	"fmt"

	"CoinAI.net/server/global/config"
	"CoinAI.net/server/okxApi/restApi/order"
	"CoinAI.net/server/router/middle"
	"CoinAI.net/server/router/result"
	"CoinAI.net/server/utils/dbUser"
	"github.com/EasyGolang/goTools/mFiber"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/gofiber/fiber/v2"
)

func Buy(c *fiber.Ctx) error {
	var json OrderParam
	mFiber.Parser(c, &json)

	UserID, err := middle.TokenAuth(c)
	if err != nil {
		return c.JSON(result.ErrToken.WithData(mStr.ToStr(err)))
	}

	UserDB, err := dbUser.NewUserDB(dbUser.NewUserOpt{
		UserID: UserID,
	})
	if err != nil {
		return c.JSON(result.ErrLogin.WithMsg(err))
	}
	// 验证密码
	err = UserDB.CheckPassword(json.Password)
	if err != nil {
		return c.JSON(result.ErrLogin.WithMsg(err))
	}

	ApiKeyList := config.AppEnv.ApiKeyList
	var ListErr error
	ApiKey := mOKX.TypeOkxKey{}
	for key, val := range ApiKeyList {
		if key == json.Index {
			if val.UserID != UserID {
				ListErr = fmt.Errorf("无权操作")
				break
			}
			ApiKey = val
		}
	}

	if ListErr != nil {
		return c.JSON(result.Fail.WithMsg(ListErr))
	}

	err = order.Buy(ApiKey)
	if err != nil {
		return c.JSON(result.Fail.WithData(err))
	}
	return c.JSON(result.Succeed.WithData("做多"))
}
