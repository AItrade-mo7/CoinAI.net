package order

import (
	"CoinAI.net/server/okxApi/restApi/order"
	"CoinAI.net/server/router/result"
	"github.com/gofiber/fiber/v2"
)

func Buy(c *fiber.Ctx) error {
	order.Buy()
	return c.JSON(result.Succeed.WithData("开多"))
}
