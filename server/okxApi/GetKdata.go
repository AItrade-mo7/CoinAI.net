package okxApi

import (
	"fmt"

	"CoinAI.net/server/global"
	"CoinAI.net/server/okxApi/binanceApi"
	"CoinAI.net/server/okxApi/restApi/kdata"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mCount"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mTime"
)

type GetKdataOpt struct {
	InstID  string `bson:"InstID"`
	Current int    `bson:"Current"` // 当前页码 0 为
	After   int64  `bson:"After"`   // 时间 默认为当前时间
	Size    int    `bson:"Size"`    // 数量 默认为100
}

func GetKdata(opt GetKdataOpt) (KdataList []mOKX.TypeKd) {
	KdataList = []mOKX.TypeKd{}
	SPOT := okxInfo.Inst[opt.InstID]
	if len(SPOT.InstID) < 3 {
		return
	}

	nowUnix := mTime.GetUnixInt64() - mTime.UnixTimeInt64.Minute*16
	if opt.After > nowUnix {
		opt.After = 0 // 当前
	} else {
		// 历史
		if opt.Size > 100 {
			opt.Size = 100
		}
	}

	BinanceList := binanceApi.GetKdata(binanceApi.GetKdataParam{
		Symbol:  SPOT.Symbol,
		Current: opt.Current,
		After:   opt.After,
		Size:    opt.Size,
	})
	if len(BinanceList) != opt.Size {
		global.LogErr("BinanceList 长度不正确", len(BinanceList), mJson.Format(opt))
		return
	}

	var OKXList []mOKX.TypeKd
	if (opt.After) > 0 || opt.Current > 0 {
		OKXList = kdata.GetHistoryKdata(kdata.HistoryKdataParam{
			InstID:  SPOT.InstID,
			Current: opt.Current,
			After:   opt.After,
			Size:    opt.Size,
		})
	} else {
		OKXList = kdata.GetKdata(SPOT.InstID, opt.Size)
	}

	if len(OKXList) != opt.Size {
		global.LogErr("OKXList 未获取到数据", len(OKXList), mJson.Format(opt))
		return
	}

	List, err := KdataMerge(KdataMergeOpt{
		OKXList:     OKXList,
		BinanceList: BinanceList,
	})
	if err != nil {
		global.LogErr(err, mJson.Format(opt))
		return
	}

	KdataList = List

	return
}

type KdataMergeOpt struct {
	OKXList     []mOKX.TypeKd
	BinanceList []mOKX.TypeKd
}

func KdataMerge(opt KdataMergeOpt) (Kdata []mOKX.TypeKd, resErr error) {
	OKXList := opt.OKXList
	BinanceList := opt.BinanceList
	Kdata = []mOKX.TypeKd{}
	resErr = nil

	if len(OKXList) != len(BinanceList) {
		resErr = fmt.Errorf("okxApi.KdataMerge len %+v %+v %+v", len(OKXList), len(BinanceList), opt)
		return
	}

	if OKXList[len(OKXList)-1].TimeStr != BinanceList[len(BinanceList)-1].TimeStr {
		global.Log.Println(
			"Merge [last]",
			OKXList[len(OKXList)-1].TimeStr,
			BinanceList[len(BinanceList)-1].TimeStr,
			"Merge [0]",
			OKXList[0].TimeStr,
			BinanceList[0].TimeStr,
		)
	}

	for _, item := range OKXList {
		OkxItem := item

		for _, BinanceItem := range BinanceList {
			if OkxItem.TimeUnix == BinanceItem.TimeUnix {
				VolCcy := mCount.Add(BinanceItem.VolCcy, OkxItem.VolCcy)
				OkxItem.VolCcy = VolCcy
				Vol := mCount.Add(BinanceItem.Vol, OkxItem.Vol)
				OkxItem.Vol = Vol
				OkxItem.DataType = "Merge"
				break
			}
		}

		Kdata = append(Kdata, OkxItem)
	}

	return
}
