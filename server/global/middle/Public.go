package middle

import (
	"path"
	"strings"

	"CoinAI.net/server/global/config"
	"CoinAI.net/server/router/result"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/gofiber/fiber/v2"
)

func Public(c *fiber.Ctx) error {
	// 添加访问头
	c.Set("Data-Path", config.SysName)

	findWss := strings.Contains(c.Path(), "/wss")
	if findWss {
		return c.Next()
	}

	filenameWithSuffix := path.Base(c.Path())
	fileSuffix := path.Ext(filenameWithSuffix)
	if len([]rune(fileSuffix)) < 2 { // 后缀名小于2的时候允许验证
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
	}

	isCrawler := CrawlerIS(c)
	if isCrawler {
		return c.JSON(result.ErrLogin.WithData("设备异常"))
	}

	return c.Next()
}
