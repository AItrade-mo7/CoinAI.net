package hunter

import (
	"fmt"

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

// RsiRegion 是否为降序
func Is_RsiRegion_GoDown(preArr []okxInfo.TradeKdType) bool {
	cacheArr := []int{}

	/*
		cacheArr 长度原本为 1
		当前与 cacheArr 最后一位 比对 ，一定相等 则 长度为2
		当前 与 cacheArr 最后一位 比对 , 不一定相等
				若此时大于等于 则 长度为 3
				若此时  小  于 则 长度为2 并结束
	*/
	for i := len(preArr) - 1; i >= 0; i-- {
		item := preArr[i]
		if len(cacheArr) < 1 {
			cacheArr = append(cacheArr, item.RsiRegion)
		}

		cacheLast := cacheArr[len(cacheArr)-1]

		if item.RsiRegion >= cacheLast {
			cacheArr = append(cacheArr, item.RsiRegion)
		} else {
			break
		}
	}
	if len(cacheArr) < 3 {
		return false
	}
	cacheLast := cacheArr[len(cacheArr)-1]
	cacheFirst := cacheArr[0]

	if cacheLast-cacheFirst > 0 {
		return true
	}

	return false
}

// RsiRegion 是否为升序
func Is_RsiRegion_GoUp(preArr []okxInfo.TradeKdType) bool {
	cacheArr := []int{}
	/*
		cacheArr 长度原本为 1
		当前与 cacheArr 最后一位 比对 ，一定相等 则 长度为2
		当前 与 cacheArr 最后一位 比对 , 不一定相等
				若此时 则 长度为 3
				若此时 小  于 则 长度为2 并结束
	*/
	for i := len(preArr) - 1; i >= 0; i-- {
		item := preArr[i]
		if len(cacheArr) < 1 {
			cacheArr = append(cacheArr, item.RsiRegion)
		}

		cacheLast := cacheArr[len(cacheArr)-1]

		if item.RsiRegion <= cacheLast {
			cacheArr = append(cacheArr, item.RsiRegion)
		} else {
			break
		}
	}
	if len(cacheArr) < 3 {
		return false
	}
	cacheLast := cacheArr[len(cacheArr)-1]
	cacheFirst := cacheArr[0]

	if cacheLast-cacheFirst < 0 {
		return true
	}

	return false
}

// preArr 的 RsiRegion 是否有大于2  的存在
func Is_RsiRegion_Gte2(preArr []okxInfo.TradeKdType) (result bool) {
	result = false
	if len(preArr) < 3 {
		return
	}

	for i := len(preArr) - 1; i >= 0; i-- {
		RsiRegion := preArr[i].RsiRegion
		valAbs := mCount.Abs(fmt.Sprint(RsiRegion))
		if mCount.Le(valAbs, "2") >= 0 {
			result = true
			break
		}
	}

	return
}

// CAP_EMA 是否为升序
func Is_CAP_EMA_GoUp(preArr []okxInfo.TradeKdType) bool {
	cacheArr := []string{}
	/*
		cacheArr 长度原本为 1
		当前与 cacheArr 最后一位 比对 ，一定相等 则 长度为2
		当前 与 cacheArr 最后一位 比对 , 不一定相等
				若此时 则 长度为 3
				若此时 小  于 则 长度为2 并结束
	*/
	for i := len(preArr) - 1; i >= 0; i-- {
		item := preArr[i]
		if len(cacheArr) < 1 {
			cacheArr = append(cacheArr, item.CAP_EMA)
		}
		cacheLast := cacheArr[len(cacheArr)-1]

		if mCount.Le(item.CAP_EMA, cacheLast) <= 0 {
			cacheArr = append(cacheArr, item.CAP_EMA)
		} else {
			break
		}
	}
	if len(cacheArr) < 4 {
		return false
	}
	cacheLast := cacheArr[len(cacheArr)-1]
	cacheFirst := cacheArr[0]

	if mCount.Le(cacheLast, cacheFirst) < 0 {
		return true
	}

	return false
}

// CAP_EMA 是否为升序
func Is_CAP_EMA_GoDown(preArr []okxInfo.TradeKdType) bool {
	cacheArr := []string{}
	/*
		cacheArr 长度原本为 1
		当前与 cacheArr 最后一位 比对 ，一定相等 则 长度为2
		当前 与 cacheArr 最后一位 比对 , 不一定相等
				若此时 则 长度为 3
				若此时 小  于 则 长度为2 并结束
	*/
	for i := len(preArr) - 1; i >= 0; i-- {
		item := preArr[i]
		if len(cacheArr) < 1 {
			cacheArr = append(cacheArr, item.CAP_EMA)
		}
		cacheLast := cacheArr[len(cacheArr)-1]

		if mCount.Le(item.CAP_EMA, cacheLast) >= 0 {
			cacheArr = append(cacheArr, item.CAP_EMA)
		} else {
			break
		}
	}
	if len(cacheArr) < 4 {
		return false
	}
	cacheLast := cacheArr[len(cacheArr)-1]
	cacheFirst := cacheArr[0]

	if mCount.Le(cacheLast, cacheFirst) > 0 {
		return true
	}

	return false
}
