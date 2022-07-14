package middle

import (
	"strings"

	"Hunter.net/server/router/result"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/gofiber/fiber/v2"
)

func Public(c *fiber.Ctx) error {
	// 添加访问头
	AddHeader(c)

	findWss := strings.Contains(c.Path(), "/wss")
	if findWss {
		return c.Next()
	}

	// 授权验证
	err := EncryptAuth(c)
	if err != nil {
		return c.JSON(result.ErrAuth.WithData(mStr.ToStr(err)))
	}

	// Token 验证
	_, err = TokenAuth(c)
	if err != nil {
		return c.JSON(result.ErrToken.WithData(mStr.ToStr(err)))
	}

	return c.Next()
}

func AddHeader(c *fiber.Ctx) error {
	c.Set("Data-Path", "Hunter.net")
	return nil
}
