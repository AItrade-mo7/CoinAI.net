package positions

import (
	"CoinAI.net/server/global"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mStr"
)

// 持仓频道

func Start() {
	wss := mOKX.WssOKX(mOKX.OptWssOKX{
		FetchType: 1,
		Event: func(s string, a any) {
			global.WssLog.Println("positions.Start", s, mStr.ToStr(a))
		},
		OKXKey: mOKX.TypeOkxKey{
			ApiKey:     okxInfo.OkxKey.ApiKey,
			SecretKey:  okxInfo.OkxKey.SecretKey,
			Passphrase: okxInfo.OkxKey.Passphrase,
		},
	})

	wss.Read(func(msg []byte) {
		global.WssLog.Println("读数据", mStr.ToStr(msg))
	})
}
