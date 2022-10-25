package main

import (
	_ "embed"
	"fmt"
	"time"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/okxApi"
	"CoinAI.net/server/ready"
	jsoniter "github.com/json-iterator/go"
)

//go:embed package.json
var AppPackage []byte

func main() {
	jsoniter.Unmarshal(AppPackage, &config.AppInfo)
	// 初始化系统参数
	global.Start()

	// 数据准备
	ready.Start()

	OkxKey := config.GetOKXKey(0)

	for i := 0; i < 20; i++ {
		time.Sleep(time.Second / 3)
		fmt.Println("======开始=======")
		OKXAccount, err := okxApi.NewAccount(okxApi.AccountParam{
			OkxKey: OkxKey,
		})
		if err != nil {
			fmt.Println("1111", err)
		}
		err = OKXAccount.SetPositionMode()
		if err != nil {
			fmt.Println("2222", err)
		}
		err = OKXAccount.SetLeverage()
		if err != nil {
			fmt.Println("3333", err)
		}
		err = OKXAccount.GetMaxSize()
		if err != nil {
			fmt.Println("4444", err)
		}
		err = OKXAccount.GetBalance()
		if err != nil {
			fmt.Println("5555", err)
		}
	}

	// 启动路由
	// router.Start()
}
