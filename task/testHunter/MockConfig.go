package testHunter

import (
	"CoinAI.net/server/hunter"
	"github.com/EasyGolang/goTools/mStr"
)

type MockOptType struct {
	MockOpt       BillingType
	TradeKdataOpt hunter.TradeKdataOpt
}

func MockConfig(EmaPArr []int) []MockOptType {
	MockConfigArr := []MockOptType{}

	CAP := 4 //  3 或者 4

	for _, emaP := range EmaPArr {
		MockConfigArr = append(MockConfigArr,
			MockOptType{
				BillingType{
					MockName:  "MA_" + mStr.ToStr(emaP) + "_CAP_" + mStr.ToStr(CAP),
					InitMoney: "1000", // 初始资金
					Level:     "1",    // 杠杆倍数
					Charge:    "0.05", // 吃单标准手续费率 0.05%
				},
				hunter.TradeKdataOpt{
					MA_Period:      emaP,
					RSI_Period:     18,
					RSI_EMA_Period: 14,
					CAP_Period:     CAP,
				},
			},
		)
	}

	return MockConfigArr
}
