package hunter

import (
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mCount"
	"github.com/EasyGolang/goTools/mOKX"
)

func CAPIdxToText(idx int) string {
	if idx > 0 {
		return "Buy"
	}
	if idx < 0 {
		return "Sell"
	}
	return "nil"
}

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

func GetCAPIdx(now okxInfo.TradeKdType) int {
	now_EMA_diff := mCount.Le(now.CAP_EMA, "0") // 1 0 -1  EMA
	now_MA_diff := mCount.Le(now.CAP_MA, "0")   // -1 0 1  MA

	nowDiff := now_EMA_diff
	if now_MA_diff == now_EMA_diff {
		nowDiff = now_MA_diff + now_EMA_diff
	}

	return nowDiff
}

/*
3   大于70  超买区
2   60-70   多买区
1   50-60   上震荡区
-1  40-50   下震荡区
-2  30-40   多卖区
-3  小于 30  超卖区
*/
func GetRsiRegion(now okxInfo.TradeKdType) int {
	RSI := now.RSI_18
	// 1 50-60
	if mCount.Le(RSI, "50") > 0 && mCount.Le(RSI, "60") <= 0 {
		return 1
	}

	// 2 60-70
	if mCount.Le(RSI, "60") > 0 && mCount.Le(RSI, "70") < 0 {
		return 2
	}

	// 3 大于70
	if mCount.Le(RSI, "70") >= 0 {
		return 3
	}

	if mCount.Le(RSI, "50") == 0 {
		return 0
	}

	// -1 40-50
	if mCount.Le(RSI, "40") >= 0 && mCount.Le(RSI, "50") < 0 {
		return -1
	}

	// -2 30-40
	if mCount.Le(RSI, "30") > 0 && mCount.Le(RSI, "40") < 0 {
		return -2
	}

	// -3 30-40
	if mCount.Le(RSI, "30") <= 0 {
		return -3
	}

	return 0
}

