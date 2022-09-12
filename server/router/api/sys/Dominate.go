package sys

import (
	"CoinAI.net/server/router/result"
	"github.com/EasyGolang/goTools/mFiber"
	"github.com/gofiber/fiber/v2"
)

// 认主模式，后面在写这个逻辑
func Dominate(c *fiber.Ctx) error {
	var json SysAuthParam
	mFiber.Parser(c, &json)

	/*
		0. 确认当前主机是否有主
		x. 确认当前的主机是否绑定Key
		1. 拿到 当前 Token
		2. 提取 UserID
		3. 去数据库 读取用户信息并验证密码
		4. 密码验证通过则验证填写的 Okx 信息

	*/

	return c.JSON(result.Succeed.WithMsg("认主模式"))
}
