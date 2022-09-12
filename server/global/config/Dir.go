package config

import (
	"os"

	"github.com/EasyGolang/goTools/mPath"
	"github.com/EasyGolang/goTools/mStr"
)

type DirType struct {
	App      string // APP 根目录
	Log      string // 日志文件目录
	JsonData string // json 数据存放目录
}

var Dir DirType

type FileType struct {
	AppEnv   string // ./app_env.json
	Restart  string // ./restart.sh
	Shutdown string // ./shutdown.sh
}

var File FileType

func DirInit() {
	Dir.App, _ = os.Getwd()

	Dir.Log = mStr.Join(
		Dir.App,
		mStr.ToStr(os.PathSeparator),
		"logs",
	)

	Dir.JsonData = mStr.Join(
		Dir.App,
		mStr.ToStr(os.PathSeparator),
		"jsonData",
	)

	File.AppEnv = mStr.Join(
		Dir.App,
		mStr.ToStr(os.PathSeparator),
		"app_env.json",
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

	// 检测 logs 目录
	if !mPath.Exists(Dir.Log) {
		os.MkdirAll(Dir.Log, 0o777)
	}
	if !mPath.Exists(Dir.JsonData) {
		os.MkdirAll(Dir.JsonData, 0o777)
	}
}
