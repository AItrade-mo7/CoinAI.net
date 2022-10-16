package api

import (
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/router/result"
	"github.com/EasyGolang/goTools/mFetch"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
)

func SetKey(c *fiber.Ctx) error {
	// 在这里请求数据
	GithubReqData, _ := mFetch.NewHttp(mFetch.HttpOpt{
		Origin: "https://raw.githubusercontent.com",
		Path:   "/AITrade-mo7/CoinAIPackage/main/package.json",
	}).Get()

	var GithubInfo struct {
		Name    string `bson:"name"`
		Version string `bson:"version"`
	}
	jsoniter.Unmarshal(GithubReqData, &GithubInfo)

	ConfigData := make(map[string]any)
	ConfigData["AppEnv"] = config.AppEnv
	ConfigData["GithubInfo"] = GithubInfo

	return c.JSON(result.Succeed.WithData(ConfigData))
}
