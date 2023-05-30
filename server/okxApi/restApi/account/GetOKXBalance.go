package account

import (
	"fmt"
	"go.uber.org/zap"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/global/dbType"
	"github.com/EasyGolang/goTools/mCount"
	"github.com/EasyGolang/goTools/mFile"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/EasyGolang/goTools/mTime"
	jsoniter "github.com/json-iterator/go"
)

type BalanceReq []struct {
	Details []BalanceDetails
	UTime   string
}
type BalanceDetails struct {
	UTime string
	Ccy   string // 币种
	DisEq string // 美金层面折算
}

// 查看账户余额
func GetOKXBalance(OKXKey dbType.OkxKeyType) (resData []dbType.AccountBalance, resErr error) {
	resData = []dbType.AccountBalance{}
	resErr = nil

	if len(OKXKey.ApiKey) < 10 {
		resErr = fmt.Errorf("account.GetOKXBalance OKXKey.ApiKey 不能为空 %+v Name:%+v", OKXKey.ApiKey, OKXKey.Name)
		global.LogErr("该错误已同步至用户邮箱", resErr)
		LogErr(OKXKey, resErr)
		return
	}

	res, err := mOKX.FetchOKX(mOKX.OptFetchOKX{
		Path:   "/api/v5/account/balance",
		Method: "GET",
		OKXKey: mOKX.TypeOkxKey{
			ApiKey:     OKXKey.ApiKey,
			SecretKey:  OKXKey.SecretKey,
			Passphrase: OKXKey.Passphrase,
		},
	})
	// 打印接口日志
	global.OKXLogo.Info("account.GetOKXBalance",
		zap.Error(err),
		zap.String("res", mStr.ToStr(res)),
		zap.String("name", mJson.ToStr(OKXKey.Name)),
	)

	if err != nil {
		resErr = fmt.Errorf("account.GetOKXBalance1 %s Name:%+v", mStr.ToStr(err), OKXKey.Name)
		global.LogErr("该错误已同步至用户邮箱", resErr)
		LogErr(OKXKey, resErr)
		return
	}

	var resObj mOKX.TypeReq
	jsoniter.Unmarshal(res, &resObj)

	if resObj.Code != "0" {
		resErr = fmt.Errorf("account.GetOKXBalance2 %s Name:%+v", mStr.ToStr(res), OKXKey.Name)
		global.LogErr("该错误已同步至用户邮箱", resErr)
		LogErr(OKXKey, resErr)
		return
	}
	var Data BalanceReq
	jsoniter.Unmarshal(mJson.ToJson(resObj.Data), &Data)

	if len(Data) > 0 {
		Details := Data[0].Details
		for _, val := range Details {
			myTime := mTime.MsToTime(val.UTime, "0")
			NewBalance := dbType.AccountBalance{
				TimeUnix: mTime.ToUnixMsec(myTime),
				TimeStr:  mTime.UnixFormat(val.UTime),
				CcyName:  val.Ccy,
				Balance:  mCount.Cent(val.DisEq, 2),
			}
			if mCount.Le(NewBalance.Balance, "0.1") > -1 {
				resData = append(resData, NewBalance)
			}
		}
	}

	Balance_file := mStr.Join(config.Dir.JsonData, "/Balance.json")
	mFile.Write(Balance_file, string(res))

	return
}
