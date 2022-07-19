package order

import (
	"CoinAI.net/server/router/result"
	"github.com/gofiber/fiber/v2"
)

func Close(c *fiber.Ctx) error {
	return c.JSON(result.Succeed.WithData("Close"))
}
