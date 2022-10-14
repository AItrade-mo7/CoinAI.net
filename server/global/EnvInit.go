package global

import (
	"CoinAI.net/server/global/config"
	"github.com/EasyGolang/goTools/mFile"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mPath"
	"github.com/spf13/viper"
)

func AppEnvInit() {
	// 检查配置文件在不在
	isUserEnvPath := mPath.Exists(config.File.AppEnv)
	if !isUserEnvPath {
		WriteAppEnv()
	}

	viper.SetConfigFile(config.File.AppEnv)

	err := viper.ReadInConfig()
	if err != nil {
		LogErr("AppEnv 读取配置文件出错:", err)
	}
	viper.Unmarshal(&config.AppEnv)

	if len(config.AppEnv.Port) < 1 {
		config.AppEnv.Port = "9453"
	}
	WriteAppEnv()
}

func WriteAppEnv() {
	// 如果不存在 app_env.json 则创建写入
	mFile.Write(config.File.AppEnv, mJson.ToStr(config.AppEnv))
}
