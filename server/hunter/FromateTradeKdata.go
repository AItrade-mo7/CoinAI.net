package hunter

import (
	"fmt"

	"CoinAI.net/server/global"
	"CoinAI.net/server/okxInfo"
)

func FormatTradeKdata() {
	if len(okxInfo.NowKdata) < 100 {
		global.LogErr("hunter.FormatTradeKdata 数据不足")
		return
	}

	for _, Kdata := range okxInfo.NowKdata {
		fmt.Println(Kdata.TimeUnix)
	}
}
