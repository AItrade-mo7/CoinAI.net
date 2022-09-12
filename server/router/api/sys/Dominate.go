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

	return c.JSON(result.Succeed.WithMsg("认主模式"))
}
