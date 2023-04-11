package analyConfig

import (
	"fmt"
	"os"

	"CoinAI.net/task/testHunter"
	"github.com/EasyGolang/goTools/mPath"
	"github.com/EasyGolang/goTools/mStr"
	jsoniter "github.com/json-iterator/go"
)

type GetWinConfigOpt struct {
	OutPutDir string
	InstID    string
}

func GetWinConfig(opt GetWinConfigOpt) {
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

	for _, item := range BillingArr {
		fmt.Println(item.MockName)
	}
}
