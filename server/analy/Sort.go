package analy

import (
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mCount"
)

// 按照 最大振幅 排序
func Sort_MaxHLPer(data []okxInfo.HLAnalySelectType) []okxInfo.HLAnalySelectType {
	size := len(data)
	list := make([]okxInfo.HLAnalySelectType, size)
	copy(list, data)

	var swapped bool
	for i := size - 1; i > 0; i-- {
		swapped = false
		for j := 0; j < i; j++ {
			a := list[j+1].DiffMaxAvg
			b := list[j].DiffMaxAvg
			if mCount.Le(a, b) < 0 {
				list[j], list[j+1] = list[j+1], list[j]
				swapped = true
			}
		}
		if !swapped {
			break
		}
	}
	// 设置 Idx 并翻转
	listIDX := []okxInfo.HLAnalySelectType{}
	j := 0
	for i := len(list) - 1; i > -1; i-- {
		Kdata := list[i]
		listIDX = append(listIDX, Kdata)
		j++
	}
	return listIDX
}
