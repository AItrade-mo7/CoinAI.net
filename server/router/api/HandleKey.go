package api

import (
	"CoinAI.net/server/router/result"
	"github.com/EasyGolang/goTools/mFiber"
	"github.com/gofiber/fiber/v2"
)

type HandleKeyParam struct {
	Index    int
	Password string
	Type     string // del //  embed
}

func HandleKey(c *fiber.Ctx) error {
	var json HandleKeyParam
	mFiber.Parser(c, &json)

	// UserID, err := middle.TokenAuth(c)
	// if err != nil {
	// 	return c.JSON(result.ErrToken.WithData(mStr.ToStr(err)))
	// }

	// UserDB, err := dbUser.NewUserDB(dbUser.NewUserOpt{
	// 	UserID: UserID,
	// })
	// if err != nil {
	// 	return c.JSON(result.ErrLogin.WithMsg(err))
	// }
	// // 验证密码
	// err = UserDB.CheckPassword(json.Password)
	// if err != nil {
	// 	return c.JSON(result.ErrLogin.WithMsg(err))
	// }

	// ApiKeyList := config.AppEnv.ApiKeyList

	// NewApiKey := []mOKX.TypeOkxKey{}

	// var ListErr error
	// for key, val := range ApiKeyList {
	// OkxKey := val
	// if key == json.Index {
	// if val.UserID != UserID {
	// 	ListErr = fmt.Errorf("无权操作")
	// 	break
	// }
	// if json.Type == "embed" {
	// 	OkxKey.IsTrade = !OkxKey.IsTrade
	// }
	// if json.Type == "del" {
	// 	continue
	// }
	// }
	// NewApiKey = append(NewApiKey, OkxKey)
	// }

	// if ListErr != nil {
	// 	return c.JSON(result.Fail.WithMsg(ListErr))
	// }

	// config.AppEnv.ApiKeyList = []mOKX.TypeOkxKey{}
	// config.AppEnv.ApiKeyList = NewApiKey

	// global.WriteAppEnv()

	return c.JSON(result.Succeed.WithData("操作完成"))
}
