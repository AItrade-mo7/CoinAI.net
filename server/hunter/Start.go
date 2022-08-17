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
	global.HunterLog.Println("加载设置", mJson.Format(okxInfo.HunterRun))

	global.HunterLog.Println("设置持仓模式")

	global.HunterLog.Println("登录监听持仓频道") // WSS

	global.HunterLog.Println("监听 大盘 K 线") // rest 监听

	global.HunterLog.Println("监听 交易 K 线") // rest 监听
}
