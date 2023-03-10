package sys

import (
	"CoinAI.net/server/router/result"
	"github.com/EasyGolang/goTools/mFiber"
	"github.com/gofiber/fiber/v2"
)

func Remove(c *fiber.Ctx) error {
	var json SysAuthParam
	mFiber.Parser(c, &json)

	// UserID, err := middle.TokenAuth(c)
	// if err != nil {
	// 	return c.JSON(result.ErrDB.WithData(mStr.ToStr(err)))
	// }
	// if UserID != config.AppEnv.UserID {
	// 	return c.JSON(result.Fail.WithMsg("无权操作"))
	// }

	// UserInfo, err := dbUser.NewUserDB(dbUser.NewUserOpt{
	// 	UserID: UserID,
	// })
	// if err != nil {
	// 	return c.JSON(result.ErrDB.WithData(mStr.ToStr(err)))
	// }

	// err = UserInfo.CheckPassword(json.Password)
	// if err != nil {
	// 	return c.JSON(result.ErrDB.WithMsg(mStr.ToStr(err)))
	// }

	// err = verifyCode.CheckEmailCode(verifyCode.CheckEmailCodeParam{
	// 	Email: UserInfo.AccountData.Email,
	// 	Code:  json.Code,
	// })
	// if err != nil {
	// 	return c.JSON(result.ErrDB.WithMsg(mStr.ToStr(err)))
	// }
	// UserInfo.DB.Close()

	// go shellControl.SysRemove()

	return c.JSON(result.Succeed.WithMsg("指令已发送"))
}
