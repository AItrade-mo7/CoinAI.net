package global

import (
	"os"

	"CoinAI.net/server/global/config"
	"CoinAI.net/server/utils/reqDataCenter"
	"github.com/EasyGolang/goTools/mFile"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mPath"
	jsoniter "github.com/json-iterator/go"
)

func AppEnvInit() {
	// 检查并读取配置文件
	isEnvPath := mPath.Exists(config.File.AppEnv)
	if isEnvPath {
		byteCont, _ := os.ReadFile(config.File.AppEnv)
		jsoniter.Unmarshal(byteCont, &config.AppEnv)
	}

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
