package config

import (
	"fmt"
	"os"

	"github.com/EasyGolang/goTools/mPath"
	"github.com/EasyGolang/goTools/mStr"
)

type DirType struct {
	Home     string // Home 根目录
	App      string // APP 根目录
	Log      string // 日志文件目录
	JsonData string // json 数据存放目录
}

var Dir DirType

type FileType struct {
	AppEnv   string // ./app_env.json
	ReBoot   string // ./ReBoot.sh
	Shutdown string // ./Shutdown.sh
}

var File FileType

func DirInit() {
	Dir.Home = mPath.HomePath()

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

	File.ReBoot = mStr.Join(
		Dir.App,
		mStr.ToStr(os.PathSeparator),
		"ReBoot.sh",
	)

	File.Shutdown = mStr.Join(
		Dir.App,
		mStr.ToStr(os.PathSeparator),
		"Shutdown.sh",
	)

	// 检测 logs 目录
	if !mPath.Exists(Dir.Log) {
		os.MkdirAll(Dir.Log, 0o777)
	}
	if !mPath.Exists(Dir.JsonData) {
		os.MkdirAll(Dir.JsonData, 0o777)
	}

	if !mPath.Exists(File.Shutdown) {
		err := fmt.Errorf("缺少文件 %+v", File.Shutdown)
		panic(err)
	}

	if !mPath.Exists(File.ReBoot) {
		err := fmt.Errorf("缺少文件 %+v", File.ReBoot)
		panic(err)
	}
}
