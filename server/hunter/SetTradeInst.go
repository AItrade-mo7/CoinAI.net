package hunter

import (
	"fmt"

	"CoinAI.net/server/global"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mJson"
)

func (_this *HunterObj) SetTradeInst() (resErr error) {
	resErr = nil
	if len(okxInfo.NowTicker.TickerVol) < 4 {
		resErr = fmt.Errorf("hunter.SetTradeInst TickerVol %+v", len(okxInfo.NowTicker.TickerVol))
		return
	}
	if len(okxInfo.NowTicker.AnalyWhole) < 3 {
		resErr = fmt.Errorf("hunter.SetTradeInst AnalyWhole %+v", len(okxInfo.NowTicker.AnalyWhole))
		return
	}
	if len(okxInfo.NowTicker.AnalySingle) < 4 {
		resErr = fmt.Errorf("hunter.SetTradeInst AnalySingle %+v", len(okxInfo.NowTicker.AnalySingle))
		return
	}

	if len(okxInfo.NowTicker.MillionCoin) < 2 {
		resErr = fmt.Errorf("hunter.SetTradeInst MillionCoin %+v", len(okxInfo.NowTicker.MillionCoin))
		return
	}
	// 在这里按照涨幅的绝对值排个序SetTradeInst

	HLPerList := Sort_HLPer(okxInfo.NowTicker.MillionCoin)
	HLPerInstID := "ETH-USDT"
	// 在这里 取出 最末尾的三个 然后取第一个

	InstIDList := []string{}
	rangeCount := 0
	for i := len(HLPerList) - 1; i >= 0; i-- {
		item := HLPerList[i]
		rangeCount++
		if rangeCount > _this.HLPerLevel {
			break
		}
		InstIDList = append(InstIDList, item.InstID)
	}

	if len(InstIDList) > 0 {
		HLPerInstID = InstIDList[len(InstIDList)-1]
	}

	CoinId := HLPerInstID

	if len(CoinId) < 1 {
		resErr = fmt.Errorf("hunter.SetTradeInst CoinId %+v", CoinId)
		return
	}

	KdataInst := okxInfo.Inst[CoinId]
	TradeInst := okxInfo.Inst[CoinId+"-SWAP"]

	if KdataInst.State == "live" &&
		TradeInst.State == "live" &&
		len(TradeInst.InstID) > 1 &&
		len(KdataInst.InstID) > 1 &&
		KdataInst.InstType == "SPOT" &&
		TradeInst.InstType == "SWAP" {
	} else {
		resErr = fmt.Errorf(
			"hunter.SetTradeInst KdataInst:%+v TradeInst:%+v",
			mJson.Format(KdataInst),
			mJson.Format(TradeInst),
		)
		return
	}

	_this.KdataInst = KdataInst
	_this.TradeInst = TradeInst

	global.TradeLog.Println(_this.HunterName, "hunter.SetTradeInst", InstIDList, CoinId)

	return
}
