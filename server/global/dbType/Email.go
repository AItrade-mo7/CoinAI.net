package dbType

import (
	"time"

	"github.com/EasyGolang/goTools/mTime"
)

// 发送验证码的表结构 ========= EmailCode ============
type EmailCodeTable struct {
	Email    string `bson:"Email"`
	Code     string `bson:"Code"`
	SendTime int64  `bson:"SendTime"`
}

var DBKdataStart = int64(1578171600000) // 2020-01-05T05:00:00

var MinTime = mTime.UnixTimeInt64.Day * 2190 // 6年前

// 2006-1-2
func ParseTime(str string) int64 {
	t1, _ := time.ParseInLocation("2006-1-2", str, time.Local)
	return mTime.ToUnixMsec(t1)
}
