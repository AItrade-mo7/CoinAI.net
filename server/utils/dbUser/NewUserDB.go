package dbUser

import (
	"fmt"

	"Hunter.net/server/global"
	"Hunter.net/server/global/config"
	"Hunter.net/server/global/dbType"
	"github.com/EasyGolang/goTools/mMongo"
	"go.mongodb.org/mongo-driver/bson"
)

type NewUserOpt struct {
	Email  string
	UserID string
}

type AccountType struct {
	UserID      string `bson:"UserID"`
	AccountData dbType.AccountTable
	DB          *mMongo.DB
}

func NewUserDB(opt NewUserOpt) (resData *AccountType, resErr error) {
	resData = &AccountType{}
	resErr = nil

	db := mMongo.New(mMongo.Opt{
		UserName: config.SysEnv.MongoUserName,
		Password: config.SysEnv.MongoPassword,
		Address:  config.SysEnv.MongoAddress,
		DBName:   "Hunter",
	}).Connect().Collection("Account")

	resData.DB = db

	err := db.Ping()
	if err != nil {
		db.Close()
		errStr := fmt.Errorf("用户数据读取失败,数据库连接错误 %+v", err)
		global.LogErr(errStr)
		resErr = errStr
		return
	}

	var result dbType.AccountTable
	FK := bson.D{{
		Key:   "Email",
		Value: opt.Email,
	}}
	if len(opt.UserID) > 3 {
		FK = bson.D{{
			Key:   "UserID",
			Value: opt.UserID,
		}}
	}
	db.Table.FindOne(db.Ctx, FK).Decode(&result)

	resData.UserID = result.UserID
	resData.AccountData = result

	return
}
