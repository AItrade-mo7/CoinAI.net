package account

import (
	"CoinAI.net/server/global/config"
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

type AccountBalance struct {
	TimeUnix int64  `bson:"TimeUnix"`
	TimeStr  string `bson:"TimeStr"`
	CcyName  string `bson:"CcyName"` // 币种
	Balance  string `bson:"Balance"` // 币种
}

func GetOKXBalance(OKXKey mOKX.TypeOkxKey) (resData []AccountBalance, resErr error) {
	resData = []AccountBalance{}
	resErr = nil

	res, err := mOKX.FetchOKX(mOKX.OptFetchOKX{
		Path:   "/api/v5/account/balance",
		Method: "GET",
		OKXKey: OKXKey,
	})
	if err != nil {
		resErr = err
		return
	}

	var resObj mOKX.TypeReq
	jsoniter.Unmarshal(res, &resObj)

	if resObj.Code != "0" {
		resErr = err
		return
	}
	var Data BalanceReq
	jsoniter.Unmarshal(mJson.ToJson(resObj.Data), &Data)

	if len(Data) > 0 {
		Details := Data[0].Details
		for _, val := range Details {
			myTime := mTime.MsToTime(val.UTime, "0")
			NewBalance := AccountBalance{
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
