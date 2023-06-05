package config

import (
	"os"

	"github.com/EasyGolang/goTools/mPath"
	"github.com/EasyGolang/goTools/mStr"
)

var Dir struct {
	Home     string // Home 根目录
	App      string // APP 根目录
	Log      string // 日志文件目录
	JsonData string // json 数据存放目录
}

type FileType struct {
	AppEnv   string // ./app_env.json
	Reboot   string // ./Reboot.sh
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

	File.Reboot = mStr.Join(
		Dir.App,
		mStr.ToStr(os.PathSeparator),
		"Reboot.sh",
	)

	File.Shutdown = mStr.Join(
		Dir.App,
		mStr.ToStr(os.PathSeparator),
		"Shutdown.sh",
	)

	// 检测 JsonData 目录
	isJsonDataPath := mPath.Exists(Dir.JsonData)
	if !isJsonDataPath {
		// 不存在则创建 JsonData 目录
		os.MkdirAll(Dir.JsonData, 0o777)
	}

	// 检测 logs 目录
	isLogPath := mPath.Exists(Dir.Log)
	if !isLogPath {
		// 不存在则创建 logs 目录
		os.MkdirAll(Dir.Log, 0o777)
	}

	//if !mPath.Exists(File.AppEnv) {
	//	err := fmt.Errorf("缺少文件 %+v", File.AppEnv)
	//	panic(err)
	//}

	//if !mPath.Exists(File.Shutdown) {
	//	err := fmt.Errorf("缺少文件 %+v", File.Shutdown)
	//	panic(err)
	//}
	//
	//if !mPath.Exists(File.Reboot) {
	//	err := fmt.Errorf("缺少文件 %+v", File.Reboot)
	//	panic(err)
	//}
}
