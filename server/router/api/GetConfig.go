package api

import (
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/router/result"
	"github.com/EasyGolang/goTools/mFetch"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
)

func GetConfig(c *fiber.Ctx) error {
	// 在这里请求数据
	GithubReqData := mFetch.NewHttp(mFetch.HttpOpt{
		Origin: "https://raw.githubusercontent.com",
		Path:   "/mo7static/CoinAI.net/main/package.json",
	}).Get()

	var GithubInfo struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}
	jsoniter.Unmarshal(GithubReqData, &GithubInfo)

	ConfigData := make(map[string]any)
	ConfigData["SysEnv"] = config.SysEnv
	ConfigData["AppEnv"] = config.AppEnv
	ConfigData["AppInfo"] = config.AppInfo
	ConfigData["GithubInfo"] = GithubInfo

	return c.JSON(result.Succeed.WithData(ConfigData))
}
