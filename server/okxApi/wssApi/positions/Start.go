package positions

import (
	"strings"
	"time"

	"CoinAI.net/server/global"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mStr"
)

// 持仓频道

func Start() {
	global.WssLog.Println("=============== positions.Start ================")
	Wss := mOKX.WssOKX(mOKX.OptWssOKX{
		FetchType: 1,
		Event: func(lType string, a any) {
			cont := mStr.ToStr(a)
			if cont == "pong" || cont == "ping" {
				// global.WssLog.Println("positions.Start1", lType, cont)
				return
			}

			if lType == "Write" {
				loginInfo, err := Write_LoginInfo(a)
				global.WssLog.Println("positions.Start2", lType, err, mJson.ToStr(loginInfo))
				return
			}

			if lType == "Close" {
				global.WssLog.Println("positions.Start4", lType, cont)
				time.Sleep(time.Second * 2)
				Start()
				return
			}
		},
		// OKXKey: mOKX.TypeOkxKey{
		// 	ApiKey:     config.AppEnv.ApiKey,
		// 	SecretKey:  config.AppEnv.SecretKey,
		// 	Passphrase: config.AppEnv.Passphrase,
		// },
	})

	Wss.Read(func(msg []byte) {
		global.WssLog.Println("positions.Start5", mStr.ToStr(msg))

		cont := mStr.ToStr(msg)
		find := strings.Contains(cont, "login") // 登录是否成功
		if find {
			isLogin := Read_Login(msg)
			if !isLogin {
				Wss.Close("登录失败")
				return
			}
		}
	})
}
