package hunter

import (
	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/global/dbType"
	"github.com/EasyGolang/goTools/mFile"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mMongo"
	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func (_this *HunterObj) ReadOrder() {
	db, err := mMongo.New(mMongo.Opt{
		UserName: config.SysEnv.MongoUserName,
		Password: config.SysEnv.MongoPassword,
		Address:  config.SysEnv.MongoAddress,
		DBName:   "AIServe",
	}).Connect()
	if err != nil {
		global.LogErr("hunter.ReadOrder 数据库连接失败", _this.HunterName, err)
		return
	}
	defer db.Close()
	db.Collection("CoinOrder")

	findOpt := options.Find()
	findOpt.SetSort(map[string]int{
		"NowTime": -1,
	})
	findOpt.SetAllowDiskUse(true)
	findOpt.SetLimit(100)

	FK := bson.D{}
	FK = append(FK, bson.E{
		Key:   "HunterName",
		Value: _this.HunterName,
	})
	FK = append(FK, bson.E{
		Key:   "ServeID",
		Value: config.AppEnv.ServeID,
	})

	cur, err := db.Table.Find(db.Ctx, FK, findOpt)
	if err != nil {
		global.LogErr("hunter.ReadOrder 数据读取失败", _this.HunterName, err)
	}

	var OrderArr []dbType.VirtualPositionType
	for cur.Next(db.Ctx) {
		var result map[string]any
		cur.Decode(&result)

		var order dbType.VirtualPositionType
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
	}

	// 填充当前除持仓数据
	if len(_this.NowVirtualPosition.InitMoney) < 1 {
		_this.NowVirtualPosition.InitMoney = "1000"
	}

	if len(_this.NowVirtualPosition.ChargeUpl) < 1 {
		_this.NowVirtualPosition.ChargeUpl = "0.05"
	}

	if len(_this.NowVirtualPosition.Money) < 1 {
		_this.NowVirtualPosition.Money = _this.NowVirtualPosition.InitMoney
	}

	mFile.Write(_this.OutPutDirectory+"/OrderArr.json", mJson.ToStr(OrderArr))

	global.TradeLog.Info(_this.HunterName, zap.String("加载初始持仓", mJson.ToStr(_this.NowVirtualPosition)))
}
