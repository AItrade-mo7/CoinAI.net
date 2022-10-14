package global

import (
	"bytes"
	"text/template"

	"CoinAI.net/server/global/config"
	"CoinAI.net/server/tmpl"
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

	CreateReboot()
	CreateShutdown()
}

func WriteAppEnv() {
	// 如果不存在 app_env.json 则创建写入
	mFile.Write(config.File.AppEnv, mJson.ToStr(config.AppEnv))
}

func CreateReboot() {
	Body := new(bytes.Buffer)
	Tmpl := template.Must(template.New("").Parse(tmpl.Reboot))
	Tmpl.Execute(Body, tmpl.RebootParam{
		Port: config.AppEnv.Port,
		Path: config.Dir.App,
	})
	Cont := Body.String()

	mFile.Write(config.File.Reboot, Cont)
}

func CreateShutdown() {
	Body := new(bytes.Buffer)
	Tmpl := template.Must(template.New("").Parse(tmpl.Shutdown))
	Tmpl.Execute(Body, tmpl.ShutdownParam{
		Port: config.AppEnv.Port,
		Path: config.Dir.App,
	})
	Cont := Body.String()

	mFile.Write(config.File.Shutdown, Cont)
}
