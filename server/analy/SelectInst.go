package analy

import (
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mCount"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mOKX"
)

func SelectInst(AnalyKdata map[string][]mOKX.TypeKd) {
	AnalySelect := []okxInfo.AnalySelectType{}
	for _, list := range AnalyKdata {
		if len(list) == 300 {
			result := SingleAnalyInst(list)
			AnalySelect = append(AnalySelect, result)
		}
	}

	// 振幅与平均振幅
	mJson.Println(AnalySelect)
}

// 在这里判断趋势 并挑选币种
func SingleAnalyInst(list []mOKX.TypeKd) (resData okxInfo.AnalySelectType) {
	resData = okxInfo.AnalySelectType{}
	// 截取最近 8 小时振幅 , 15分钟一格， 32 个格子
	listLen := len(list)
	HLPerArr := []string{}
	var precision int32
	for key, item := range list {
		precision = mCount.GetDecimal(item.TickSz)
		if key > listLen-38 {
			HLPerArr = append(HLPerArr, item.HLPer)
		}
	}
	// 平均振幅
	HLPerAvg := mCount.Average(HLPerArr)
	// 最大振幅
	MaxHLPer := HLPerArr[0]
	for _, item := range HLPerArr {
		if mCount.Le(item, MaxHLPer) > 0 {
			MaxHLPer = item
		}
	}
	resData.MaxHLPer = MaxHLPer
	resData.HLPerAvg = mCount.CentRound(HLPerAvg, precision)
	resData.DiffMaxAvg = mCount.Sub(MaxHLPer, resData.HLPerAvg)
	resData.InstID = list[0].InstID
	return
}
