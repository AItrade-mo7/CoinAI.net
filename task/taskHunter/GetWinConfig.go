package taskHunter

import (
	"fmt"
	"os"

	"CoinAI.net/server/global/dbType"
	"CoinAI.net/server/hunter/testHunter"
	"github.com/EasyGolang/goTools/mPath"
	"github.com/EasyGolang/goTools/mStr"
	jsoniter "github.com/json-iterator/go"
)

type GetWinConfigOpt struct {
	OutPutDir string
	InstID    string
}

func GetWinConfig(opt GetWinConfigOpt) []dbType.TradeKdataOpt {
	WinArrDataPath := mStr.Join(opt.OutPutDir, "/", opt.InstID, "-WinArr.json")
	if !mPath.Exists(WinArrDataPath) {
		err := fmt.Errorf("文件不存在 %+v", opt.OutPutDir)
		panic(err)
	}
	file, err := os.ReadFile(WinArrDataPath)
	if err != nil {
		err := fmt.Errorf("读取文件出错 %+v", err)
		panic(err)
	}
	var BillingArr []testHunter.BillingType // 数据来源
	jsoniter.Unmarshal(file, &BillingArr)

	confArr := []dbType.TradeKdataOpt{}
	for _, item := range BillingArr {
		// MockName := item.MockName //  EMA_320_CAP_5_CAPMax_1_level_1

		// EMA := regText(MockName, []string{
		// 	"EMA_", "_CAP_",
		// })
		// CAP := regText(MockName, []string{
		// 	"_CAP_", "_CAPMax_",
		// })
		// CAPMax := regText(MockName, []string{
		// 	"_CAPMax_", "_level_",
		// })

		// conf := dbType.TradeKdataOpt{
		// 	EMA_Period: mCount.ToInt(EMA),
		// 	CAP_Period: mCount.ToInt(CAP),
		// 	CAP_Max:    CAPMax, // CAP 判断的边界值 0.2
		// }
		confArr = append(confArr, item.Opt)
	}

	return confArr
}

// func regText(origin string, reg []string) string {
// 	comp := mStr.Join(reg[0], "(.*?)", reg[1])
// 	flysnowRegexp := regexp.MustCompile(comp)
// 	params := flysnowRegexp.FindStringSubmatch(origin)
// 	if len(params) == 2 {
// 		return params[1]
// 	}
// 	if len(params) == 1 {
// 		return params[0]
// 	}
// 	return ""
// }
