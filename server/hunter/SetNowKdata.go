package hunter

import (
	"fmt"

	"CoinAI.net/server/global"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mTime"
)

func SetNowKdata() error {
	NowList := mOKX.GetKdata(mOKX.GetKdataOpt{
		InstID: okxInfo.KdataInst.InstID,
		Page:   0,
	})

	if len(NowList) < 100 || len(okxInfo.NowKdataList) < 200 {
		err := fmt.Errorf("hunter.SetNowKdata 数据不足")
		return err
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
			okxInfo.NowKdataList = append(okxInfo.NowKdataList, NowItem)
		}

	}

	if len(okxInfo.NowKdataList)-okxInfo.MaxLen > 0 {
		okxInfo.NowKdataList = okxInfo.NowKdataList[len(okxInfo.NowKdataList)-okxInfo.MaxLen:]
	}

	// 数据检查
	var err error = nil
	for key, val := range okxInfo.NowKdataList {
		preIndex := key - 1
		if preIndex < 0 {
			preIndex = 0
		}
		preItem := okxInfo.NowKdataList[preIndex]
		nowItem := okxInfo.NowKdataList[key]
		if key > 0 {
			if nowItem.TimeUnix-preItem.TimeUnix != mTime.UnixTimeInt64.Hour {
				okxInfo.NowKdataList = []mOKX.TypeKd{} // 清空历史数据
				err = fmt.Errorf("数据检查出错, 系统正在自行恢复 %+v %+v %+v", val.InstID, val.TimeStr, key)
				break
			}
		}
	}

	Last := okxInfo.NowKdataList[len(okxInfo.NowKdataList)-1]
	global.TradeLog.Println("更新一次最新数据: ", Last.InstID, Last.TimeStr, len(okxInfo.NowKdataList), Last.C)

	return err
}
