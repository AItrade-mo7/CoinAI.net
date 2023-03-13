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

	if len(NowList) < 100 || len(okxInfo.NowKdataList) < 100 {
		global.LogErr("hunter.SetNowKdata 数据不足")
		return
	}

	for _, NowItem := range NowList {
		Fund := false
		FundKey := 0
		for key, Item := range okxInfo.NowKdataList {
			if NowItem.TimeUnix == Item.TimeUnix { // 相等的直接替换
				Fund = true
				FundKey = key
				break
			}
		}

		if Fund {
			okxInfo.NowKdataList[FundKey] = NowItem
		} else {
			// global.RunLog.Println("新增")
			okxInfo.NowKdataList = append(okxInfo.NowKdataList, NowItem)
		}
	}

	// if len(okxInfo.NowKdataList)-okxInfo.MaxLen > 0 {
	// 	okxInfo.NowKdataList = okxInfo.NowKdataList[len(okxInfo.NowKdataList)-okxInfo.MaxLen:]
	// 	global.RunLog.Println("长度超出，裁剪", len(okxInfo.NowKdataList))
	// }

	// 数据检查
	for key, val := range okxInfo.NowKdataList {
		preIndex := key - 1
		if preIndex < 0 {
			preIndex = 0
		}
		preItem := okxInfo.NowKdataList[preIndex]
		nowItem := okxInfo.NowKdataList[key]
		if key > 0 {
			if nowItem.TimeUnix-preItem.TimeUnix != mTime.UnixTimeInt64.Hour {
				global.LogErr("数据检查出错, 系统正在自行恢复", val.InstID, val.TimeStr, key)
				okxInfo.NowKdataList = []mOKX.TypeKd{} // 清空历史数据
				Running()                              // 立即重新执行一次 Running
				break
			}
		}
	}

	// Last := okxInfo.NowKdataList[len(okxInfo.NowKdataList)-1]

	// global.RunLog.Println("更新一次最新数据", Last.TimeStr, Last.C, len(okxInfo.NowKdataList))
}
