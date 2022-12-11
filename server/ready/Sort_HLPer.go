package ready

import (
	"github.com/EasyGolang/goTools/mCount"
	"github.com/EasyGolang/goTools/mOKX"
)

// 按照 平均振幅 排序
func Sort_HLPer(data []mOKX.AnalySliceType) []mOKX.AnalySliceType {
	size := len(data)
	list := make([]mOKX.AnalySliceType, size)
	copy(list, data)

	var swapped bool
	for i := size - 1; i > 0; i-- {
		swapped = false
		for j := 0; j < i; j++ {
			a := list[j+1].HLPerAvg
			b := list[j].HLPerAvg
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
	listIDX := []mOKX.AnalySliceType{}
	j := 0
	for i := len(list) - 1; i > -1; i-- {
		Kdata := list[i]
		listIDX = append(listIDX, Kdata)
		j++
	}
	return listIDX
}
