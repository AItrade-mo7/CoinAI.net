package ready

import (
	"CoinAI.net/server/global"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mClock"
)

func Start() {
	//
	GetAnalyData()
	go mClock.New(mClock.OptType{
		Func: GetAnalyData,
		Spec: "30 0,5,10,15,20,25,30,35,40,45,50,55 * * * ? ", // 5 分的整数过 30 秒
	})
}

/*
思路：5分整数倍过30秒拉取一波行情和交易信息

*/

func GetAnalyData() {
	okxInfo.NowTicker = GetNowTickerAnaly()

	// for Key, SingleList := range okxInfo.NowTicker.AnalySingle {
	// 	fmt.Println(Key)
	// 	for _, Single := range SingleList {
	// 		fmt.Println(Single)
	// 	}
	// }

	global.RunLog.Println("拉取一次数据接口")
}
