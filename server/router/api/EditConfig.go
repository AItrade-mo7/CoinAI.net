package api

import (
	"CoinAI.net/server/router/result"
	"github.com/EasyGolang/goTools/mFiber"
	"github.com/gofiber/fiber/v2"
)

type EditConfigParam struct {
	Password     string
	ServerName   string
	Lever        int
	MaxApiKeyNum int
}

func EditConfig(c *fiber.Ctx) error {
	var json EditConfigParam
	mFiber.Parser(c, &json)

	// UserID, err := middle.TokenAuth(c)
	// if err != nil {
	// 	return c.JSON(result.ErrToken.WithData(mStr.ToStr(err)))
	// }

	// if UserID != config.AppEnv.UserID {
	// 	return c.JSON(result.Fail.WithMsg("无权操作"))
	// }

	// UserDB, err := dbUser.NewUserDB(dbUser.NewUserOpt{
	// 	UserID: UserID,
	// })
	// if err != nil {
	// 	return c.JSON(result.ErrLogin.WithMsg(err))
	// }
	// // 验证密码
	// err = UserDB.CheckPassword(json.Password)
	// if err != nil {
	// 	return c.JSON(result.ErrLogin.WithMsg(err))
	// }

	// reg, _ := regexp.Compile("[\u4e00-\u9fa5_a-zA-Z0-9_]{2,12}")
	// match := reg.MatchString(json.ServerName)
	// if match {
	// 	// config.AppEnv.Name = json.ServerName
	// } else {
	// 	return c.JSON(result.Fail.WithMsg("系统名称不符合规范!"))
	// }

	// AppEnv := config.AppEnv
	// if json.Lever >= config.LeverOpt[0] && json.Lever <= config.LeverOpt[len(config.LeverOpt)-1] {
	// 	// AppEnv.TradeLever = json.Lever
	// } else {
	// 	return c.JSON(result.Fail.WithMsg("杠杆系数不符合规范"))
	// }

	// if json.MaxApiKeyNum > 1 && json.MaxApiKeyNum > len(config.AppEnv.ApiKeyList) {
	// 	AppEnv.MaxApiKeyNum = json.MaxApiKeyNum
	// } else {
	// 	return c.JSON(result.Fail.WithMsg("最大 ApiKey 数量不正确"))
	// }

	// config.AppEnv = AppEnv
	// global.WriteAppEnv()

	return c.JSON(result.Succeed.WithData("操作完成"))
}
