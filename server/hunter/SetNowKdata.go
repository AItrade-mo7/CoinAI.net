package hunter

import (
	"fmt"

	"CoinAI.net/server/global"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mTime"
)

func (_this *HunterObj) SetNowKdata() error {
	NowList := mOKX.GetKdata(mOKX.GetKdataOpt{
		InstID: _this.KdataInst.InstID,
		Page:   0,
	})

	if len(NowList) < 100 || len(_this.NowKdataList) < 200 {
		err := fmt.Errorf(_this.HunterName, "hunter.SetNowKdata 数据不足")
		return err
	}

	fmt.Println(1111, len(_this.NowKdataList))

	for _, NowItem := range NowList {
		Fund := false
		FundKey := 0

		for key, Item := range _this.NowKdataList {
			if NowItem.TimeUnix == Item.TimeUnix { // 相等的直接替换
				Fund = true
				FundKey = key
				break
			}
		}

		if Fund {
			_this.NowKdataList[FundKey] = NowItem
		} else {
			_this.NowKdataList = append(_this.NowKdataList, NowItem)
		}

	}

	fmt.Println(2222, len(_this.NowKdataList))

	if len(_this.NowKdataList)-_this.MaxLen > 0 {
		_this.NowKdataList = _this.NowKdataList[len(_this.NowKdataList)-_this.MaxLen:]
	}

	fmt.Println(333, len(_this.NowKdataList))

	// 数据检查
	var err error = nil
	for key, val := range _this.NowKdataList {
		preIndex := key - 1
		if preIndex < 0 {
			preIndex = 0
		}
		preItem := _this.NowKdataList[preIndex]
		nowItem := _this.NowKdataList[key]
		if key > 0 {
			if nowItem.TimeUnix-preItem.TimeUnix != mTime.UnixTimeInt64.Hour {
				_this.NowKdataList = []mOKX.TypeKd{} // 清空历史数据
				err = fmt.Errorf(_this.HunterName, "数据检查出错, 系统正在自行恢复 %+v %+v %+v", val.InstID, val.TimeStr, key)
				break
			}
		}
	}

	Last := _this.NowKdataList[len(_this.NowKdataList)-1]
	global.TradeLog.Println(_this.HunterName, "更新一次最新数据: ", Last.InstID, Last.TimeStr, len(_this.NowKdataList), Last.C)

	return err
}
