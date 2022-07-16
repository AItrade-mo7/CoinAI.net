package global

import (
	"fmt"
	"time"

	"CoinAI.net/server/global/config"
	"github.com/EasyGolang/goTools/mCycle"
	"github.com/EasyGolang/goTools/mPath"
)

func Start() {
	// 初始化目录列表
	config.DirInit()

	// 初始化日志系统 保证日志可用
	mCycle.New(mCycle.Opt{
		Func:      LogInit,
		SleepTime: time.Hour * 8,
	}).Start()

	// 加载App启动配置文件
	AppEnvInit()

	// 检测文件
	isRestartShell := mPath.Exists(config.File.Restart)
	if !isRestartShell {
		errStr := fmt.Errorf("缺少文件:" + config.File.Restart)
		LogErr(errStr)
		panic(errStr)
	}
	isShutdownShell := mPath.Exists(config.File.Shutdown)
	if !isShutdownShell {
		errStr := fmt.Errorf("缺少文件:" + config.File.Shutdown)
		LogErr(errStr)
		panic(errStr)
	}

	Log.Println(`系统初始化完成`)
}
