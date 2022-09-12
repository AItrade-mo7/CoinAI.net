package global

import (
	"fmt"

	"CoinAI.net/server/global/config"
	"github.com/EasyGolang/goTools/mPath"
)

func AppEnvInit() {
	// 检查配置文件在不在
	isUserEnvPath := mPath.Exists(config.File.AppEnv)
	if !isUserEnvPath {
		errStr := fmt.Errorf("没找到 app_env.json 配置文件")
		LogErr(errStr)
		panic(errStr)
	}

	config.LoadAppEnv()
}
