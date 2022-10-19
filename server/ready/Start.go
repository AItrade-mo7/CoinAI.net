package ready

import (
	"CoinAI.net/server/global"
	"github.com/EasyGolang/goTools/mClock"
)

func Start() {
	//
	go mClock.New(mClock.OptType{
		Func: GetAnalyData,
		Spec: "30 0,5,10,15,20,25,30,35,40,45,50,55 * * * ? ", // 5 分的整数过 30 秒
	})
}

/*
思路：5分整数倍过30秒拉取一波行情和交易信息

*/

func GetAnalyData() {
	GetSWAPInst() // 获取所有产品的合约信息
	global.RunLog.Println("拉取一次 Analy 接口")
}
