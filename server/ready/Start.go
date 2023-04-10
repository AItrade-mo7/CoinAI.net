package ready

import (
	"time"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/hunter"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mClock"
	"github.com/EasyGolang/goTools/mCount"
	"github.com/EasyGolang/goTools/mFile"
	"github.com/EasyGolang/goTools/mJson"
)

func Start() {
	// 初始化数据初始值
	okxInfo.OkxInfoInit()
	// 发送启动邮件
	StartEmail()
	// 数据预填充

	GetAnalyData()

	// 准备策略
	// 策略 1
	BTCHunter := hunter.New(hunter.HunterOpt{
		HunterName: "BTC-CoinAI",
		InstID:     "BTC-USDT",
		Describe:   "以 BTC-USDT 交易对为主执行自动交易,支持的资金量更大,更加稳定",
	})

	BTCHunter.Start()

	// 策略 2
	ETHHunter := hunter.New(hunter.HunterOpt{
		HunterName: "ETH-CoinAI",
		InstID:     "ETH-USDT",
		Describe:   "以 ETH-USDT 交易对为主执行自动交易,交易次数更加频发,可以收获更高收益",
	})

	ETHHunter.Start()

	// 构建定时任务
	go mClock.New(mClock.OptType{
		Func: func() {
			RoundNum := mCount.GetRound(0, 60) // 构建请求延迟
			time.Sleep(time.Second * time.Duration(RoundNum))
			GetAnalyData()
			BTCHunter.Start() // BTC 策略
			ETHHunter.Start() // ETH 策略
		},
		Spec: "10 1,6,11,16,21,26,31,36,41,46,51,56 * * * ? ", // 每隔5分钟比标准时间晚一分钟 过 10 秒执行查询
	})
}

func GetAnalyData() {
	go global.GetMainUser()

	okxInfo.Inst = GetInstAll()

	mFile.Write(config.Dir.JsonData+"/InstAll.json", mJson.ToStr(okxInfo.Inst))

	okxInfo.NowTicker = GetNowTickerAnaly()

	mFile.Write(config.Dir.JsonData+"/NowTicker.json", mJson.ToStr(okxInfo.NowTicker))

	// 在这里检查数据

	if len(okxInfo.Inst) < 10 {
		global.LogErr("ready.GetAnalyData okxInfo.Inst 长度不足", len(okxInfo.Inst))
		return
	}

	if len(okxInfo.NowTicker.TickerVol) < 3 {
		global.LogErr("ready.GetAnalyData okxInfo.NowTicker.TickerVol 长度不足", len(okxInfo.NowTicker.TickerVol))
		return
	}
}
