package testHunter

import (
	"fmt"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mMongo"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mTime"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TestOpt struct {
	StartTime int64
	EndTime   int64
	InstID    string
}

type TestObj struct {
	StartTime int64
	EndTime   int64
	InstID    string // BTC-USDT
	KdataList []mOKX.TypeKd
}

/*
新建基础数据

主要内容：
接受时间范围和币种
产出数据源
*/
func NewDataBase(opt TestOpt) *TestObj {
	obj := TestObj{}

	NowTime := mTime.GetUnixInt64()
	earliest := mTime.TimeParse(mTime.Lay_ss, "2020-02-01T23:00:00")

	obj.EndTime = opt.EndTime
	obj.StartTime = opt.StartTime
	obj.InstID = opt.InstID

	if obj.EndTime > NowTime {
		obj.EndTime = NowTime
	}

	if obj.StartTime < earliest {
		obj.StartTime = earliest
	}

	global.Run.Println("新建数据", mJson.Format(map[string]any{
		"InstID":    obj.InstID,
		"StartTime": mTime.UnixFormat(obj.StartTime),
		"EndTime":   mTime.UnixFormat(obj.EndTime),
		"Days":      (obj.EndTime - obj.StartTime) / mTime.UnixTimeInt64.Day,
	}))

	return &obj
}

func (_this *TestObj) StuffDBKdata() error {
	total := (_this.EndTime - _this.StartTime) / mTime.UnixTimeInt64.Hour
	if total < 1 {
		return fmt.Errorf("total 数量太少")
	}
	Timeout := int(total) * 10
	if Timeout < 100 {
		Timeout = 100
	}

	db := mMongo.New(mMongo.Opt{
		UserName: config.SysEnv.MongoUserName,
		Password: config.SysEnv.MongoPassword,
		Address:  config.SysEnv.MongoAddress,
		DBName:   "CoinMarket",
		Timeout:  Timeout,
	}).Connect().Collection(_this.InstID)
	defer db.Close()
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
		return err
	}

	AllList := []mOKX.TypeKd{}
	for cur.Next(db.Ctx) {
		var result mOKX.TypeKd
		cur.Decode(&result)
		AllList = append(AllList, result)
	}

	_this.KdataList = []mOKX.TypeKd{}
	_this.KdataList = AllList

	db.Close()

	global.Run.Println("数据填充完毕", len(_this.KdataList))

	return nil
}

func (_this *TestObj) CheckKdataList() (resErr error) {
	resErr = nil

	if len(_this.KdataList) < 1 {
		resErr = fmt.Errorf("KdataList 长度不正确")
		return
	}

	for key, val := range _this.KdataList {
		preIndex := key - 1
		if preIndex < 0 {
			preIndex = 0
		}
		preItem := _this.KdataList[preIndex]
		nowItem := _this.KdataList[key]
		if key > 0 {
			if nowItem.TimeUnix-preItem.TimeUnix != mTime.UnixTimeInt64.Hour {
				resErr = fmt.Errorf("数据检查出错, %+v", nowItem.TimeUnix-preItem.TimeUnix)
				global.LogErr("数据检查出错 backTest.CheckKdataList", val.InstID, val.TimeStr, key)
				break
			}
		}
	}

	global.Run.Println("数据检查完毕", resErr)

	return
}