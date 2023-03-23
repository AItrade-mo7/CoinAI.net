package api

import (
	"fmt"

	"CoinAI.net/server/global/config"
	"CoinAI.net/server/global/dbType"
	"CoinAI.net/server/global/middle"
	"CoinAI.net/server/okxApi"
	"CoinAI.net/server/router/result"
	"CoinAI.net/server/utils/dbUser"
	"github.com/EasyGolang/goTools/mFiber"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/gofiber/fiber/v2"
)

type OrderParam struct {
	Name     string
	Password string
	Type     string // Buy   Sell   Close
}

func Order(c *fiber.Ctx) error {
	var json OrderParam
	mFiber.Parser(c, &json)

	var isType bool
	switch json.Type {
	case "Buy", "Sell", "Close":
		isType = true
	default:
		isType = false
	}
	if !isType {
		return c.JSON(result.Fail.WithMsg("Type类型有问题"))
	}

	UserID, err := middle.TokenAuth(c)
	if err != nil {
		return c.JSON(result.ErrToken.WithData(mStr.ToStr(err)))
	}

	// 验证密码
	UserDB, err := dbUser.NewUserDB(dbUser.NewUserOpt{
		UserID: UserID,
	})
	if err != nil {
		UserDB.DB.Close()
		return c.JSON(result.ErrDB.WithData(mStr.ToStr(err)))
	}
	defer UserDB.DB.Close()
	err = UserDB.CheckPassword(json.Password)
	if err != nil {
		return c.JSON(result.ErrDB.WithMsg(mStr.ToStr(err)))
	}

	// 寻找 Key
	ApiKeyList := config.AppEnv.ApiKeyList
	var OkxKey dbType.OkxKeyType
	var ListErr error
	for _, val := range ApiKeyList {
		if val.Name == json.Name {
			if val.UserID != UserID {
				ListErr = fmt.Errorf("无权操作")
				break
			}
			OkxKey = val
			break
		}
	}
	if ListErr != nil {
		return c.JSON(result.Fail.WithMsg(ListErr))
	}

	// 新建账户
	OKXAccount, err := okxApi.NewAccount(okxApi.AccountParam{
		OkxKey: OkxKey,
	})
	if err != nil {
		return c.JSON(result.ErrOKXAccount.WithMsg(err))
	}

	if json.Type == "Buy" {
		err = OKXAccount.Buy()
	}
	if json.Type == "Sell" {
		err = OKXAccount.Sell()
	}
	if json.Type == "Close" {
		err = OKXAccount.Close()
	}

	if err != nil {
		return c.JSON(result.ErrOKXAccount.WithMsg(err))
	}

	return c.JSON(result.Succeed.WithData(json.Type))
}
