package api

import (
	"CoinAI.net/server/okxApi/restApi/account"
	"CoinAI.net/server/router/middle"
	"CoinAI.net/server/router/result"
	"CoinAI.net/server/utils/dbUser"
	"github.com/EasyGolang/goTools/mFiber"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/gofiber/fiber/v2"
)

type SetKeyParam struct {
	Name       string
	ApiKey     string
	SecretKey  string
	Passphrase string
	Password   string
}

func SetKey(c *fiber.Ctx) error {
	var json SetKeyParam
	mFiber.Parser(c, &json)

	if len(json.Name) < 2 {
		return c.JSON(result.Fail.WithMsg("请填写一个备注名"))
	}
	if len(json.ApiKey) < 10 {
		return c.JSON(result.Fail.WithMsg("缺少有效的 API key"))
	}
	if len(json.SecretKey) < 10 {
		return c.JSON(result.Fail.WithMsg("缺少有效的 Secret key"))
	}
	if len(json.Passphrase) < 8 {
		return c.JSON(result.Fail.WithMsg("缺少有效的 密码短语"))
	}

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

	var ApiKey mOKX.TypeOkxKey
	ApiKey.Name = json.Name
	ApiKey.ApiKey = json.ApiKey
	ApiKey.SecretKey = json.SecretKey
	ApiKey.Passphrase = json.Passphrase

	account.GetOKXBalance(ApiKey)

	// 在这里验证Key

	return c.JSON(result.Succeed.WithData("添加一个Key"))
}
