package api

import (
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/global/middle"
	"CoinAI.net/server/router/result"
	"github.com/EasyGolang/goTools/mFiber"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/EasyGolang/goTools/mVerify"
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

	DeviceInfo := mVerify.DeviceToUA(c.Get("User-Agent"))
	ReturnData["BrowserName"] = DeviceInfo.BrowserName
	ReturnData["OsName"] = DeviceInfo.OsName

	ips := c.IPs()
	if len(ips) > 0 {
		ReturnData["IP"] = ips[0]
	}

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
