package account

import (
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mOKX"
	jsoniter "github.com/json-iterator/go"
)

type BalanceReq []struct {
	Details     []BalanceDetails `json:"details"`
	Imr         string           `json:"imr"`
	IsoEq       string           `json:"isoEq"`
	MgnRatio    string           `json:"mgnRatio"`
	Mmr         string           `json:"mmr"`
	TotalEq     string           `json:"totalEq"`
	AdjEq       string           `json:"adjEq"`
	NotionalUsd string           `json:"notionalUsd"`
	OrdFroz     string           `json:"ordFroz"`
	UTime       string           `json:"uTime"`
}
type BalanceDetails struct {
	Eq            string `json:"eq"`
	FixedBal      string `json:"fixedBal"`
	IsoEq         string `json:"isoEq"`
	OrdFrozen     string `json:"ordFrozen"`
	StgyEq        string `json:"stgyEq"`
	AvailBal      string `json:"availBal"`
	CrossLiab     string `json:"crossLiab"`
	MaxLoan       string `json:"maxLoan"`
	MgnRatio      string `json:"mgnRatio"`
	NotionalLever string `json:"notionalLever"`
	SpotInUseAmt  string `json:"spotInUseAmt"`
	Twap          string `json:"twap"`
	IsoLiab       string `json:"isoLiab"`
	IsoUpl        string `json:"isoUpl"`
	Liab          string `json:"liab"`
	UTime         string `json:"uTime"`
	AvailEq       string `json:"availEq"`
	CashBal       string `json:"cashBal"`
	EqUsd         string `json:"eqUsd"`
	FrozenBal     string `json:"frozenBal"`
	Interest      string `json:"interest"`
	Upl           string `json:"upl"`
	UplLiab       string `json:"uplLiab"`
	Ccy           string `json:"ccy"`
	DisEq         string `json:"disEq"`
}

func GetOKXBalance(ApiKey mOKX.TypeOkxKey) (resData BalanceReq) {
	res, err := mOKX.FetchOKX(mOKX.OptFetchOKX{
		Path:   "/api/v5/account/balance",
		Method: "GET",
		OKXKey: ApiKey,
	})
	if err != nil {
		return nil
	}

	var resObj mOKX.TypeReq
	jsoniter.Unmarshal(res, &resObj)

	if resObj.Code != "0" {
		return nil
	}
	var Data BalanceReq
	jsoniter.Unmarshal(mJson.ToJson(resObj.Data), &Data)

	if len(Data) < 1 {
		return nil
	}

	if len(Data[0].UTime) != 13 {
		return nil
	}

	return Data
}
