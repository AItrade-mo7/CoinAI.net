package ready

import (
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/okxInfo"
	"CoinAI.net/server/utils/dbUser"
	"github.com/EasyGolang/goTools/mClock"
)

func Start() {
	ReadUserInfo()
	SendStartEmail()

	GetAnalyData()
	go mClock.New(mClock.OptType{
		Func: GetAnalyData,
		Spec: "10 1,6,11,16,21,26,31,36,41,46,51,56 * * * ? ", // 每隔5分钟比标准时间晚一分钟 过 10 秒执行查询
	})
}

func GetAnalyData() {
	go ReadUserInfo()

	okxInfo.Inst = GetInstAll()

	okxInfo.NowTicker = GetNowTickerAnaly()

	okxInfo.Ticking <- "Tick"
}

func ReadUserInfo() {
	if len(config.AppEnv.UserID) > 10 {
		dbUser.NewUserDB(dbUser.NewUserOpt{
			UserID: config.AppEnv.UserID,
		})

		// if len(okxInfo.UserInfo.Email) > 3 {
		// 	EmailTo := []string{}
		// 	EmailTo = append(EmailTo, config.Email.Account)
		// 	EmailTo = append(EmailTo, okxInfo.UserInfo.Email)
		// 	config.Email.To = EmailTo
		// }
	}
}

func SendStartEmail() {
	// Message := mStr.Join(
	// 	"服务已启动: ", config.AppEnv.ServeID,
	// 	`<br /> <a href="https://trade.mo7.cc/CoinServe/CoinAI?id=`,
	// 	config.AppEnv.ServeID,
	// 	`"> https://trade.mo7.cc/CoinServe/CoinAI?id=`,
	// 	config.AppEnv.ServeID,
	// 	`</a> <br />`,
	// 	"用户昵称: ",
	// 	okxInfo.UserInfo.NickName,
	// 	"<br />",
	// 	"用户ID: ",
	// 	okxInfo.UserInfo.UserID,
	// 	"<br />",
	// )

	// global.Email(global.EmailOpt{
	// 	To:       config.Email.To,
	// 	Subject:  "系统提示",
	// 	Template: tmpl.SysEmail,
	// 	SendData: tmpl.SysParam{
	// 		Message: Message,
	// 		SysTime: mTime.UnixFormat(mTime.GetUnixInt64()),
	// 	},
	// }).Send()
}
