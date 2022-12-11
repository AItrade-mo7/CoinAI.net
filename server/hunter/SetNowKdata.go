package hunter

import (
	"CoinAI.net/server/global"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mTime"
)

func SetNowKdata() {
	NowList := mOKX.GetKdata(mOKX.GetKdataOpt{
		InstID: okxInfo.KdataInst.InstID,
		Page:   0,
		After:  mTime.GetUnixInt64(),
	})

	if len(NowList) < 100 || len(okxInfo.NowKdata) < 100 {
		global.LogErr("hunter.SetNowKdata 数据不足")
		return
	}

	for _, NowItem := range NowList {
		Fund := false
		FundKey := 0
		for key, Item := range okxInfo.NowKdata {
			if NowItem.TimeUnix == Item.TimeUnix { // 相等的直接替换
				Fund = true
				FundKey = key
				break
			}
		}

		if Fund {
			global.Log.Println("替换", FundKey)
			okxInfo.NowKdata[FundKey] = NowItem
		} else {
			global.Log.Println("新增")
			okxInfo.NowKdata = append(okxInfo.NowKdata, NowItem)
		}
	}

	// 数据检查
	for key, val := range okxInfo.NowKdata {
		preIndex := key - 1
		if preIndex < 0 {
			preIndex = 0
		}
		preItem := okxInfo.NowKdata[preIndex]
		nowItem := okxInfo.NowKdata[key]
		global.Log.Println(val.TimeStr, nowItem.TimeUnix-preItem.TimeUnix, key)
	}
}
