package verifyCode

import (
	"fmt"

	"Hunter.net/server/global/config"
	"Hunter.net/server/global/dbType"
	"github.com/EasyGolang/goTools/mCount"
	"github.com/EasyGolang/goTools/mEncrypt"
	"github.com/EasyGolang/goTools/mMongo"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/EasyGolang/goTools/mTime"
	"go.mongodb.org/mongo-driver/bson"
)

type CheckEmailCodeParam struct {
	Email string
	Code  string
}

func CheckEmailCode(opt CheckEmailCodeParam) (resErr error) {
	resErr = nil

	db := mMongo.New(mMongo.Opt{
		UserName: config.SysEnv.MongoUserName,
		Password: config.SysEnv.MongoPassword,
		Address:  config.SysEnv.MongoAddress,
		DBName:   "Hunter",
	}).Connect().Collection("EmailCode")

	err := db.Ping()
	if err != nil {
		db.Close()
		resErr = fmt.Errorf("EmailCode,数据库连接错误 %+v", err)
		return
	}

	var result dbType.EmailCodeTable

	FK := bson.D{{
		Key:   "Email",
		Value: opt.Email,
	}}
	db.Table.FindOne(db.Ctx, FK).Decode(&result)

	DBCode := mEncrypt.MD5(result.Code)

	if DBCode != opt.Code {
		db.Close()
		resErr = fmt.Errorf("验证码不正确")
		return
	}

	// 校验时间
	sendTime := mStr.ToStr(result.SendTime)
	nowTime := mTime.GetUnix()
	subStr := mCount.Sub(nowTime, sendTime)

	// 20 分钟
	if mCount.Le(subStr, mCount.Mul(mTime.UnixTime.Minute, "20")) > 0 {
		db.Close()
		resErr = fmt.Errorf("验证码已过期")
		return
	}
	db.Close()

	return
}
