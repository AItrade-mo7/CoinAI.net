package global

import (
	"os"

	"CoinAI.net/server/global/config"
	"CoinAI.net/server/utils/reqDataCenter"
	"github.com/EasyGolang/goTools/mFile"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mMongo"
	"github.com/EasyGolang/goTools/mPath"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/EasyGolang/goTools/mStruct"
	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AppEnvInit() {
	// 检查并读取配置文件
	isEnvPath := mPath.Exists(config.File.AppEnv)
	if isEnvPath {
		byteCont, _ := os.ReadFile(config.File.AppEnv)
		jsoniter.Unmarshal(byteCont, &config.AppEnv)
	}

	if len(config.AppEnv.Port) < 1 {
		config.AppEnv.Port = "9000"
	}
	config.AppEnv.IP = reqDataCenter.GetLocalIP()
	config.AppEnv.ServeID = mStr.Join(config.AppEnv.IP, ":", config.AppEnv.Port)

	ReadeDBAppEnv(config.AppEnv.ServeID)

	// if len(config.AppEnv.Name) < 1 {
	// 	config.AppEnv.Name = "我的 CoinAI"
	// }
	// config.AppEnv.Version = config.AppInfo.Version

	// 设置默认杠杆倍数 TradeLever
	// if config.AppEnv.TradeLever == 0 {
	// 	config.AppEnv.TradeLever = 5
	// }
	// 设置  默认 最大 ApiKey 数量
	if config.AppEnv.MaxApiKeyNum == 0 {
		config.AppEnv.MaxApiKeyNum = 32
	}

	WriteAppEnv()
}

func ReadeDBAppEnv(ServeID string) {
	db := mMongo.New(mMongo.Opt{
		UserName: config.SysEnv.MongoUserName,
		Password: config.SysEnv.MongoPassword,
		Address:  config.SysEnv.MongoAddress,
		DBName:   "AItrade",
	}).Connect().Collection("CoinAINet")
	defer db.Close()

	findOpt := options.FindOne()
	findOpt.SetSort(map[string]int{
		"TimeUnix": 1,
	})
	FK := bson.D{{
		Key:   "ServeID",
		Value: ServeID,
	}}

	var AppEnv config.AppEnvType
	db.Table.FindOne(db.Ctx, FK, findOpt).Decode(&AppEnv)

	if len(AppEnv.ServeID) > 4 && len(AppEnv.UserID) > 4 {
		config.AppEnv = AppEnv
	}
}

// 写入 config.AppEnv
func WriteAppEnv() {
	mFile.Write(config.File.AppEnv, mJson.JsonFormat(mJson.ToJson(config.AppEnv)))

	db := mMongo.New(mMongo.Opt{
		UserName: config.SysEnv.MongoUserName,
		Password: config.SysEnv.MongoPassword,
		Address:  config.SysEnv.MongoAddress,
		DBName:   "AItrade",
	}).Connect().Collection("CoinAINet")
	defer db.Close()

	findOpt := options.FindOne()
	findOpt.SetSort(map[string]int{
		"TimeUnix": -1,
	})

	FK := bson.D{{
		Key:   "ServeID",
		Value: config.AppEnv.ServeID,
	}}
	UK := bson.D{}
	mStruct.Traverse(config.AppEnv, func(key string, val any) {
		UK = append(UK, bson.E{
			Key: "$set",
			Value: bson.D{
				{
					Key:   key,
					Value: val,
				},
			},
		})
	})
	upOpt := options.Update()
	upOpt.SetUpsert(true)
	_, err := db.Table.UpdateOne(db.Ctx, FK, UK, upOpt)
	if err != nil {
		LogErr("config.AppEnv 数据更插失败", err)
	}
	Log.Println("config.AppEnv 已更新至数据库", mJson.Format(config.AppEnv))
}
