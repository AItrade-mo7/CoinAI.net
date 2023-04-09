package analyConfig

import (
	"os"

	"CoinAI.net/task/taskStart"
	"CoinAI.net/task/testHunter"
	"github.com/EasyGolang/goTools/mCount"
	"github.com/EasyGolang/goTools/mFile"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mStr"
	jsoniter "github.com/json-iterator/go"
)

func GetWinArr(opt taskStart.BackReturn) {
	var file []byte

	if len(opt.BillingPath) > 1 {
		file, _ = os.ReadFile(opt.BillingPath)
	}

	var BillingArr []testHunter.BillingType // 数据来源
	jsoniter.Unmarshal(file, &BillingArr)

	if len(opt.BillingArr) > 2 {
		BillingArr = opt.BillingArr
	}

	// 胜率 最高来排序
	WinArr := MoneySort(BillingArr)

	// 取出来最后 5 个
	WinArr = WinArr[len(WinArr)-5:]

	resultPath := mStr.Join(opt.ResultBasePath, "/", opt.InstID, "-WinArr.json")

	// 排序结果分析
	mFile.Write(resultPath, mJson.ToStr(WinArr))
}

func MoneySort(arr []testHunter.BillingType) []testHunter.BillingType {
	size := len(arr)
	list := make([]testHunter.BillingType, size)
	copy(list, arr)

	var swapped bool
	for i := size - 1; i > 0; i-- {
		swapped = false
		for j := 0; j < i; j++ {
			a := list[j+1].WinRatio
			b := list[j].WinRatio
			if mCount.Le(a, b) < 0 {
				list[j], list[j+1] = list[j+1], list[j]
				swapped = true
			}
		}
		if !swapped {
			break
		}
	}
	return list
}
