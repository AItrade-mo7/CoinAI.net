package api

import (
	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/global/dbType"
	"CoinAI.net/server/global/middle"
	"CoinAI.net/server/okxApi/restApi/account"
	"CoinAI.net/server/router/result"
	"CoinAI.net/server/utils/dbUser"
	"github.com/EasyGolang/goTools/mFiber"
	"github.com/EasyGolang/goTools/mMongo"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/EasyGolang/goTools/mVerify"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	if len(config.AppEnv.ApiKeyList) > config.AppEnv.MaxApiKeyNum {
		return c.JSON(result.Fail.WithMsg("当前服务承载的 ApiKey 已达到上限!"))
	}

	isName := mVerify.IsNickName(json.Name)
	if !isName {
		return c.JSON(result.Fail.WithMsg("备注名不规范"))
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

	var ApiKey dbType.OkxKeyType
	ApiKey.Name = json.Name
	ApiKey.ApiKey = json.ApiKey
	ApiKey.SecretKey = json.SecretKey
	ApiKey.Passphrase = json.Passphrase
	ApiKey.UserID = UserID

	// 验证 Key 可用性
	_, err = account.GetOKXBalance(mOKX.TypeOkxKey{
		ApiKey:     ApiKey.ApiKey,
		SecretKey:  ApiKey.SecretKey,
		Passphrase: ApiKey.Passphrase,
	})
	if err != nil {
		return c.JSON(result.ErrLogin.WithMsg("Api Key 验证失败!"))
	}

	// 在这里检查 ApiKey 是否重复
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
		Key:   "ApiKeyList.ApiKey",
		Value: ApiKey.ApiKey,
	}}
	var dbAPPEnvKey dbType.AppEnvType
	db.Table.FindOne(db.Ctx, FK, findOpt).Decode(&dbAPPEnvKey)
	if len(dbAPPEnvKey.ServeID) > 0 {
		return c.JSON(result.ErrLogin.WithMsg("ApiKey已存在!!!"))
	}
	// 检查API的名字是否重复
	FK = bson.D{{
		Key:   "ApiKeyList.Name",
		Value: ApiKey.Name,
	}}
	var dbAPPEnvName dbType.AppEnvType
	db.Table.FindOne(db.Ctx, FK, findOpt).Decode(&dbAPPEnvName)
	if len(dbAPPEnvName.ServeID) > 0 {
		return c.JSON(result.ErrLogin.WithMsg("备注名已存在!"))
	}

	config.AppEnv.ApiKeyList = append(config.AppEnv.ApiKeyList, ApiKey)

	global.WriteAppEnv()

	return c.JSON(result.Succeed.WithData("成功添加"))
}
