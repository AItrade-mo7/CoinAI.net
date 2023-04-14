package sys

import (
	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/global/middle"
	"CoinAI.net/server/router/result"
	"CoinAI.net/server/utils/dbUser"
	"CoinAI.net/server/utils/taskPush"
	"github.com/EasyGolang/goTools/mFiber"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/gofiber/fiber/v2"
)

type TheOpenParam struct {
	Password string
	Code     string
}

func TheOpen(c *fiber.Ctx) error {
	var json SysAuthParam
	mFiber.Parser(c, &json)

	UserID, err := middle.TokenAuth(c)
	if err != nil {
		return c.JSON(result.ErrDB.WithData(mStr.ToStr(err)))
	}
	if UserID != config.MainUser.UserID {
		return c.JSON(result.Fail.WithMsg("无权操作"))
	}

	if len(json.Code) < 1 {
		return c.JSON(result.Fail.WithMsg("需要验证码"))
	}
	if len(json.Password) < 16 {
		return c.JSON(result.ErrLogin.With("密码格式不正确", "可能原因:密码没有加密传输！"))
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

	// 验证邮箱验证码
	err = taskPush.CheckEmailCode(taskPush.CheckEmailCodeParam{
		Email: UserDB.Data.Email,
		Code:  json.Code,
	})
	if err != nil {
		return c.JSON(result.Fail.WithMsg(err))
	}

	UserDB.DB.Close()

	config.AppEnv.IsPublic = !config.AppEnv.IsPublic
	global.WriteAppEnv()

	taskPush.DelEmailCode(UserDB.Data.Email)

	if config.AppEnv.IsPublic {
		return c.JSON(result.Succeed.WithMsg("已经公开此服务"))
	} else {
		return c.JSON(result.Succeed.WithMsg("次服务已隐藏"))
	}
}
