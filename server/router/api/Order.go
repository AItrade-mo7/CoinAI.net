package api

import (
	"fmt"

	"CoinAI.net/server/global/config"
	"CoinAI.net/server/okxApi"
	"CoinAI.net/server/router/middle"
	"CoinAI.net/server/router/result"
	"CoinAI.net/server/utils/dbUser"
	"github.com/EasyGolang/goTools/mFiber"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/gofiber/fiber/v2"
)

type OrderParam struct {
	Index    int
	Password string
	Type     string
}

func Order(c *fiber.Ctx) error {
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

	OkxKey := config.GetOKXKey(json.Index)

	if UserID != OkxKey.UserID {
		return c.JSON(result.Fail.WithMsg("无权操作"))
	}

	// 新建账户
	OKXAccount, err := okxApi.NewAccount(okxApi.AccountParam{
		OkxKey: OkxKey,
	})
	if err != nil {
		return c.JSON(result.ErrOKXAccount.WithMsg(err))
	}

	// 设置持仓模式
	err = OKXAccount.SetPositionMode()
	if err != nil {
		return c.JSON(result.ErrOKXAccount.WithMsg(err))
	}

	// 设置杠杆倍数
	err = OKXAccount.SetLeverage()
	if err != nil {
		return c.JSON(result.ErrOKXAccount.WithMsg(err))
	}

	// 获取最大可开仓数量
	err = OKXAccount.GetMaxSize()
	if err != nil {
		return c.JSON(result.ErrOKXAccount.WithMsg(err))
	}
	fmt.Println("响应结束", err)

	// if json.Type == "Buy" {
	// 	//
	// }
	// if json.Type == "Sell" {
	// 	//
	// }

	// if json.Type == "Close" {
	// 	//
	// }

	// if json.Type == "BuySPOT" {
	// 	//
	// }

	return c.JSON(result.Succeed.WithData(json.Type))
}
