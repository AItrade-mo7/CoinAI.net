package backTest

import (
	"fmt"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/global/dbType"
	"CoinAI.net/server/hunter"
	"CoinAI.net/server/okxInfo"
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
	KdataList []mOKX.TypeKd
}

func NewTest(opt TestOpt) *TestObj {
	obj := TestObj{}
	now := mTime.GetUnixInt64()

	if opt.EndTime < dbType.MinTime { // 如果太小，则变成当前
		opt.EndTime = now
	}

	if opt.StartTime < dbType.MinTime { // 如果太小，则自动变成过去一个月
		opt.StartTime = mTime.UnixTimeInt64.Day * 30
	}

	if opt.StartTime < dbType.DBKdataStart {
		opt.StartTime = dbType.DBKdataStart
	}

	if opt.StartTime > opt.EndTime {
		global.LogErr("backTest.NewTest 开始时间不可以大于结束时间")
	}

	obj.StartTime = opt.StartTime
	obj.EndTime = opt.EndTime

	obj.TableName = opt.CcyName + "USDT"
	obj.InstID = opt.CcyName + "-USDT"

	return &obj
}

func (_this *TestObj) GetDBKdata() *TestObj {
	total := (_this.EndTime - _this.StartTime) / mTime.UnixTimeInt64.Hour

	if total < 1 {
		return nil
	}

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

	_this.KdataList = []mOKX.TypeKd{}

	for cur.Next(db.Ctx) {
		var result mOKX.TypeKd
		cur.Decode(&result)
		_this.KdataList = append(_this.KdataList, result)
	}

	return _this
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

	return
}

var (
	EMA_Arr = []string{}
	MA_Arr  = []string{}
)

func (_this *TestObj) MockData() {
	// 在这里开始执行模拟数据流的流动
	// 执行一次清理
	okxInfo.TradeKdataList = []okxInfo.TradeKdType{}
	EMA_Arr = []string{}
	MA_Arr = []string{}
	FormatEnd := []mOKX.TypeKd{}

	for _, Kdata := range _this.KdataList {
		// 回填历史数据
		okxInfo.NowKdataList = append(okxInfo.NowKdataList, Kdata)
		if len(okxInfo.NowKdataList) >= 100 {
			// 开始执行分析
			FormatEnd = append(FormatEnd, Kdata)
			TradeKdata := hunter.NewTradeKdata(Kdata, FormatEnd)
			okxInfo.TradeKdataList = append(okxInfo.TradeKdataList, TradeKdata)

			hunter.Analy()
		}
	}
}
