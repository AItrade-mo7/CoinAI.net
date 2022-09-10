package positions

import (
	"strings"

	"CoinAI.net/server/global"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mStr"
)

// 持仓频道

func Start() {
	global.WssLog.Println("=============== positions.Start ================")
	wss := mOKX.WssOKX(mOKX.OptWssOKX{
		FetchType: 1,
		Event: func(lType string, a any) {
			cont := mStr.ToStr(a)

			if cont == "pong" || cont == "ping" {
				return
			}

			if lType == "Write" {
				loginInfo, err := Write_LoginInfo(a)
				global.WssLog.Println("positions.Start", lType, err, mStr.ToStr(loginInfo), cont)
				return
			}

			if lType != "Read" {
				global.WssLog.Println("positions.Start", lType, cont)
				return
			}

			if lType != "Close" {
				global.WssLog.Println("positions.Start", lType, "链接关闭，重启", cont)
				Start()
				return
			}
		},
		OKXKey: mOKX.TypeOkxKey{
			ApiKey:     okxInfo.OkxKey.ApiKey,
			SecretKey:  okxInfo.OkxKey.SecretKey,
			Passphrase: okxInfo.OkxKey.Passphrase,
		},
	})

	wss.Read(func(msg []byte) {
		global.WssLog.Println("positions.Start", "wss.Read", mStr.ToStr(msg))

		cont := mStr.ToStr(msg)
		find := strings.Contains(cont, "login") // 是否为SWAP
		if find {
			isLogin := Read_Login(msg)
			if !isLogin {
				wss.Close("登录失败")
				return
			}
		}
	})
}
