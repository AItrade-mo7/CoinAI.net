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
			okxInfo.NowKdata[FundKey] = NowItem
		} else {
			global.RunLog.Println("新增")
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
		if key > 0 {
			if nowItem.TimeUnix-preItem.TimeUnix != mTime.UnixTimeInt64.Hour {
				global.LogErr("数据检查出错, 系统正在自行恢复", val.InstID, val.TimeStr, key)
				okxInfo.NowKdata = []mOKX.TypeKd{} // 清空历史数据
				Running()                          // 立即重新执行一次 Running
				break
			}
		}
	}
}
