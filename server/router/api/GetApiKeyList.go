package api

import (
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/router/result"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/gofiber/fiber/v2"
)

func GetApiKeyList(c *fiber.Ctx) error {
	return c.JSON(result.Succeed.WithData(GetFuzzyApiKey()))
}

func GetFuzzyApiKey() []mOKX.TypeOkxKey {
	ApiKeyList := config.AppEnv.ApiKeyList

	NewKeyList := []mOKX.TypeOkxKey{}

	for _, val := range ApiKeyList {
		NewKeyList = append(NewKeyList, mOKX.TypeOkxKey{
			Name:       val.Name,
			ApiKey:     GetKeyFuzzy(val.ApiKey, 5),
			SecretKey:  GetKeyFuzzy(val.SecretKey, 5),
			Passphrase: GetKeyFuzzy(val.Passphrase, 1),
			IsTrade:    val.IsTrade,
			UserID:     val.UserID,
		})
	}

	return NewKeyList
}

func GetKeyFuzzy(Ket string, num int) string {
	return Ket[0:num] + "****" + Ket[len(Ket)-num-1:len(Ket)-1]
}
