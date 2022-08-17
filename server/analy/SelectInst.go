package analy

import (
	"fmt"

	"github.com/EasyGolang/goTools/mOKX"
)

func SelectInst(AnalyKdata map[string][]mOKX.TypeKd) {
	for key, list := range AnalyKdata {

		// text
		if key != "EOS-USDT" {
			continue
		}
		// text

		fmt.Println(key)
		if len(list) == 300 {
			SingleAnalyInst(list)
		}
	}
}

// 在这里判断趋势
func SingleAnalyInst(list []mOKX.TypeKd) {
	for _, item := range list {
		fmt.Println(item.Time, item.C, item.O)
	}
}
