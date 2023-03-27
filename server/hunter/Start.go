package hunter

import (
	"time"

	"CoinAI.net/server/global"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mTime"
)

func Start() {
	for ok := range okxInfo.HunterTicking {
		global.TradeLog.Println(" ===== hunter.Start 执行 ===== ", mTime.UnixFormat(mTime.GetUnixInt64()), ok)
		Running()
	}
}

func Running() {
	global.TradeLog.Println(" === hunter.Running === ", okxInfo.KdataInst.InstID)

	err := SetTradeInst() // 如果均为空则设置一下
	if err != nil {
		global.LogErr(err)
		return
	}

	// if len(okxInfo.KdataInst.InstID) < 2 || len(okxInfo.TradeInst.InstID) < 2 {
	// 	global.LogErr("hunter.Running", "okxInfo.TradeInst.InstID 为空")
	// 	return
	// }

	// FileBaseKdata()

	// SetNowKdata()

	// FormatTradeKdata()

	// Analy()
}

func FileBaseKdata() {
	Page := 2 // 如果数组为空，则填充 300 条进去
	if len(okxInfo.NowKdataList) < 100 {
		// 回填历史数据 1 组
		for i := Page; i >= 0; i-- {
			time.Sleep(time.Second / 3)
			List := mOKX.GetKdata(mOKX.GetKdataOpt{
				InstID: okxInfo.KdataInst.InstID,
				Page:   i,
				After:  mTime.GetUnixInt64(),
			})

			okxInfo.NowKdataList = append(okxInfo.NowKdataList, List...)

			// for _, val := range List {
			// 	okxInfo.NowKdataList = append(okxInfo.NowKdataList, val)
			// }
		}
		// global.RunLog.Println("历史数据回填完毕", len(okxInfo.NowKdataList))
	} else { // 如果不为空 则检查当前的数组和持仓币种的关系
		// 在这里执行重启
		if okxInfo.KdataInst.InstID != okxInfo.NowKdataList[len(okxInfo.NowKdataList)-1].InstID {
			okxInfo.NowKdataList = []mOKX.TypeKd{} // 清空历史数据
			Running()                              // 立即重新执行一次 Running
			warnStr := "切换监听币种为: " + okxInfo.KdataInst.InstID
			// global.RunLog.Println(warnStr)
			go SendEmail(warnStr)
		}
	}
}

func SendEmail(Message string) {
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
