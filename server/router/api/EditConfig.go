package api

import (
	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/global/dbType"
	"CoinAI.net/server/global/middle"
	"CoinAI.net/server/router/result"
	"CoinAI.net/server/utils/dbUser"
	"github.com/EasyGolang/goTools/mFiber"
	"github.com/EasyGolang/goTools/mMongo"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/EasyGolang/goTools/mVerify"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EditConfigParam struct {
	Password     string
	SysName      string
	MaxApiKeyNum int
}

func EditConfig(c *fiber.Ctx) error {
	var json EditConfigParam
	mFiber.Parser(c, &json)

	UserID, err := middle.TokenAuth(c)
	if err != nil {
		return c.JSON(result.ErrDB.WithData(mStr.ToStr(err)))
	}
	if UserID != config.MainUser.UserID {
		return c.JSON(result.Fail.WithMsg("无权操作"))
	}
	if len(json.Password) < 1 {
		return c.JSON(result.Fail.WithMsg("需要密码"))
	}
	if len(json.SysName) < 1 {
		return c.JSON(result.Fail.WithMsg("需要密码"))
	}

	if json.MaxApiKeyNum < len(config.AppEnv.ApiKeyList) {
		return c.JSON(result.Fail.WithMsg("ApiKey数量错误"))
	}

	isName := mVerify.IsNickName(json.SysName)
	if !isName {
		return c.JSON(result.Fail.WithMsg("系统名称不符合规范!"))
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

	// 检查是否修改了服务器名字
	if config.AppEnv.SysName != json.SysName {
		// 检查名称是否重复
		db := mMongo.New(mMongo.Opt{
			UserName: config.SysEnv.MongoUserName,
			Password: config.SysEnv.MongoPassword,
			Address:  config.SysEnv.MongoAddress,
			DBName:   "AIServe",
		}).Connect().Collection("CoinAI")
		defer db.Close()
		findOpt := options.FindOne()
		findOpt.SetSort(map[string]int{
			"TimeUnix": -1,
		})
		FK := bson.D{{
			Key:   "SysName",
			Value: json.SysName,
		}}
		var DBAppEnv dbType.AppEnvType
		db.Table.FindOne(db.Ctx, FK, findOpt).Decode(&DBAppEnv)
		if len(DBAppEnv.ServeID) > 3 {
			c.JSON(result.Succeed.WithData("系统名称重复!"))
		}
		config.AppEnv.SysName = json.SysName
	}

	config.AppEnv.MaxApiKeyNum = json.MaxApiKeyNum

	global.WriteAppEnv()

	return c.JSON(result.Succeed.WithData("操作完成"))
}
