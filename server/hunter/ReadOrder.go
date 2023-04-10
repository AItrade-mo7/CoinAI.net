package hunter

import (
	"fmt"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/global/dbType"
	"github.com/EasyGolang/goTools/mFile"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mMongo"
	"github.com/EasyGolang/goTools/mTime"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (_this *HunterObj) ReadOrder() {
	db := mMongo.New(mMongo.Opt{
		UserName: config.SysEnv.MongoUserName,
		Password: config.SysEnv.MongoPassword,
		Address:  config.SysEnv.MongoAddress,
		DBName:   "AIServe",
	}).Connect().Collection("CoinOrder")
	defer db.Close()
	findOpt := options.Find()
	findOpt.SetSort(map[string]int{
		"NowTime": -1,
	})
	findOpt.SetAllowDiskUse(true)
	findOpt.SetLimit(100)

	FK := bson.D{{
		Key:   "HunterName",
		Value: _this.HunterName,
	}}

	cur, err := db.Table.Find(db.Ctx, FK, findOpt)
	if err != nil {
		global.LogErr("hunter.ReadOrder 数据读取失败", _this.HunterName, err)
	}

	var CoinOrder []dbType.CoinOrderTable
	for cur.Next(db.Ctx) {
		var result dbType.CoinOrderTable
		cur.Decode(&result)
		CoinOrder = append(CoinOrder, result)
		fmt.Println(mTime.TimeGet(result.CreateTime).TimeStr)
	}

	mFile.Write(_this.OutPutDirectory+"/CoinOrder.json", mJson.ToStr(CoinOrder))
}
