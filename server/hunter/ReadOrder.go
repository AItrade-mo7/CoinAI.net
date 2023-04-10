package hunter

import (
	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mFile"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mMongo"
	jsoniter "github.com/json-iterator/go"
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

	var OrderArr []okxInfo.VirtualPositionType
	for cur.Next(db.Ctx) {
		var result map[string]any
		cur.Decode(&result)

		var order okxInfo.VirtualPositionType
		jsoniter.Unmarshal(mJson.ToJson(result), &order)
		OrderArr = append(OrderArr, order)
	}

	for i := len(OrderArr) - 1; i >= 0; i-- {
		item := OrderArr[i]
		_this.OrderArr = append(_this.OrderArr, item)
		_this.PositionArr = append(_this.PositionArr, item)
	}

	if len(OrderArr) > 0 {
		_this.NowVirtualPosition = OrderArr[0]

		mJson.Println(_this.NowVirtualPosition)
	}

	mFile.Write(_this.OutPutDirectory+"/OrderArr.json", mJson.ToStr(OrderArr))
}
