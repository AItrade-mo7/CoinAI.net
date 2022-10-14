package global

import (
	"CoinAI.net/server/global/config"
	"github.com/EasyGolang/goTools/mPath"
	"github.com/spf13/viper"
)

func AppEnvInit() {
	// 检查配置文件在不在
	isUserEnvPath := mPath.Exists(config.File.AppEnv)
	if !isUserEnvPath {
		CreateAppEnvFile()
	}

	viper.SetConfigFile(config.File.AppEnv)

	err := viper.ReadInConfig()
	if err != nil {
		LogErr("AppEnv 读取配置文件出错:", err)
	}
	viper.Unmarshal(&config.AppEnv)
}

func CreateAppEnvFile() {
}
