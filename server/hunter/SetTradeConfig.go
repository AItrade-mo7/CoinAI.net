package hunter

import (
	"fmt"

	"CoinAI.net/server/global"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mJson"
)

// 设置最优参数
func (_this *HunterObj) SetTradeConfig() (resErr error) {
	resErr = nil

	LastKdata := _this.NowKdataList[len(_this.NowKdataList)-1]

	if len(LastKdata.InstID) < 1 || LastKdata.InstID != _this.InstID {
		resErr = fmt.Errorf(_this.HunterName+" hunter.SetTradeConfig InstID 不正确 %+v", LastKdata.InstID)
		return
	}

	_this.TradeKdataOpt = okxInfo.CoinTradeConfig[LastKdata.InstID]

	global.TradeLog.Println(_this.HunterName, "已设置交易参数", mJson.ToStr(_this.TradeKdataOpt))

	return
}
