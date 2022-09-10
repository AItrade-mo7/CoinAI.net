package hunter

import (
	"CoinAI.net/server/global"
	"CoinAI.net/server/okxApi/wssApi/positions"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mJson"
)

/*

监听 K 线 并 执行 下单  平仓 操作

*/

func Start() {
	global.RunLog.Println("======== Hunter =========")

	global.RunLog.Println("监听持仓频道")
	go positions.Start()

	global.RunLog.Println("加载设置", mJson.Format(okxInfo.HunterRun))
}
