package dbUser

import (
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/global/dbType"
	"github.com/EasyGolang/goTools/mMongo"
	"go.mongodb.org/mongo-driver/bson"
)

type NewUserOpt struct {
	Email  string
	UserID string
}

type AccountType struct {
	UserID string `bson:"UserID"` // 用户 ID
	Data   dbType.UserTable
	DB     *mMongo.DB
}

func NewUserDB(opt NewUserOpt) (resData *AccountType, resErr error) {
	resData = &AccountType{}
	resErr = nil
	db, err := mMongo.New(mMongo.Opt{
		UserName: config.SysEnv.MongoUserName,
		Password: config.SysEnv.MongoPassword,
		Address:  config.SysEnv.MongoAddress,
		DBName:   "Account",
	}).Connect()
	if err != nil {
		resErr = err
		return
	}
	defer db.Close()
	db.Collection("User")

	resData.DB = db

	FK := bson.D{{
		Key:   "UserEmail",
		Value: opt.Email,
	}}

	if len(opt.UserID) > 0 {
		FK = bson.D{{
			Key:   "UserID",
			Value: opt.UserID,
		}}
	}

	var result dbType.UserTable
	db.Table.FindOne(db.Ctx, FK).Decode(&result)

	resData.UserID = result.UserID
	resData.Data = result

	return
}
