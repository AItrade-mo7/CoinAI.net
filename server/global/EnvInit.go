package global

import (
	"fmt"
	"os"

	"CoinAI.net/server/global/config"
	"CoinAI.net/server/global/dbType"
	"CoinAI.net/server/utils/dbUser"
	"CoinAI.net/server/utils/taskPush"
	"github.com/EasyGolang/goTools/mFile"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mMongo"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/EasyGolang/goTools/mStruct"
	"github.com/EasyGolang/goTools/mTime"
	"github.com/EasyGolang/goTools/mVerify"
	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AppEnvInit() {
	// 检查并读取配置文件
	byteCont, _ := os.ReadFile(config.File.AppEnv)
	jsoniter.Unmarshal(byteCont, &config.AppEnv)

	if len(config.AppEnv.UserID) < 1 {
		err := fmt.Errorf("启动错误，缺少 AppEnv.UserID 字段: %v", mJson.ToStr(config.AppEnv))
		LogErr(err)
		panic(err)
	}

	if len(config.AppEnv.Port) < 1 {
		err := fmt.Errorf("启动错误，缺少 AppEnv.Port 字段: %v", mJson.ToStr(config.AppEnv))
		LogErr(err)
		panic(err)
	}
	// 回填 用户 信息
	GetMainUser()

	// 回填 IP
	config.AppEnv.IP = GetLocalAPI()

	if !mVerify.IsIP(config.AppEnv.IP) {
		err := fmt.Errorf("启动错误， 系统 IP 获取失败 %+v", config.AppEnv.IP)
		LogErr(err)
		panic(err)
	}
	config.AppEnv.ServeID = mStr.Join(config.AppEnv.IP, ":", config.AppEnv.Port)

	ReadeDBAppEnv()

	if len(config.AppEnv.SysName) < 1 {
		config.AppEnv.SysName = mStr.Join(config.MainUser.NickName, "的 CoinAI")
	}

	if config.AppEnv.CreateTime < mTime.TimeParse(mTime.Lay_Y, "2022") { // 表示没有创建时间
		config.AppEnv.CreateTime = mTime.GetUnixInt64()
	}

	// 设置  默认 最大 ApiKey 数量
	if config.AppEnv.MaxApiKeyNum == 0 {
		config.AppEnv.MaxApiKeyNum = 32
	}
	config.AppEnv.SysVersion = config.AppInfo.Version
	config.AppEnv.Type = config.SysName

	WriteAppEnv()
}

func ReadeDBAppEnv() {
	db := mMongo.New(mMongo.Opt{
		UserName: config.SysEnv.MongoUserName,
		Password: config.SysEnv.MongoPassword,
		Address:  config.SysEnv.MongoAddress,
		DBName:   "AIServe",
	}).Connect().Collection("CoinAI")
	defer db.Close()

	findOpt := options.FindOne()
	findOpt.SetSort(map[string]int{
		"TimeUnix": 1,
	})
	FK := bson.D{{
		Key:   "ServeID",
		Value: config.AppEnv.ServeID,
	}}

	var AppEnv dbType.AppEnvType
	db.Table.FindOne(db.Ctx, FK, findOpt).Decode(&AppEnv)

	if len(AppEnv.ServeID) > 4 && len(AppEnv.UserID) > 4 {
		if config.AppEnv.UserID == AppEnv.UserID {
			config.AppEnv = AppEnv
		} else {
			err := fmt.Errorf("启动错误，当前服务归属用户不正确")
			LogErr(err)
			panic(err)
		}
	}
}

// 写入 config.AppEnv
func WriteAppEnv() {
	mFile.Write(config.File.AppEnv, mJson.JsonFormat(mJson.ToJson(config.AppEnv)))

	db := mMongo.New(mMongo.Opt{
		UserName: config.SysEnv.MongoUserName,
		Password: config.SysEnv.MongoPassword,
		Address:  config.SysEnv.MongoAddress,
		DBName:   "AIServe",
	}).Connect().Collection("CoinAI")
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
	UK = append(UK, bson.E{
		Key: "$set",
		Value: bson.D{
			{
				Key:   "UpdateTime",
				Value: mTime.GetUnixInt64(),
			},
		},
	})
	upOpt := options.Update()
	upOpt.SetUpsert(true)
	_, err := db.Table.UpdateOne(db.Ctx, FK, UK, upOpt)
	if err != nil {
		LogErr("config.AppEnv 数据更插失败", err)
	}
	Log.Println("config.AppEnv 已更新至数据库", mJson.Format(config.AppEnv))
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

	var resData struct {
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

	jsoniter.Unmarshal(res, &resData)
	if resData.Code < 0 {
		LogErr(resData.Msg)
		return ""
	}
	ip = resData.Data.IP
	return
}

func GetMainUser() {
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
	// 清空通知邮箱
	config.NoticeEmail = []string{}
	// 清空 MainUser
	config.MainUser = dbType.UserTable{}

	// 回填通知邮箱
	NoticeEmail := []string{}
	NoticeEmail = append(NoticeEmail, UserDB.Data.Email)
	config.NoticeEmail = NoticeEmail

	// 回填用户信息
	config.MainUser = UserDB.Data
}
