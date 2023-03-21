package api

import (
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/global/dbType"
	"CoinAI.net/server/router/result"
	"github.com/EasyGolang/goTools/mFetch"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
)

func GetConfig(c *fiber.Ctx) error {
	// 在这里请求数据
	GithubReqData, _ := mFetch.NewHttp(mFetch.HttpOpt{
		Origin: config.GithubPackagePath.Origin,
		Path:   config.GithubPackagePath.Path,
	}).Get()

	var GithubInfo struct {
		Name    string `bson:"name"`
		Version string `bson:"version"`
	}
	jsoniter.Unmarshal(GithubReqData, &GithubInfo)

	ConfigData := make(map[string]any)

	AppEnv := config.AppEnv
	AppEnv.ApiKeyList = GetFuzzyApiKey()
	ConfigData["AppEnv"] = AppEnv
	ConfigData["GithubInfo"] = GithubInfo

	return c.JSON(result.Succeed.WithData(ConfigData))
}

func GetFuzzyApiKey() []dbType.OkxKeyType {
	ApiKeyList := config.AppEnv.ApiKeyList

	NewKeyList := []dbType.OkxKeyType{}

	for _, val := range ApiKeyList {
		NewKeyList = append(NewKeyList, dbType.OkxKeyType{
			Name:       val.Name,
			ApiKey:     mStr.GetKeyFuzzy(val.ApiKey, 3),
			SecretKey:  mStr.GetKeyFuzzy(val.SecretKey, 3),
			Passphrase: mStr.GetKeyFuzzy(val.Passphrase, 1),
			UserID:     val.UserID,
		})
	}

	return NewKeyList
}
