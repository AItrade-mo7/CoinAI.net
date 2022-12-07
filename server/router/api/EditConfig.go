package api

import (
	"regexp"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/okxInfo"
	"CoinAI.net/server/router/middle"
	"CoinAI.net/server/router/result"
	"CoinAI.net/server/utils/dbUser"
	"github.com/EasyGolang/goTools/mFiber"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/gofiber/fiber/v2"
)

type EditConfigParam struct {
	Password   string
	ServerName string
	Lever      int
}

func EditConfig(c *fiber.Ctx) error {
	var json EditConfigParam
	mFiber.Parser(c, &json)

	UserID, err := middle.TokenAuth(c)
	if err != nil {
		return c.JSON(result.ErrToken.WithData(mStr.ToStr(err)))
	}

	if UserID != config.AppEnv.UserID {
		return c.JSON(result.Fail.WithMsg("无权操作"))
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

	reg, _ := regexp.Compile("[\u4e00-\u9fa5_a-zA-Z0-9_]{2,12}")
	match := reg.MatchString(json.ServerName)
	if match {
		config.AppEnv.Name = json.ServerName
	} else {
		return c.JSON(result.Fail.WithMsg("系统名称不符合规范!"))
	}

	if json.Lever > 1 && json.Lever < 11 {
		okxInfo.TradeLever = json.Lever
	} else {
		return c.JSON(result.Fail.WithMsg("杠杆系数不符合规范"))
	}


	global.WriteAppEnv()

	return c.JSON(result.Succeed.WithData("操作完成"))
}
