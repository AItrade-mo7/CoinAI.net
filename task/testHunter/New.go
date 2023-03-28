package testHunter

import (
	"CoinAI.net/server/global"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mTime"
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

func New(opt TestOpt) *TestObj {
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

	global.Run.Println("新建回测", mJson.Format(map[string]any{
		"opt":       obj,
		"StartTime": mTime.UnixFormat(obj.StartTime),
		"EndTime":   mTime.UnixFormat(obj.EndTime),
		"Days":      (obj.EndTime - obj.StartTime) / mTime.UnixTimeInt64.Day,
	}))

	return &obj
}
