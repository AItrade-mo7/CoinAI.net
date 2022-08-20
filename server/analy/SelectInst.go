package analy

import (
	"fmt"

	"github.com/EasyGolang/goTools/mCount"
	"github.com/EasyGolang/goTools/mOKX"
)

func SelectInst(AnalyKdata map[string][]mOKX.TypeKd) {
	for _, list := range AnalyKdata {
		if len(list) == 300 {
			SingleAnalyInst(list)
		}
	}
}

// 在这里判断趋势 并挑选币种
func SingleAnalyInst(list []mOKX.TypeKd) {
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
	fmt.Println("精度", precision)
	// 平均振幅
	HLPerAvg := mCount.Average(HLPerArr)
	// 最大振幅
	MaxHLPer := HLPerArr[0]
	for _, item := range HLPerArr {
		if mCount.Le(item, MaxHLPer) > 0 {
			MaxHLPer = item
		}
	}

	fmt.Println(list[0].InstID, "平均振幅", mCount.CentRound(HLPerAvg, precision))
	fmt.Println(list[0].InstID, "最大振幅", MaxHLPer)
}
