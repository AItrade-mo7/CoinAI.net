package ready

import (
	"bytes"
	"text/template"
	"time"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/okxInfo"
	"CoinAI.net/server/tmpl"
	"github.com/EasyGolang/goTools/mClock"
	"github.com/EasyGolang/goTools/mOKX"
)

func Start() {
	SetMarket()
	go mClock.New(mClock.OptType{
		Func: SetMarket,
		Spec: "40 0,15,30,45 * * * ? ",
	})
}

func SetMarket() {
	global.RunLog.Println("============= 开始执行周期任务 ==============")

	err := CheckAccount()
	if err != nil {
		return
	}
	global.RunLog.Println("获取 SWAP 品信息")
	GetSWAPInst()

	global.RunLog.Println("获取市场行情")
	GetCoinMarket()

	global.RunLog.Println("获取币种历史数据")
	okxInfo.AnalyKdata_SPOT = make(map[string][]mOKX.TypeKd)
	okxInfo.AnalyKdata_SWAP = make(map[string][]mOKX.TypeKd)
	AnalyKdata_SPOT := make(map[string][]mOKX.TypeKd)
	AnalyKdata_SWAP := make(map[string][]mOKX.TypeKd)
	for _, item := range okxInfo.MarketTicker.List {
		// 开始设置 SWAP
		SwapInst := mOKX.TypeInst{}
		for _, SWAP := range okxInfo.SWAP_inst {
			if SWAP.Uly == item.InstID {
				SwapInst = SWAP
				break
			}
		}
		if len(SwapInst.InstID) < 3 {
			continue
		}

		SPOT_list := GetCoinAnalyKdata(item.InstID)
		SWAP_list := GetCoinAnalyKdata(SwapInst.InstID)

		if len(SPOT_list) == 300 {
			AnalyKdata_SPOT[SwapInst.InstID] = SPOT_list
		}
		if len(SWAP_list) == 300 {
			AnalyKdata_SWAP[SwapInst.InstID] = SWAP_list
		}
	}

	okxInfo.AnalyKdata_SPOT = AnalyKdata_SPOT
	okxInfo.AnalyKdata_SWAP = AnalyKdata_SWAP

	global.RunLog.Println("开始获取历史分析结果列表")
	GetAnalyList()
}

// 用户信息检查
func CheckAccount() (resErr error) {
	global.RunLog.Println("开始获取用户数据")
	GetUserInfo()

	global.RunLog.Println("开始获取 OkxKey 数据")

	resErr = nil

	global.RunLog.Println("发送 启动邮件 邮件")

	Body := new(bytes.Buffer)
	Tmpl := template.Must(template.New("").Parse(tmpl.StartSlice))
	Tmpl.Execute(Body, tmpl.StartSliceParam{
		CoinServeID: config.AppEnv.IP + ":" + config.AppEnv.Port,
	})
	Message := Body.String()

	global.Email(global.EmailOpt{
		To: []string{
			okxInfo.UserInfo.Email,
		},
		Subject:  "CoinServe 启动成功",
		Template: tmpl.SysEmail,
		SendData: tmpl.SysParam{
			NickName:     okxInfo.UserInfo.NickName,
			Message:      Message,
			SysTime:      time.Now(),
			SecurityCode: okxInfo.UserInfo.SecurityCode,
		},
	}).Send()

	return
}
