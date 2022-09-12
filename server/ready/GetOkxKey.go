package ready

import (
	"fmt"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/global/dbType"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mMongo"
	"go.mongodb.org/mongo-driver/bson"
)

func GetOkxKey() {
	db := mMongo.New(mMongo.Opt{
		UserName: config.SysEnv.MongoUserName,
		Password: config.SysEnv.MongoPassword,
		Address:  config.SysEnv.MongoAddress,
		DBName:   "AITrade",
	}).Connect().Collection("CoinServe")

	err := db.Ping()
	if err != nil {
		db.Close()
		errStr := fmt.Errorf("数据读取失败,数据库连接错误 %+v", err)
		global.LogErr(errStr)
		return
	}

	FK := bson.D{{
		Key: "CoinServeID",
	}}

	var result dbType.CoinServeTable
	db.Table.FindOne(db.Ctx, FK).Decode(&result)

	if len(result.CoinServeID) > 6 {
		okxInfo.CoinServe = result
	}
}
