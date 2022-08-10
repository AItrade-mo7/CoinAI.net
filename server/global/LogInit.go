package global

import (
	"fmt"
	"log"
	"os"

	"CoinAI.net/server/global/config"
	"github.com/EasyGolang/goTools/mLog"
	"github.com/EasyGolang/goTools/mPath"
	"github.com/EasyGolang/goTools/mTime"
)

var (
	Log    *log.Logger // 系统日志& 重大错误或者事件
	ErrLog *log.Logger // 重大错误或者事件
	WssLog *log.Logger // 系统日志& 重大错误或者事件
)

func LogInit() {
	// 检测 logs 目录
	isLogPath := mPath.Exists(config.Dir.Log)
	if !isLogPath {
		// 不存在则创建 logs 目录
		os.Mkdir(config.Dir.Log, 0o777)
	}

	// 创建一个log
	Log = mLog.NewLog(mLog.NewLogParam{
		Path: config.Dir.Log,
		Name: "Sys",
	})

	ErrLog = mLog.NewLog(mLog.NewLogParam{
		Path: config.Dir.Log,
		Name: "Err",
	})

	WssLog = mLog.NewLog(mLog.NewLogParam{
		Path: config.Dir.Log,
		Name: "Wss",
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

	ErrLog.Println(str)
}
