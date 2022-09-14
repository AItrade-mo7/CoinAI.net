package global

import (
	"fmt"

	"CoinAI.net/server/global/config"
	"github.com/EasyGolang/goTools/mPath"
	"github.com/spf13/viper"
)

func AppEnvInit() {
	// 检查配置文件在不在
	isUserEnvPath := mPath.Exists(config.File.AppEnv)
	if !isUserEnvPath {
		errStr := fmt.Errorf("没找到 app_env.json 配置文件")
		LogErr(errStr)
		panic(errStr)
	}

	viper.SetConfigFile(config.File.AppEnv)

	err := viper.ReadInConfig()
	if err != nil {
		errStr := fmt.Errorf("AppEnv 读取配置文件出错: %+v", err)
		LogErr(errStr)
		panic(errStr)
	}
	viper.Unmarshal(&config.AppEnv)
}
