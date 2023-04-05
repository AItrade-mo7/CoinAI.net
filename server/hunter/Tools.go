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

/*
3   大于70  超买区
2   60-70   多买区
1   50-60   上震荡区
-1  40-50   下震荡区
-2  30-40   多卖区
-3  小于 30  超卖区

func GetRsiRegion(now okxInfo.TradeKdType) int {
	RSI_EMA := now.RSI_EMA
	// 1 50-60
	if mCount.Le(RSI_EMA, "50") > 0 && mCount.Le(RSI_EMA, "60") <= 0 {
		return 1
	}

	// 2 60-70
	if mCount.Le(RSI_EMA, "60") > 0 && mCount.Le(RSI_EMA, "70") < 0 {
		return 2
	}

	// 3 大于70
	if mCount.Le(RSI_EMA, "70") >= 0 {
		return 3
	}

	if mCount.Le(RSI_EMA, "50") == 0 {
		return 0
	}

	// -1 40-50
	if mCount.Le(RSI_EMA, "40") >= 0 && mCount.Le(RSI_EMA, "50") < 0 {
		return -1
	}

	// -2 30-40
	if mCount.Le(RSI_EMA, "30") > 0 && mCount.Le(RSI_EMA, "40") < 0 {
		return -2
	}

	// -3 30-40
	if mCount.Le(RSI_EMA, "30") <= 0 {
		return -3
	}

	return 0
}
*/
/*
// RsiRegion EMA 是否为降序
func Is_RsiRegion_GoDown(preArr []okxInfo.TradeKdType) []int {
	cacheArr := []int{}
	downArr := []int{}

		cacheArr 长度原本为 1
		当前与 cacheArr 最后一位 比对 ，一定相等 则 长度为2
		当前 与 cacheArr 最后一位 比对 , 不一定相等
				若此时大于等于 则 长度为 3
				若此时  小  于 则 长度为2 并结束
	for i := len(preArr) - 1; i >= 0; i-- {
		item := preArr[i]
		if len(cacheArr) < 1 {
			cacheArr = append(cacheArr, item.RsiEmaRegion)
		}

		cacheLast := cacheArr[len(cacheArr)-1]
		if item.RsiEmaRegion >= cacheLast {
			cacheArr = append(cacheArr, item.RsiEmaRegion)

			if item.RsiEmaRegion > cacheArr[0] {
				downArr = append(downArr, item.RsiEmaRegion)
			}
		} else {
			break
		}
	}

	return downArr
}
*/

/*
// RsiRegion EMA 是否为升序
func Is_RsiRegion_GoUp(preArr []okxInfo.TradeKdType) []int {
	cacheArr := []int{}
	upArr := []int{}
		cacheArr 长度原本为 1
		当前与 cacheArr 最后一位 比对 ，一定相等 则 长度为2
		当前 与 cacheArr 最后一位 比对 , 不一定相等
				若此时 则 长度为 3
				若此时 小  于 则 长度为2 并结束
	for i := len(preArr) - 1; i >= 0; i-- {
		item := preArr[i]
		if len(cacheArr) < 1 {
			cacheArr = append(cacheArr, item.RsiEmaRegion)
		}

		cacheLast := cacheArr[len(cacheArr)-1]

		if item.RsiEmaRegion <= cacheLast {
			cacheArr = append(cacheArr, item.RsiEmaRegion)

			if item.RsiEmaRegion < cacheArr[0] {
				upArr = append(upArr, item.RsiEmaRegion)
			}
		} else {
			break
		}
	}

	return upArr
}
*/
/*
// preArr 的 RsiRegion 是否有大于2  的存在
func Is_RsiRegion_Gte2(preArr []okxInfo.TradeKdType) (result bool) {
	result = false
	if len(preArr) < 3 {
		return
	}

	for i := len(preArr) - 1; i >= 0; i-- {
		RsiRegion := preArr[i].RsiEmaRegion
		valAbs := mCount.Abs(fmt.Sprint(RsiRegion))
		if mCount.Le(valAbs, "2") >= 0 {
			result = true
			break
		}
	}

	return
}
*/

// CAP_EMA 是否为升序
func Is_CAP_EMA_GoUp(preArr []okxInfo.TradeKdType) []string {
	cacheArr := []string{}
	upArr := []string{}

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
			if mCount.Le(item.CAP_EMA, cacheArr[0]) < 0 {
				upArr = append(upArr, item.CAP_EMA)
			}
		} else {
			break
		}
	}

	return upArr
}

// CAP_EMA 是否为升序
func Is_CAP_EMA_GoDown(preArr []okxInfo.TradeKdType) []string {
	cacheArr := []string{}
	downArr := []string{}

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
			if mCount.Le(item.CAP_EMA, cacheArr[0]) > 0 {
				downArr = append(downArr, item.CAP_EMA)
			}
		} else {
			break
		}
	}

	return downArr
}

func GetCAPIdx(now okxInfo.TradeKdType) int {
	now_EMA_diff := mCount.Le(now.CAP_EMA, "0") // 1 0 -1  EMA
	now_MA_diff := mCount.Le(now.CAP_MA, "0")   // -1 0 1  MA
	nowDiff := now_EMA_diff

	if now_MA_diff == now_EMA_diff {
		nowDiff = now_MA_diff + now_EMA_diff
	}

	// diffAdd := mCount.Add(now.CAP_EMA, now.CAP_MA)
	// nowDiff := mCount.Le(diffAdd, "0")

	return nowDiff
}
