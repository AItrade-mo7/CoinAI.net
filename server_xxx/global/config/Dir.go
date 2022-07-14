package config

import (
	"os"

	"github.com/EasyGolang/goTools/mStr"
)

type DirType struct {
	App string // APP 根目录
	Log string // 日志文件目录
}

var Dir DirType

type FileType struct {
	AppEnv   string // ./app_env.yaml
	Restart  string // ./restart.yaml
	Shutdown string // ./shutdown.yaml
}

var File FileType

func DirInit() {
	Dir.App, _ = os.Getwd()

	Dir.Log = mStr.Join(
		Dir.App,
		mStr.ToStr(os.PathSeparator),
		"logs",
	)

	File.AppEnv = mStr.Join(
		Dir.App,
		mStr.ToStr(os.PathSeparator),
		"app_env.yaml",
	)

	File.Restart = mStr.Join(
		Dir.App,
		mStr.ToStr(os.PathSeparator),
		"restart.sh",
	)

	File.Shutdown = mStr.Join(
		Dir.App,
		mStr.ToStr(os.PathSeparator),
		"shutdown.sh",
	)
}
