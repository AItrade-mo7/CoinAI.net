package api

import (
	"fmt"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/global/dbType"
	"CoinAI.net/server/global/middle"
	"CoinAI.net/server/router/result"
	"CoinAI.net/server/utils/dbUser"
	"github.com/EasyGolang/goTools/mFiber"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/gofiber/fiber/v2"
)

type HandleKeyParam struct {
	Name     string
	Password string
	Type     string // delete
}

func HandleKey(c *fiber.Ctx) error {
	var json HandleKeyParam
	mFiber.Parser(c, &json)

	var isType bool
	switch json.Type {
	case "disable", "enable", "delete":
		isType = true
	default:
		isType = false
	}
	if !isType {
		return c.JSON(result.Fail.WithMsg("Type类型有问题"))
	}

	// 验证用户和密码
	UserID, err := middle.TokenAuth(c)
	if err != nil {
		return c.JSON(result.ErrToken.WithData(mStr.ToStr(err)))
	}
	UserDB, err := dbUser.NewUserDB(dbUser.NewUserOpt{
		UserID: UserID,
	})
	if err != nil {
		return c.JSON(result.ErrToken.WithData(mStr.ToStr(err)))
	}
	if err != nil {
		UserDB.DB.Close()
		return c.JSON(result.ErrDB.WithData(mStr.ToStr(err)))
	}
	defer UserDB.DB.Close()
	err = UserDB.CheckPassword(json.Password)
	if err != nil {
		return c.JSON(result.ErrDB.WithMsg(mStr.ToStr(err)))
	}
	UserDB.DB.Close()

	// 开始写入数据

	// 寻找Key
	ApiKeyList := config.AppEnv.ApiKeyList
	NewApiKey := []dbType.OkxKeyType{}

	var ListErr error
	for _, val := range ApiKeyList {
		NewKey := val
		if val.Name == json.Name {
			if val.UserID != UserID {
				ListErr = fmt.Errorf("无权操作")
				break
			}
			if json.Type == "delete" {
				continue
			}
		}
		NewApiKey = append(NewApiKey, NewKey)
	}

	if ListErr != nil {
		return c.JSON(result.Fail.WithMsg(ListErr))
	}

	config.AppEnv.ApiKeyList = []dbType.OkxKeyType{}
	config.AppEnv.ApiKeyList = NewApiKey

	global.WriteAppEnv()

	return c.JSON(result.Succeed.WithData("操作完成"))
}
