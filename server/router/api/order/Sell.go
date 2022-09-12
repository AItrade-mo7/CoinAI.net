package order

import (
	"CoinAI.net/server/okxApi/restApi/order"
	"CoinAI.net/server/router/result"
	"github.com/gofiber/fiber/v2"
)

func Sell(c *fiber.Ctx) error {
	order.Sell()
	return c.JSON(result.Succeed.WithData("Sell"))
}
