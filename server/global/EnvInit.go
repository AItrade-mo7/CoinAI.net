package global

import (
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/utils/reqDataCenter"
	"github.com/EasyGolang/goTools/mFile"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/spf13/viper"
)

func AppEnvInit() {
	// 检查配置文件在不在

	viper.SetConfigFile(config.File.AppEnv)

	viper.Unmarshal(&config.AppEnv)

	if len(config.AppEnv.Port) < 1 {
		config.AppEnv.Port = "9453"
	}
	if len(config.AppEnv.IP) < 1 {
		config.AppEnv.IP = reqDataCenter.GetLocalIP()
	}

	WriteAppEnv()
}

func WriteAppEnv() {
	// 如果不存在 app_env.json 则创建写入

	mFile.Write(config.File.AppEnv, mJson.JsonFormat(mJson.ToJson(config.AppEnv)))
}
