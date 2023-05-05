package ready

import (
	"time"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/global/dbType"
	"CoinAI.net/server/hunter"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mClock"
	"github.com/EasyGolang/goTools/mCount"
	"github.com/EasyGolang/goTools/mFile"
	"github.com/EasyGolang/goTools/mJson"
)

var HunterOptArr = []hunter.HunterOpt{
	{
		HunterName: "BTC-CAP",
		InstID:     "BTC-USDT",
		Describe:   "以 BTC-USDT 交易对为主执行自动交易,支持的资金量更大,更加稳定",
		TradeKdataOpt: dbType.TradeKdataOpt{
			EMA_Period:    294, // 参数已确定  2023-04-11 18:14
			CAP_Period:    6,
			CAP_Max:       "3",
			CAP_Min:       "-0.5",
			MaxTradeLever: 4,
		},
	},
	{
		HunterName: "ETH-CAP",
		InstID:     "ETH-USDT",
		Describe:   "以 ETH-USDT 交易对为主执行自动交易,交易次数更加频发,可以收获更高收益",
		TradeKdataOpt: dbType.TradeKdataOpt{
			EMA_Period:    80, // 参数确定时间 2023-4-11 20:28:37
			CAP_Period:    2,
			CAP_Max:       "0.5",
			CAP_Min:       "-0.5",
			MaxTradeLever: 3,
		},
	},
}

func Start() {
	// 初始化 Hunter 全局 初始值
	okxInfo.OkxInfoInit()

	// 发送启动邮件
	StartEmail()

	// 数据预填充
	GetAnalyData()
	go mClock.New(mClock.OptType{
		Func: func() {
			RoundNum := mCount.GetRound(0, 60) // 构建请求延迟
			time.Sleep(time.Second * time.Duration(RoundNum))
			GetAnalyData()
		},
		Spec: "10 1,6,11,16,21,26,31,36,41,46,51,56 * * * ? ", // 每隔5分钟比标准时间晚一分钟 过 10 秒执行查询
	})

	for _, conf := range HunterOptArr {
		hunter.New(conf).Start()
	}

	CheckOKXAccount()
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

func CheckOKXAccount() {
	for _, conf := range HunterOptArr {
		HunterName := conf.HunterName
		MaxTradeLever := conf.TradeKdataOpt.MaxTradeLever

		for key, ApiKey := range config.AppEnv.ApiKeyList {
			if ApiKey.Hunter == HunterName {
				// 检查杠杆倍率 并重制
				if ApiKey.TradeLever < 1 {
					config.AppEnv.ApiKeyList[key].TradeLever = 1
				}
				if ApiKey.TradeLever > MaxTradeLever {
					config.AppEnv.ApiKeyList[key].TradeLever = MaxTradeLever
				}
			}
		}
	}

	global.WriteAppEnv()
}
