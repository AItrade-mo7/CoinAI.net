package ready

import (
	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mClock"
	"github.com/EasyGolang/goTools/mFile"
	"github.com/EasyGolang/goTools/mJson"
)

func Start() {
	StartEmail()

	GetAnalyData()
	go mClock.New(mClock.OptType{
		Func: GetAnalyData,
		Spec: "10 1,6,11,16,21,26,31,36,41,46,51,56 * * * ? ", // 每隔5分钟比标准时间晚一分钟 过 10 秒执行查询
	})
}

func GetAnalyData() {
	go global.GetMainUser()

	okxInfo.Inst = GetInstAll()

	mFile.Write(config.Dir.JsonData+"/InstAll.json", mJson.ToStr(okxInfo.Inst))

	okxInfo.NowTicker = GetNowTickerAnaly()

	mFile.Write(config.Dir.JsonData+"/NowTicker.json", mJson.ToStr(okxInfo.NowTicker))

	// 挑选交易币种
	// 在这里先写死
	if len(okxInfo.NowTicker.TickerVol) < 1 {
		global.LogErr("ready.GetAnalyData  okxInfo.NowTicker.TickerVol 长度不足", len(okxInfo.NowTicker.TickerVol))
		return
	}

	TradeInstID := okxInfo.NowTicker.TickerVol[0].InstID
	// 设置当前交易品信息并获取它的K线
	okxInfo.NowKdataList = GetNowKdata(TradeInstID)

	SetTradeInst()

	okxInfo.Ticking <- "Tick"
}

func SetTradeInst() {
	// 这里 一定为 合约交易对
	InstID := okxInfo.NowTicker.TickerVol[0].InstID
	okxInfo.TradeInst = okxInfo.Inst[InstID+"-SWAP"]
}
