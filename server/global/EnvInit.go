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

	Log.Println(`第一次读取文件`, mJson.Format(config.AppEnv))

	if len(config.AppEnv.Port) < 1 {
		config.AppEnv.Port = "9000"
	}
	if len(config.AppEnv.IP) < 1 {
		config.AppEnv.IP = reqDataCenter.GetLocalIP()
	}

	config.AppEnv.Version = config.AppInfo.Version
	config.AppEnv.Name = config.AppInfo.Name

	WriteAppEnv()

	Log.Println(`第二次读取文件`, mJson.Format(config.AppEnv))
}

// 写入 config.AppEnv
func WriteAppEnv() {
	mFile.Write(config.File.AppEnv, mJson.JsonFormat(mJson.ToJson(config.AppEnv)))
}
