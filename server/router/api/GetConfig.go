package api

import (
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/router/result"
	"github.com/EasyGolang/goTools/mFetch"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
)

func GetConfig(c *fiber.Ctx) error {
	// 在这里请求数据
	GithubReqData, _ := mFetch.NewHttp(mFetch.HttpOpt{
		Origin: "https://raw.githubusercontent.com",
		Path:   "/AItrade-mo7/CoinAIPackage/main/package.json",
	}).Get()

	var GithubInfo struct {
		Name    string `bson:"name"`
		Version string `bson:"version"`
	}
	jsoniter.Unmarshal(GithubReqData, &GithubInfo)

	ConfigData := make(map[string]any)

	AppEnv := config.AppEnv
	// AppEnv.ApiKeyList = GetFuzzyApiKey()
	ConfigData["AppEnv"] = AppEnv
	ConfigData["GithubInfo"] = GithubInfo
	ConfigData["LeverOpt"] = config.LeverOpt

	return c.JSON(result.Succeed.WithData(ConfigData))
}

func GetFuzzyApiKey() []mOKX.TypeOkxKey {
	// ApiKeyList := config.AppEnv.ApiKeyList

	NewKeyList := []mOKX.TypeOkxKey{}

	// for _, val := range ApiKeyList {
	// 	NewKeyList = append(NewKeyList, mOKX.TypeOkxKey{
	// 		Name:       val.Name,
	// 		ApiKey:     GetKeyFuzzy(val.ApiKey, 5),
	// 		SecretKey:  GetKeyFuzzy(val.SecretKey, 5),
	// 		Passphrase: GetKeyFuzzy(val.Passphrase, 1),
	// 		// IsTrade:    val.IsTrade,
	// 		UserID: val.UserID,
	// 	})
	// }

	return NewKeyList
}

func GetKeyFuzzy(Ket string, num int) string {
	return Ket[0:num] + "****" + Ket[len(Ket)-num-1:len(Ket)-1]
}
