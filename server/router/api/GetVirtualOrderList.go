package api

import (
	"CoinAI.net/server/okxInfo"
	"CoinAI.net/server/router/result"
	"github.com/EasyGolang/goTools/mFiber"
	"github.com/gofiber/fiber/v2"
)

type GetVirtualOrderListParam struct {
	HunterName string
}

func GetVirtualOrderList(c *fiber.Ctx) error {
	var json GetVirtualOrderListParam
	mFiber.Parser(c, &json)

	if len(json.HunterName) < 1 {
		return c.JSON(result.Fail.WithMsg("HunterName 不能为空"))
	}

	Hunter := okxInfo.HunterData{}
	for key, item := range okxInfo.NowHunterData {
		if key == json.HunterName {
			Hunter = item
		}
	}
	if len(Hunter.HunterName) < 1 {
		return c.JSON(result.Fail.WithMsg("该Hunter不存在!"))
	}

	return c.JSON(result.Succeed.WithData(json))
}
