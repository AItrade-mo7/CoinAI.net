package hunter

import (
	"fmt"

	"CoinAI.net/server/okxInfo"
)

func SetTradeInst() (resErr error) {
	resErr = nil
	if len(okxInfo.NowTicker.TickerVol) < 4 {
		resErr = fmt.Errorf("榜单数量不足 ready.SetTradeInst TickerVol %+v", len(okxInfo.NowTicker.TickerVol))
		return
	}
	if len(okxInfo.NowTicker.AnalyWhole) < 3 {
		resErr = fmt.Errorf("切片数量不足 ready.SetTradeInst AnalyWhole %+v", len(okxInfo.NowTicker.AnalyWhole))
		return
	}
	if len(okxInfo.NowTicker.AnalySingle) < 4 {
		resErr = fmt.Errorf("币种分析样本不足 ready.SetTradeInst AnalySingle %+v", len(okxInfo.NowTicker.AnalySingle))
		return
	}

	if len(okxInfo.NowTicker.MillionCoin) < 2 {
		resErr = fmt.Errorf("数据异常 ready.SetTradeInst MillionCoin %+v", len(okxInfo.NowTicker.MillionCoin))
		return
	}
	// 在这里按照涨幅的绝对值排个序SetTradeInst

	HLPerList := Sort_HLPer(okxInfo.NowTicker.MillionCoin)

	CoinId := HLPerList[len(HLPerList)-1].InstID

	for _, val := range HLPerList {
		if val.InstID == "ETH-USDT" || val.InstID == "BTC-USDT" {
			CoinId = val.InstID
			break
		}
	}

	if len(CoinId) < 0 {
		resErr = fmt.Errorf("数据异常 ready.SetTradeInst CoinId %+v", CoinId)
		return
	}

	KdataInst := okxInfo.Inst[CoinId]
	TradeInst := okxInfo.Inst[CoinId+"-SWAP"]

	if KdataInst.State == "live" && TradeInst.State == "live" {
	} else {
		resErr = fmt.Errorf("数据异常 ready.SetTradeInst State %+v", KdataInst.State)
		return
	}

	if KdataInst.InstType == "SPOT" && TradeInst.InstType == "SWAP" {
	} else {
		resErr = fmt.Errorf("数据异常 ready.SetTradeInst InstType %+v", KdataInst.InstType)
		return
	}

	okxInfo.KdataInst = KdataInst
	okxInfo.TradeInst = TradeInst

	return
}
