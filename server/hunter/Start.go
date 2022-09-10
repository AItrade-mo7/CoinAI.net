package hunter

import (
	"CoinAI.net/server/global"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mJson"
)

/*

监听 K 线 并 执行 下单  平仓 操作

*/

func Start() {
	global.RunLog.Println("加载设置", mJson.Format(okxInfo.HunterRun))

	global.RunLog.Println("设置持仓模式")

	global.RunLog.Println("登录监听持仓频道") // WSS

	global.RunLog.Println("监听 大盘 K 线") // rest 监听

	global.RunLog.Println("监听 交易 K 线") // rest 监听
}
