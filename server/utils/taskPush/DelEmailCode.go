package taskPush

import (
	"fmt"

	"CoinAI.net/server/global/config"
	"github.com/EasyGolang/goTools/mMongo"
	"github.com/EasyGolang/goTools/mVerify"
	"go.mongodb.org/mongo-driver/bson"
)

// 删除指定邮箱当前的验证码
func DelEmailCode(email string) error {
	isEmail := mVerify.IsEmail(email)
	if !isEmail {
		emailErr := fmt.Errorf("json.Email 格式不正确 %+v", email)
		return emailErr
	}

	db := mMongo.New(mMongo.Opt{
		UserName: config.SysEnv.MongoUserName,
		Password: config.SysEnv.MongoPassword,
		Address:  config.SysEnv.MongoAddress,
		DBName:   "Message",
	}).Connect().Collection("VerifyCode")
	defer db.Close()
	// 查找参数设置
	FK := bson.D{{
		Key:   "Email",
		Value: email,
	}}
	_, err := db.Table.DeleteOne(db.Ctx, FK)
	db.Close()
	if err != nil {
		return err
	}

	return nil
}
