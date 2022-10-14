package mOKX

import (
	"github.com/EasyGolang/goTools/mCount"
)

// 成交量排序
func SortVolume(data []TypeTicker) []TypeTicker {
	size := len(data)
	list := make([]TypeTicker, size)
	copy(list, data)

	var swapped bool
	for i := size - 1; i > 0; i-- {
		swapped = false
		for j := 0; j < i; j++ {
			a := list[j+1].Volume
			b := list[j].Volume
			if mCount.Le(a, b) < 0 {
				list[j], list[j+1] = list[j+1], list[j]
				swapped = true
			}
		}
		if !swapped {
			break
		}
	}

	// 设置 VolIdx 并翻转
	listIDX := []TypeTicker{}
	j := 0
	for i := len(list) - 1; i > -1; i-- {
		Kdata := list[i]
		listIDX = append(listIDX, Kdata)
		j++
	}

	return listIDX
}
