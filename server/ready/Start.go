package ready

import (
	"fmt"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/okxInfo"
	"CoinAI.net/server/tmpl"
	"CoinAI.net/server/utils/dbUser"
	"github.com/EasyGolang/goTools/mClock"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/EasyGolang/goTools/mTime"
)

func Start() {
	ReadUserInfo()
	SendStartEmail()

	GetAnalyData()
	go mClock.New(mClock.OptType{
		Func: GetAnalyData,
		Spec: "20 1,6,11,16,21,26,31,36,41,46,51,56 * * * ? ", // 每隔5分钟比标准时间晚一分钟 过 20 秒执行查询
	})
}

func GetAnalyData() {
	go ReadUserInfo()

	okxInfo.Inst = GetInstAll()

	okxInfo.NowTicker = GetNowTickerAnaly()

	err := SetTradeInst()
	if err != nil {
		global.LogErr(err)
		return
	}

	okxInfo.Ticking <- "Tick"
}

func SetTradeInst() (resErr error) {
	resErr = nil
	if len(okxInfo.NowTicker.TickerVol) < 4 {
		resErr = fmt.Errorf("榜单数量不足 ready.SetTradeInst TickerVol %+v", len(okxInfo.NowTicker.TickerVol))
		return
	}
	if len(okxInfo.NowTicker.AnalyWhole) < 3 {
		resErr = fmt.Errorf("切片数量不足 ready.SetTradeInst AnalyWhole %+v", len(okxInfo.NowTicker.AnalyWhole))
		return
	}
	if len(okxInfo.NowTicker.AnalySingle) < 4 {
		resErr = fmt.Errorf("币种分析样本不足 ready.SetTradeInst AnalySingle %+v", len(okxInfo.NowTicker.AnalySingle))
		return
	}

	if len(okxInfo.NowTicker.MillionCoin) < 2 {
		resErr = fmt.Errorf("数据异常 ready.SetTradeInst MillionCoin %+v", len(okxInfo.NowTicker.MillionCoin))
		return
	}
	// 在这里按照涨幅的绝对值排个序SetTradeInst

	HLPerList := Sort_HLPer(okxInfo.NowTicker.MillionCoin)

	CoinId := HLPerList[len(HLPerList)-1].InstID

	for _, val := range HLPerList {
		if val.InstID == "ETH-USDT" || val.InstID == "BTC-USDT" {
			CoinId = val.InstID
			break
		}
	}

	if len(CoinId) < 0 {
		resErr = fmt.Errorf("数据异常 ready.SetTradeInst CoinId %+v", CoinId)
		return
	}

	KdataInst := okxInfo.Inst[CoinId]
	TradeInst := okxInfo.Inst[CoinId+"-SWAP"]

	if KdataInst.State == "live" && TradeInst.State == "live" {
	} else {
		resErr = fmt.Errorf("数据异常 ready.SetTradeInst State %+v", KdataInst.State)
		return
	}

	if KdataInst.InstType == "SPOT" && TradeInst.InstType == "SWAP" {
	} else {
		resErr = fmt.Errorf("数据异常 ready.SetTradeInst InstType %+v", KdataInst.InstType)
		return
	}

	okxInfo.KdataInst = KdataInst
	okxInfo.TradeInst = TradeInst

	return
}

func ReadUserInfo() {
	if len(config.AppEnv.UserID) > 10 {
		dbUser.NewUserDB(dbUser.NewUserOpt{
			UserID: config.AppEnv.UserID,
		})

		if len(okxInfo.UserInfo.Email) > 3 {
			EmailTo := []string{}
			EmailTo = append(EmailTo, config.Email.Account)
			EmailTo = append(EmailTo, okxInfo.UserInfo.Email)
			config.Email.To = EmailTo
		}
	}
}

func SendStartEmail() {
	Message := mStr.Join(
		"服务已启动: ", config.AppEnv.ServeID,
		`<br /> <a href="https://trade.mo7.cc/CoinServe/CoinAI?id=`,
		config.AppEnv.ServeID,
		`"> https://trade.mo7.cc/CoinServe/CoinAI?id=`,
		config.AppEnv.ServeID,
		`</a> <br />`,
		"用户昵称: ",
		okxInfo.UserInfo.NickName,
		"<br />",
		"用户ID: ",
		okxInfo.UserInfo.UserID,
		"<br />",
	)

	global.Email(global.EmailOpt{
		To:       config.Email.To,
		Subject:  "系统提示",
		Template: tmpl.SysEmail,
		SendData: tmpl.SysParam{
			Message: Message,
			SysTime: mTime.UnixFormat(mTime.GetUnixInt64()),
		},
	}).Send()
}
