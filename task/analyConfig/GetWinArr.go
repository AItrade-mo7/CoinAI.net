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

	// 收益最高的排序
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
			a := list[j+1].Money
			b := list[j].Money
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

/*

2021-01
2022-01
最优

BTC
EMA_86_CAP_2_level_1
EMA_84_CAP_2_level_1
EMA_80_CAP_2_level_1
EMA_88_CAP_4_level_1
EMA_82_CAP_5_level_1
EMA_88_CAP_5_level_1

//参数 范围 82 83  84  85 86  87 88

ETH
EMA_78_CAP_2_level_1
EMA_80_CAP_2_level_1
EMA_74_CAP_3_level_1
EMA_76_CAP_2_level_1

//参数 范围 76 77  78  79  80  81 82

接下来是带上杠杆，求最优值

*/
