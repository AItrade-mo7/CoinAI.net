package analyConfig

import (
	"os"

	"CoinAI.net/server/global"
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

	AfterN := 40

	// Money最高来排序
	MoneyArr := MoneySort(BillingArr)
	for _, item := range MoneyArr {
		if mCount.Le(item.ResultMoney, "1000") > 0 {
			Tmp := `结算最高:
参数名称: ${MockName}
InstID: ${InstID}
开仓频率: ${OrderRate} 
胜率: ${WinRatio}
盈亏比: ${PLratio}
最终金钱: ${ResultMoney}
平仓后历史最低余额: ${MinMoney}
持仓过程中最低盈利比率: ${PositionMinRatio}
杠杆倍率: ${Level}
总手续费: ${ChargeAdd}
`
			Data := map[string]string{
				"MockName":         item.MockName,
				"InstID":           item.InstID,
				"OrderRate":        item.OrderRate,
				"WinRatio":         item.WinRatio,
				"PLratio":          item.PLratio,
				"ResultMoney":      item.ResultMoney,
				"MinMoney":         mJson.ToStr(item.MinMoney),
				"PositionMinRatio": mJson.ToStr(item.PositionMinRatio),
				"Level":            item.Level,
				"ChargeAdd":        item.ChargeAdd,
			}
			global.Run.Println(mStr.Temp(Tmp, Data))
		}
	}

	MoneyArr = MoneyArr[0:AfterN] // 取前 20 位
	resultPath := mStr.Join(opt.ResultBasePath, "/", opt.InstID, "-MoneyArr.json")
	mFile.Write(resultPath, mJson.ToStr(MoneyArr))

	// 胜率 最高来排序
	WinArr := MoneySort(BillingArr)
	for _, item := range WinArr {
		if mCount.Le(item.ResultMoney, "1000") > 0 {
			Tmp := `胜率最高:
参数名称: ${MockName}
InstID: ${InstID}
开仓频率: ${OrderRate} 
胜率: ${WinRatio}
盈亏比: ${PLratio}
最终金钱: ${ResultMoney}
平仓后历史最低余额: ${MinMoney}
持仓过程中最低盈利比率: ${PositionMinRatio}
杠杆倍率: ${Level}
总手续费: ${ChargeAdd}
`
			Data := map[string]string{
				"MockName":         item.MockName,
				"InstID":           item.InstID,
				"OrderRate":        item.OrderRate,
				"WinRatio":         item.WinRatio,
				"PLratio":          item.PLratio,
				"ResultMoney":      item.ResultMoney,
				"MinMoney":         mJson.ToStr(item.MinMoney),
				"PositionMinRatio": mJson.ToStr(item.PositionMinRatio),
				"Level":            item.Level,
				"ChargeAdd":        item.ChargeAdd,
			}
			global.TradeLog.Println(mStr.Temp(Tmp, Data))
		}
	}
	WinArr = WinArr[0:AfterN]
	resultPath = mStr.Join(opt.ResultBasePath, "/", opt.InstID, "-WinArr.json")
	mFile.Write(resultPath, mJson.ToStr(WinArr))
}

// Money 排序
func MoneySort(arr []testHunter.BillingType) []testHunter.BillingType {
	size := len(arr)
	list := make([]testHunter.BillingType, size)
	copy(list, arr)

	var swapped bool
	for i := size - 1; i > 0; i-- {
		swapped = false
		for j := 0; j < i; j++ {
			a := list[j+1].ResultMoney
			b := list[j].ResultMoney
			if mCount.Le(a, b) < 0 {
				list[j], list[j+1] = list[j+1], list[j]
				swapped = true
			}
		}
		if !swapped {
			break
		}
	}

	rList := make(
		[]testHunter.BillingType,
		len(list),
		len(list)*2,
	)
	j := 0
	for i := len(list) - 1; i > -1; i-- {
		rList[j] = list[i]
		j++
	}
	return rList
}

// Win 排序
func WinSort(arr []testHunter.BillingType) []testHunter.BillingType {
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

	rList := make(
		[]testHunter.BillingType,
		len(list),
		len(list)*2,
	)
	j := 0
	for i := len(list) - 1; i > -1; i-- {
		rList[j] = list[i]
		j++
	}
	return rList
}
