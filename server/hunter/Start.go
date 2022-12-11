package hunter

import (
	"fmt"

	"CoinAI.net/server/global"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mTime"
)

func Start() {
	for ok := range okxInfo.Ticking {
		global.RunLog.Println("hunter.Start 执行", mTime.UnixFormat(mTime.GetUnixInt64()), ok)
		Running()
	}
}

func Running() {
	// 检测数据
	if len(okxInfo.TradeInst.InstID) < 2 {
		global.LogErr("hunter.Running", "okxInfo.TradeInst.InstID 为空")
		return
	}

	FileBaseKdata()
}

func FileBaseKdata() {
	if len(okxInfo.TradeKdata) < 100 {
		// 回填历史数据 1 组
		List := mOKX.GetKdata(mOKX.GetKdataOpt{
			InstID: okxInfo.TradeInst.InstID,
			After:  mTime.GetUnixInt64(),
			Page:   1,
		})

		for _, val := range List {
			fmt.Println(val.TimeStr)
		}
	}
}
