package global

import (
	"fmt"
	"github.com/EasyGolang/goTools/mPath"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"

	"CoinAI.net/server/global/config"
	"CoinAI.net/server/utils/taskPush"
	"github.com/EasyGolang/goTools/mLog"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/EasyGolang/goTools/mTime"
)

var (
	Log      *zap.Logger // 系统日志
	Run      *zap.Logger // 运行日志
	WssLog   *zap.Logger // Wss 数据
	KdataLog *zap.Logger // Kdata 日志
	TradeLog *zap.Logger // 交易模块
	OKXLogo  *zap.Logger // 交易模块
)

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func LogInit() {
	// 创建一个log
	encoder := getEncoder()

	loggerCreater := func(path, name string) *zap.Logger {
		// 检测 logs 目录
		isLogPath := mPath.Exists(path)
		if !isLogPath {
			// 不存在则创建 logs 目录
			_ = os.Mkdir(path, 0o777)
		}
		file := path + "/" + name + "-T" + time.Now().Format("06年1月02日15时") + ".log"
		logF, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o777)
		if nil != err {
			return nil
		}
		return zap.New(zapcore.NewTee(zapcore.NewCore(encoder, zapcore.AddSync(logF), zapcore.ErrorLevel)), zap.AddCaller())
	}
	Log = loggerCreater(config.Dir.Log, "Sys")
	WssLog = loggerCreater(config.Dir.Log, "Wss")
	KdataLog = loggerCreater(config.Dir.Log, "Kdata")
	Run = loggerCreater(config.Dir.Log, "Run")
	TradeLog = loggerCreater(config.Dir.Log, "Trade")
	OKXLogo = loggerCreater(config.Dir.Log, "OKX")

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
	Log.Error(str)

	message := ""
	if len(sum) > 0 {
		message = mStr.ToStr(sum[0])
	}
	content := mStr.ToStr(sum)

	TitleName := config.SysName

	if len(config.AppEnv.SysName) > 0 {
		TitleName = config.AppEnv.SysName
	}

	err := taskPush.SysEmail(taskPush.SysEmailOpt{
		From:        config.SysName,
		To:          config.NoticeEmail,
		Subject:     "系统错误",
		Title:       TitleName + " 系统出错",
		Message:     message,
		Content:     content,
		Description: "出现系统错误",
	})
	Log.Error("邮件已发送: ", zap.String("err", err.Error()))
}
