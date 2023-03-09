package global

import (
	"fmt"
	"log"

	"CoinAI.net/server/global/config"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mLog"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/EasyGolang/goTools/mTime"
)

var (
	Log      *log.Logger // 系统日志
	WssLog   *log.Logger // Wss 数据
	KdataLog *log.Logger // Kdata 日志
	RunLog   *log.Logger // 运行过程
	TradeLog *log.Logger // 交易API
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
	KdataLog = mLog.NewLog(mLog.NewLogParam{
		Path: config.Dir.Log,
		Name: "Kdata",
	})

	RunLog = mLog.NewLog(mLog.NewLogParam{
		Path: config.Dir.Log,
		Name: "Run",
	})

	TradeLog = mLog.NewLog(mLog.NewLogParam{
		Path: config.Dir.Log,
		Name: "Trade",
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
	str := fmt.Sprintf("系统错误: %+v", sum)
	Log.Println(str)

	message := ""
	if len(sum) > 0 {
		message = mStr.ToStr(sum[0])
	}
	content := mJson.Format(sum)

	err := taskPush.SysEmail(taskPush.SysEmailOpt{
		From:        config.SysName,
		Subject:     "系统错误",
		Title:       config.SysName + " 系统出错",
		Message:     message,
		Content:     content,
		Description: "出现系统错误",
	})
	Log.Println("邮件已发送", err)
}
