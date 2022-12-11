package backTest

import (
	"fmt"
	"time"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mMongo"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mTime"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 任意时间点之间的回测函数

type TestOpt struct {
	StartTime int64
	EndTime   int64
	CcyName   string
}

type TestObj struct {
	StartTime int64
	EndTime   int64
	TableName string
	InstID    string
}

func NewTest(opt TestOpt) *TestObj {
	obj := TestObj{}
	now := mTime.GetUnixInt64()

	if opt.EndTime < now-mTime.UnixTimeInt64.Day*2190 {
		opt.EndTime = now
	}

	obj.StartTime = opt.StartTime
	obj.EndTime = opt.EndTime

	obj.TableName = opt.CcyName + "USDT"
	obj.InstID = opt.CcyName + "-USDT"

	return &obj
}

func (_this *TestObj) GetDBKdata() *TestObj {
	total := (_this.EndTime - _this.StartTime) / mTime.UnixTimeInt64.Hour

	Timeout := int(total) * 10
	if Timeout < 100 {
		Timeout = 100
	}

	db := mMongo.New(mMongo.Opt{
		UserName: config.SysEnv.MongoUserName,
		Password: config.SysEnv.MongoPassword,
		Address:  config.SysEnv.MongoAddress,
		DBName:   "AITrade",
		Timeout:  Timeout,
	}).Connect().Collection(_this.TableName)
	defer global.RunLog.Println("关闭数据库连接", _this.TableName)
	defer db.Close()

	fmt.Println(config.SysEnv)
	fmt.Println(_this.TableName)
	fmt.Println(Timeout)

	findOpt := options.Find()
	findOpt.SetSort(map[string]int{
		"TimeUnix": 1,
	})
	findOpt.SetAllowDiskUse(true)

	FK := bson.D{}
	FK = append(FK, bson.E{
		Key: "TimeUnix",
		Value: bson.D{
			{
				Key:   "$gte", // 大于或等于
				Value: _this.StartTime,
			}, {
				Key:   "$lte", // 小于或等于
				Value: _this.EndTime,
			},
		},
	})

	cur, err := db.Table.Find(db.Ctx, FK, findOpt)
	if err != nil {
		db.Close()
		return nil
	}

	var dbKdataList []mOKX.TypeKd
	for cur.Next(db.Ctx) {
		var result mOKX.TypeKd
		cur.Decode(&result)
		fmt.Println(result.TimeStr)
		global.RunLog.Println(mJson.ToStr(result))
		dbKdataList = append(dbKdataList, result)
	}

	return _this
}

func GetTimeUnix(str string) int64 {
	t1, _ := time.ParseInLocation("2006-01-02", str, time.Local)
	return mTime.ToUnixMsec(t1)
}
