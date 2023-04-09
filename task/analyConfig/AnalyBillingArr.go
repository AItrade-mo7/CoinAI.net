package analyConfig

import (
	"fmt"
	"os"

	"CoinAI.net/task/testHunter"
	"github.com/EasyGolang/goTools/mCount"
	"github.com/EasyGolang/goTools/mFile"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mStr"
	jsoniter "github.com/json-iterator/go"
)

var ResultPath = "/root/AItrade/CoinAI.net/task/analyConfig/data"

func AnalyBillingArr(InstID string) {
	fmt.Println("开始分析")
	BtcFilePath := mStr.Join(ResultPath, "/", InstID, "-BillingArr.json")
	file, _ := os.ReadFile(BtcFilePath)

	var BillingArr []testHunter.BillingType
	jsoniter.Unmarshal(file, &BillingArr)

	MoneySortArr := MoneySort(BillingArr)

	WinArr := []testHunter.BillingType{}
	for _, Billing := range MoneySortArr {
		// fmt.Println(Billing.Money)
		if mCount.Le(Billing.Money, "6600") > 0 {
			WinArr = append(WinArr, Billing)
		}
	}

	resultPath := mStr.Join(ResultPath, "/", InstID, "-WinArr.json")

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
BTC
EMA_86_CAP_2_level_1
EMA_84_CAP_2_level_1
EMA_80_CAP_2_level_1
EMA_88_CAP_4_level_1
EMA_82_CAP_5_level_1
EMA_88_CAP_5_level_1

//参数 范围  83  84  85 86   87
*/

/*
ETH
EMA_78_CAP_2_level_1
EMA_80_CAP_2_level_1
EMA_74_CAP_3_level_1
EMA_76_CAP_2_level_1

//参数 范围  77  78  79  80  81
*/
