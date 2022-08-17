package analy

import (
	"fmt"

	"github.com/EasyGolang/goTools/mOKX"
)

func SelectInst(AnalyKdata map[string][]mOKX.TypeKd) {
	for key, list := range AnalyKdata {
		fmt.Println(key)
		if len(list) == 300 {
			SingleAnalyInst(list)
		}
	}
}

func SingleAnalyInst(list []mOKX.TypeKd) {
	for _, item := range list {
		fmt.Println(item.Time)
	}
}
