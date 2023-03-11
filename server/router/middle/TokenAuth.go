package middle

import (
	"errors"
	"fmt"
	"time"

	"CoinAI.net/server/global/config"
	"github.com/EasyGolang/goTools/mEncrypt"
	"github.com/EasyGolang/goTools/mMongo"
	"github.com/gofiber/fiber/v2"
)

func TokenAuth(c *fiber.Ctx) (Message string, err error) {
	Message = ""
	err = nil

	Token := c.Get("Token")
	if len(Token) < 1 {
		err = errors.New("缺少Token")
		return
	}

	Claims, AuthOk := mEncrypt.ParseToken(Token, config.SecretKey)
	if !AuthOk {
		err = errors.New("Token验证失败")
		return
	}

	Message = Claims.Message
	UserID := Message
	if len(UserID) != 32 {
		err = errors.New("Token解析失败")
		return
	}

	ExpiresAt := Claims.StandardClaims.ExpiresAt
	now := time.Now().Unix()

	if ExpiresAt-now < 0 {
		err = errors.New("Token过期,请重新登录")
		return
	}

	// 数据库验证
	db := mMongo.New(mMongo.Opt{
		UserName: config.SysEnv.MongoUserName,
		Password: config.SysEnv.MongoPassword,
		Address:  config.SysEnv.MongoAddress,
		DBName:   "AItrade",
	}).Connect().Collection("Token")
	defer db.Close()
	err = db.Ping()
	if err != nil {
		db.Close()
		err = fmt.Errorf("token 验证失败1")
		return
	}

	// var dbRes dbType.TokenTable
	// FK := bson.D{{
	// 	Key:   "UserID",
	// 	Value: UserID,
	// }}
	// db.Table.FindOne(db.Ctx, FK).Decode(&dbRes)

	// if dbRes.Token != Token {
	// 	db.Close()
	// 	err = fmt.Errorf("token 验证失败2")
	// 	return
	// }

	db.Close()

	return
}
