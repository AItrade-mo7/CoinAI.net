package api

import (
	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/okxApi/restApi/account"
	"CoinAI.net/server/okxInfo"
	"CoinAI.net/server/router/middle"
	"CoinAI.net/server/router/result"
	"CoinAI.net/server/utils/dbUser"
	"github.com/EasyGolang/goTools/mFiber"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/gofiber/fiber/v2"
)

type SetKeyParam struct {
	Name       string
	ApiKey     string
	SecretKey  string
	Passphrase string
	Password   string
}

func SetKey(c *fiber.Ctx) error {
	var json SetKeyParam
	mFiber.Parser(c, &json)

	if len(config.AppEnv.ApiKeyList) > okxInfo.MaxApiKeyNum {
		return c.JSON(result.Fail.WithMsg("当前服务承载的 ApiKey 已达到上限!"))
	}

	if len(json.Name) < 2 {
		return c.JSON(result.Fail.WithMsg("请填写一个备注名"))
	}
	if len(json.ApiKey) < 10 {
		return c.JSON(result.Fail.WithMsg("缺少有效的 API key"))
	}
	if len(json.SecretKey) < 10 {
		return c.JSON(result.Fail.WithMsg("缺少有效的 Secret key"))
	}
	if len(json.Passphrase) < 8 {
		return c.JSON(result.Fail.WithMsg("缺少有效的 密码短语"))
	}

	UserID, err := middle.TokenAuth(c)
	if err != nil {
		return c.JSON(result.ErrToken.WithData(mStr.ToStr(err)))
	}

	UserDB, err := dbUser.NewUserDB(dbUser.NewUserOpt{
		UserID: UserID,
	})
	if err != nil {
		return c.JSON(result.ErrLogin.WithMsg(err))
	}
	// 验证密码
	err = UserDB.CheckPassword(json.Password)
	if err != nil {
		return c.JSON(result.ErrLogin.WithMsg(err))
	}

	var ApiKey mOKX.TypeOkxKey
	ApiKey.Name = json.Name
	ApiKey.ApiKey = json.ApiKey
	ApiKey.SecretKey = json.SecretKey
	ApiKey.Passphrase = json.Passphrase
	ApiKey.IsTrade = true
	ApiKey.UserID = UserID

	// 验证 Key 准确性
	_, err = account.GetOKXBalance(ApiKey)
	if err != nil {
		return c.JSON(result.ErrLogin.WithMsg("Api Key 验证失败!"))
	}

	ApiKeyList := config.AppEnv.ApiKeyList
	isRepeat := false
	isName := false
	for _, val := range ApiKeyList {

		if ApiKey.ApiKey == val.ApiKey || ApiKey.SecretKey == val.SecretKey {
			isRepeat = true
			break
		}
		if ApiKey.Name == val.Name {
			isName = true
			break
		}
	}
	if isName {
		return c.JSON(result.ErrLogin.WithMsg("备注名重复"))
	}

	if isRepeat {
		return c.JSON(result.ErrLogin.WithMsg("Api key已存在"))
	}

	config.AppEnv.ApiKeyList = append(config.AppEnv.ApiKeyList, ApiKey)

	global.WriteAppEnv()

	return c.JSON(result.Succeed.WithData("添加一个Key"))
}
