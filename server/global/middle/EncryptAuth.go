package middle

import (
	"errors"
	"strings"

	"CoinAI.net/server/global/config"
	"github.com/EasyGolang/goTools/mEncrypt"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/gofiber/fiber/v2"
)

func EncryptAuth(c *fiber.Ctx) error {
	EncStr := c.Get("Auth-Encrypt")
	if len([]rune(EncStr)) < 20 {
		return errors.New("需要授权码")
	}

	enData := mEncrypt.MD5(mStr.ToStr(c.Body()))
	headersAuth := mStr.Join(c.Path(), c.Get("User-Agent"), enData)
	shaStr := config.Encrypt(headersAuth)
	isFind := strings.Contains(shaStr, EncStr)

	if !isFind {
		return errors.New("授权验证错误")
	}

	return nil
}
