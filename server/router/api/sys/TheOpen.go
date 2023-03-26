package sys

import (
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

	if len(json.Password) < 1 {
		return c.JSON(result.Fail.WithMsg("需要密码"))
	}
	if len(json.Code) < 1 {
		return c.JSON(result.Fail.WithMsg("需要验证码"))
	}
	if len(json.Password) != 32 {
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
	/*
		开始填写业务逻辑
		1. 如果为 非 public 则 公开
		2. 如果为 公开 的 则判断 Key 数量，再选择关闭
	*/

	if config.AppEnv.IsPublic {
		config.AppEnv.IsPublic = false

		return c.JSON(result.Succeed.WithMsg("隐藏卫星服务"))
	} else {
		config.AppEnv.IsPublic = true

		return c.JSON(result.Succeed.WithMsg("公开卫星服务"))
	}
}
