package order

import (
	"CoinAI.net/server/okxApi/restApi/order"
	"CoinAI.net/server/router/result"
	"github.com/gofiber/fiber/v2"
)

func BuySPOT(c *fiber.Ctx) error {
	order.BuySPOT()
	return c.JSON(result.Succeed.WithData("买入现货"))
}
