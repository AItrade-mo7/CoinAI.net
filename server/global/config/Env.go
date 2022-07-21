package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var SysEnv struct {
	MongoAddress  string
	MongoPassword string
	MongoUserName string
}

func LoadSysEnv() {
	SysEnv.MongoAddress = "trade.mo7.cc:17017"
	SysEnv.MongoPassword = "asdasd55555"
	SysEnv.MongoUserName = "mo7"
}

var AppEnv struct {
	Port        string `json:"Port"`
	UserID      string `json:"UserID"`
	CoinServeID string `json:"CoinServeID"`
	RunMod      int    // 0 则为正常模式 ， 1 则为数据模拟模式
}

var AppInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func LoadAppEnv() {
	viper.SetConfigFile(File.AppEnv)

	err := viper.ReadInConfig()
	if err != nil {
		errStr := fmt.Errorf("AppEnv 读取配置文件出错: %+v", err)
		LogErr(errStr)
		panic(errStr)
	}
	viper.Unmarshal(&AppEnv)
}
