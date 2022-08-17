package analy

import (
	"fmt"

	"github.com/EasyGolang/goTools/mOKX"
)

func SelectInst(AnalyKdata map[string][]mOKX.TypeKd) {
	for key, list := range AnalyKdata {
		fmt.Println(key, len(list))
	}
}
