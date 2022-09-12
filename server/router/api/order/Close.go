package order

import (
	"CoinAI.net/server/okxApi/restApi/order"
	"CoinAI.net/server/router/result"
	"github.com/gofiber/fiber/v2"
)

func Close(c *fiber.Ctx) error {
	order.Close()
	return c.JSON(result.Succeed.WithData("Close"))
}
