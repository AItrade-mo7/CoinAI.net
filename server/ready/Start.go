package ready

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"CoinAI.net/server/global"
	"CoinAI.net/server/okxInfo"
	"CoinAI.net/server/tmpl"
	"github.com/EasyGolang/goTools/mClock"
	"github.com/EasyGolang/goTools/mOKX"
)

func Start() {
	SetMarket()
	go mClock.New(mClock.OptType{
		Func: SetMarket,
		Spec: "30 0,15,30,45 * * * ? ",
	})
}

func SetMarket() {
	err := CheckAccount()
	if err != nil {
		return
	}
	// 获取产品信息
	GetSWAPInst()

	// 获取市场行情
	GetCoinMarket()

	// 获取币种历史数据
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

	// 开始获取历史分析结果
	GetAnalyList()
}

// 用户信息检查
func CheckAccount() (resErr error) {
	GetUserInfo()
	GetOkxKey()

	resErr = nil
	if len(okxInfo.CoinServe.OkxKeyID) < 10 {
		resErr = fmt.Errorf("读取 dbData.CoinServe 失败 %+v", okxInfo.CoinServe)
		global.LogErr(resErr)
		return
	}

	if len(okxInfo.UserInfo.OkxKeyList) < 1 {
		resErr = fmt.Errorf("读取 dbData.UserInfo 失败 %+v", okxInfo.UserInfo)
		global.LogErr(resErr)
		return
	}

	for _, val := range okxInfo.UserInfo.OkxKeyList {
		if okxInfo.CoinServe.OkxKeyID == val.OkxKeyID {
			okxInfo.OkxKey = val
			break
		}
	}

	if len(okxInfo.OkxKey.OkxKeyID) < 10 {
		resErr = fmt.Errorf("读取 dbData.OkxKey 失败 %+v", okxInfo.OkxKey)
		global.LogErr(resErr)
		return
	}

	Body := new(bytes.Buffer)
	Tmpl := template.Must(template.New("").Parse(tmpl.StartSlice))
	Tmpl.Execute(Body, tmpl.StartSliceParam{
		CoinServeID: okxInfo.CoinServe.CoinServeID,
	})
	Message := Body.String()

	go global.Email(global.EmailOpt{
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
