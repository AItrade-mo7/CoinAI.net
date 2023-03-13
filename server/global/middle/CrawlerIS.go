package middle

import (
	"strings"

	"github.com/EasyGolang/goTools/mStr"
	"github.com/EasyGolang/goTools/mVerify"
	"github.com/gofiber/fiber/v2"
)

func CrawlerIS(c *fiber.Ctx) bool {
	Referer := mStr.ToStr(c.Context().Referer())
	RefererPreStr := strings.HasPrefix(Referer, "http")
	RefererSubStr := strings.HasSuffix(Referer, "/")
	DeviceInfo := mVerify.DeviceToUA(c.Get("User-Agent"))
	if len([]rune(Referer)) > 3 && RefererPreStr && RefererSubStr && len(DeviceInfo.BrowserName) > 2 && len(DeviceInfo.OsName) > 2 {
		return false // 不是爬虫
	} else {
		return true // 是爬虫
	}
}
