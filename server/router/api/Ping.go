package api

import (
	"CoinServe.net/server/global/config"
	"CoinServe.net/server/router/middle"
	"CoinServe.net/server/router/result"
	"github.com/EasyGolang/goTools/mRes/mFiber"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/gofiber/fiber/v2"
)

func Ping(c *fiber.Ctx) error {
	json := mFiber.Parser(c)

	ReturnData := make(map[string]any)
	ReturnData["ResParam"] = json
	ReturnData["Method"] = c.Method()
	ReturnData["AppInfo"] = config.AppInfo

	ReturnData["UserAgent"] = c.Get("User-Agent")
	ReturnData["FullPath"] = c.BaseURL() + c.OriginalURL()
	ReturnData["ContentType"] = c.Get("Content-Type")

	// 获取 token
	token := c.Get("Token")
	if len(token) > 0 {
		// Token 验证
		_, err := middle.TokenAuth(c)
		if err != nil {
			return c.JSON(result.ErrToken.WithData(mStr.ToStr(err)))
		}
		ReturnData["Token"] = token
	}

	return c.JSON(result.Succeed.WithData(ReturnData))
}
