package global

import (
	"fmt"
	"os"

	"CoinAI.net/server/global/config"
	"CoinAI.net/server/utils/dbUser"
	"CoinAI.net/server/utils/taskPush"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/EasyGolang/goTools/mVerify"
	jsoniter "github.com/json-iterator/go"
)

func AppEnvInit() {
	// 检查并读取配置文件
	byteCont, _ := os.ReadFile(config.File.AppEnv)
	jsoniter.Unmarshal(byteCont, &config.AppEnv)

	if len(config.AppEnv.Port) < 1 {
		err := fmt.Errorf("启动错误，缺少 AppEnv.Port 字段: %v", mJson.ToStr(config.AppEnv))
		LogErr(err)
		panic(err)
	}

	if len(config.AppEnv.UserID) < 1 {
		err := fmt.Errorf("启动错误，缺少 AppEnv.UserID 字段: %v", mJson.ToStr(config.AppEnv))
		LogErr(err)
		panic(err)
	}

	// 在这里 获取 用户 信息
	UserDB, err := dbUser.NewUserDB(dbUser.NewUserOpt{
		UserID: config.AppEnv.UserID,
	})
	if err != nil {
		UserDB.DB.Close()
		err := fmt.Errorf("启动错误，数据库链接失败: %v", err)
		LogErr(err)
		panic(err)
	}
	defer UserDB.DB.Close()

	if len(UserDB.UserID) < 1 {
		err := fmt.Errorf("启动错误，用户未找到: %v", UserDB.UserID)
		LogErr(err)
		panic(err)
	}

	config.AppEnv.IP = GetLocalAPI()

	if !mVerify.IsIP(config.AppEnv.IP) {
		LogErr("启动错误，缺少 AppEnv.UserID 字段")
		panic("ip 获取失败")
	}

	config.AppEnv.ServeID = mStr.Join(config.AppEnv.IP, ":", config.AppEnv.Port)

	// ReadeDBAppEnv(config.AppEnv.ServeID)

	// if len(config.AppEnv.Name) < 1 {
	// 	config.AppEnv.Name = "我的 CoinAI"
	// }
	// config.AppEnv.Version = config.AppInfo.Version

	// 设置默认杠杆倍数 TradeLever
	// if config.AppEnv.TradeLever == 0 {
	// 	config.AppEnv.TradeLever = 5
	// }
	// 设置  默认 最大 ApiKey 数量
	// if config.AppEnv.MaxApiKeyNum == 0 {
	// 	config.AppEnv.MaxApiKeyNum = 32
	// }

	// WriteAppEnv()
}

// func ReadeDBAppEnv(ServeID string) {
// 	db := mMongo.New(mMongo.Opt{
// 		UserName: config.SysEnv.MongoUserName,
// 		Password: config.SysEnv.MongoPassword,
// 		Address:  config.SysEnv.MongoAddress,
// 		DBName:   "AItrade",
// 	}).Connect().Collection("CoinAINet")
// 	defer db.Close()

// 	findOpt := options.FindOne()
// 	findOpt.SetSort(map[string]int{
// 		"TimeUnix": 1,
// 	})
// 	FK := bson.D{{
// 		Key:   "ServeID",
// 		Value: ServeID,
// 	}}

// 	var AppEnv config.AppEnvType
// 	db.Table.FindOne(db.Ctx, FK, findOpt).Decode(&AppEnv)

// 	if len(AppEnv.ServeID) > 4 && len(AppEnv.UserID) > 4 {
// 		config.AppEnv = AppEnv
// 	}
// }

// // 写入 config.AppEnv
// func WriteAppEnv() {
// 	mFile.Write(config.File.AppEnv, mJson.JsonFormat(mJson.ToJson(config.AppEnv)))

// 	db := mMongo.New(mMongo.Opt{
// 		UserName: config.SysEnv.MongoUserName,
// 		Password: config.SysEnv.MongoPassword,
// 		Address:  config.SysEnv.MongoAddress,
// 		DBName:   "AItrade",
// 	}).Connect().Collection("CoinAINet")
// 	defer db.Close()

// 	findOpt := options.FindOne()
// 	findOpt.SetSort(map[string]int{
// 		"TimeUnix": -1,
// 	})

// 	FK := bson.D{{
// 		Key:   "ServeID",
// 		Value: config.AppEnv.ServeID,
// 	}}
// 	UK := bson.D{}
// 	mStruct.Traverse(config.AppEnv, func(key string, val any) {
// 		UK = append(UK, bson.E{
// 			Key: "$set",
// 			Value: bson.D{
// 				{
// 					Key:   key,
// 					Value: val,
// 				},
// 			},
// 		})
// 	})
// 	upOpt := options.Update()
// 	upOpt.SetUpsert(true)
// 	_, err := db.Table.UpdateOne(db.Ctx, FK, UK, upOpt)
// 	if err != nil {
// 		LogErr("config.AppEnv 数据更插失败", err)
// 	}
// 	Log.Println("config.AppEnv 已更新至数据库", mJson.Format(config.AppEnv))
// }

// 方法
// func GetOKXKey(Index int) mOKX.TypeOkxKey {
// 	// ReturnKey := mOKX.TypeOkxKey{}
// 	// for key, val := range AppEnv.ApiKeyList {
// 	// 	if key == Index {
// 	// 		ReturnKey = val
// 	// 		break
// 	// 	}
// 	// }
// 	return ReturnKey
// }

type PublicPingType struct {
	Code int64 `json:"Code"`
	Data struct {
		APIInfo struct {
			Name    string `json:"Name"`
			Version string `json:"Version"`
		} `json:"ApiInfo"`
		IP        string         `json:"IP"`
		Method    string         `json:"Method"`
		Path      string         `json:"Path"`
		ResParam  map[string]any `json:"ResParam"`
		UserAgent string         `json:"UserAgent"`
	} `json:"Data"`
	Msg string `json:"Msg"`
}

func GetLocalAPI() (ip string) {
	res, err := taskPush.Request(taskPush.RequestOpt{
		Origin: config.SysEnv.MessageBaseUrl,
		Path:   "/ping",
	})
	if err != nil {
		LogErr(err)
		return ""
	}
	var resData PublicPingType
	jsoniter.Unmarshal(res, &resData)
	if resData.Code < 0 {
		LogErr(resData.Msg)
		return ""
	}
	ip = resData.Data.IP
	return
}
