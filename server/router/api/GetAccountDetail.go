package api

import (
	"CoinAI.net/server/router/result"
	"github.com/gofiber/fiber/v2"
)

func GetAccountDetail(c *fiber.Ctx) error {
	// 在这里请求数据

	return c.JSON(result.Succeed.WithData("接口开发中"))
}
