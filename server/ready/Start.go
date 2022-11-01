package ready

import (
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
		Spec: "30 0,5,10,15,20,25,30,35,40,45,50,55 * * * ? ", // 5 分的整数过 30 秒
	})
}

func GetAnalyData() {
	go ReadUserInfo()

	okxInfo.InstAll = GetInstAll()

	okxInfo.NowTicker = GetNowTickerAnaly()

	if len(okxInfo.NowTicker.TickerVol) > 0 {
		okxInfo.TradeInst = okxInfo.NowTicker.TickerVol[1]
	}

	global.RunLog.Println("拉取一次数据接口")
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
		"服务已启动,",
		`<br /> <a href="https://trade.mo7.cc/CoinServe"> https://trade.mo7.cc/CoinServe </a> <br />`,
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
			SysTime: mTime.IsoTime(),
		},
	}).Send()
}
