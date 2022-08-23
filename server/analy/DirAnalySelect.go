package analy

import (
	"fmt"

	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mJson"
)

func DirAnalySelect() {
	fmt.Println(okxInfo.WholeDir)
	mJson.Println(okxInfo.AnalySingle)
}
