package analy

import (
	"fmt"

	"CoinAI.net/server/okxInfo"
)

func SelectInst() {
	for key, list := range okxInfo.AnalyKdata {
		fmt.Println(key, len(list))
	}
}
