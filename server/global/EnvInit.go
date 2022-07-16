package global

import (
	"fmt"

	"CoinAI.net/server/global/config"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mPath"
)

func AppEnvInit() {
	config.LoadSysEnv()
	Log.Println("加载 SysEnv : ", mJson.JsonFormat(mJson.ToJson(config.SysEnv)))

	// 检查配置文件在不在
	isUserEnvPath := mPath.Exists(config.File.AppEnv)
	if !isUserEnvPath {
		errStr := fmt.Errorf("没找到 app_env.yaml 配置文件")
		LogErr(errStr)
		panic(errStr)
	}

	config.LoadAppEnv()

	Log.Println("加载 AppEnv : ", mJson.JsonFormat(mJson.ToJson(config.AppEnv)))
}
