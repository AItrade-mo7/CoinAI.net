package hunter

import (
	"CoinAI.net/server/global"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mTime"
)

func Start() {
	for ok := range okxInfo.Ticking {
		global.RunLog.Println("hunter.Start 执行", mTime.UnixFormat(mTime.GetUnixInt64()), ok)
		Running()
	}
}

func FileBaseKdata() {
	// List := mOKX.GetKdata(mOKX.GetKdataOpt{
	// 	InstID: item.InstID,
	// })
}

func Running() {
	// mFile.Write(config.Dir.JsonData+"/NowTicker.json", string(mJson.ToJson(okxInfo.Inst)))
}
