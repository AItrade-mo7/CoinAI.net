package hunter

import (
	"CoinAI.net/server/global"
	"github.com/EasyGolang/goTools/mTime"
)

func Start() {
	global.RunLog.Println("hunter.Start 执行", mTime.UnixFormat(mTime.GetUnixInt64()))
}
