package api

import (
	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/router/result"
	"github.com/EasyGolang/goTools/mFetch"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
)

type MainUserType struct {
	UserID   string `bson:"UserID"`   // 用户 ID
	Email    string `bson:"Email"`    // 用户主要的 Email
	Avatar   string `bson:"Avatar"`   // 用户头像
	NickName string `bson:"NickName"` // 用户昵称
}

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
	AppEnv.ApiKeyList = global.GetFuzzyApiKey()
	ConfigData["AppEnv"] = AppEnv
	ConfigData["GithubInfo"] = GithubInfo
	ConfigData["MainUser"] = GetMainUser()
	// 当前管理员信息
	return c.JSON(result.Succeed.WithData(ConfigData))
}

func GetMainUser() (resData MainUserType) {
	resData = MainUserType{}
	jsoniter.Unmarshal(mJson.ToJson(config.MainUser), &resData)
	return
}
