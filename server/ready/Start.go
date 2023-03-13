package ready

import (
	"CoinAI.net/server/global/config"
	"github.com/EasyGolang/goTools/mJson"
)

func Start() {
	StartEmail()

	mJson.Println(config.AppEnv)
	mJson.Println(config.MainUser)
	mJson.Println(config.NoticeEmail)

	// GetAnalyData()
	// go mClock.New(mClock.OptType{
	// 	Func: GetAnalyData,
	// 	Spec: "10 1,6,11,16,21,26,31,36,41,46,51,56 * * * ? ", // 每隔5分钟比标准时间晚一分钟 过 10 秒执行查询
	// })
}

func GetAnalyData() {
	// go ReadUserInfo()

	// okxInfo.Inst = GetInstAll()

	// okxInfo.NowTicker = GetNowTickerAnaly()

	// okxInfo.Ticking <- "Tick"
}
