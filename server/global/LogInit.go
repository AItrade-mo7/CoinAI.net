package global

import (
	"fmt"
	"log"

	"CoinAI.net/server/global/config"
	"CoinAI.net/server/tmpl"
	"github.com/EasyGolang/goTools/mLog"
	"github.com/EasyGolang/goTools/mTime"
)

var (
	Log    *log.Logger // 系统日志
	WssLog *log.Logger // Wss 数据
	RunLog *log.Logger // 运行过程
)

func LogInit() {
	// 创建一个log
	Log = mLog.NewLog(mLog.NewLogParam{
		Path: config.Dir.Log,
		Name: "Sys",
	})

	WssLog = mLog.NewLog(mLog.NewLogParam{
		Path: config.Dir.Log,
		Name: "Wss",
	})

	RunLog = mLog.NewLog(mLog.NewLogParam{
		Path: config.Dir.Log,
		Name: "Run",
	})

	// 设定清除log
	mLog.Clear(mLog.ClearParam{
		Path:      config.Dir.Log,
		ClearTime: mTime.UnixTimeInt64.Day * 10,
	})

	// 将方法重载到 config 内部,便于执行
	config.LogErr = LogErr
}

func LogErr(sum ...any) {
	str := fmt.Sprintf("系统错误 : %+v", sum)
	Email := Email(EmailOpt{
		To:       config.Email.To,
		Subject:  "系统出错",
		Template: tmpl.SysEmail,
		SendData: tmpl.SysParam{
			Message: str,
			SysTime: mTime.IsoTime(),
		},
	})
	Log.Println(str)
	go Email.Send()
}
