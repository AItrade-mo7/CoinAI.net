package main

import (
	_ "embed"
	"fmt"

	"CoinAI.net/server/global"
	"CoinAI.net/task/testHunter"
	"github.com/EasyGolang/goTools/mTime"
)

func main() {
	// 初始化系统参数
	global.Start()

	// 新建回测
	back := testHunter.New(testHunter.TestOpt{
		StartTime: mTime.TimeParse(mTime.Lay_ss, "2022-12-12T00:00:00"),
		EndTime:   mTime.TimeParse(mTime.Lay_ss, "2023-01-01T00:00:00"),
		InstID:    "BTC-USDT",
	})
	err := back.StuffDBKdata()
	if err != nil {
		fmt.Println("出错", err)
	}
	err = back.CheckKdataList() // 检查数据是否出错
	if err != nil {
		fmt.Println("出错", err)
	}
}
