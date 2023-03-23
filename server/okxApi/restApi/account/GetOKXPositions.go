package account

import (
	"fmt"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/global/dbType"
	"github.com/EasyGolang/goTools/mFile"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mStr"
	jsoniter "github.com/json-iterator/go"
)

type PositionsData struct {
	AvgPx       string `bson:"AvgPx"`       // 开仓均价
	CTime       string `bson:"CTime"`       // 持仓创建时间
	Ccy         string `bson:"Ccy"`         // 币种
	InstID      string `bson:"InstID"`      // InstID
	InstType    string `bson:"InstType"`    // SWAP
	Interest    string `bson:"Interest"`    // 利息
	Last        string `bson:"Last"`        // 当前最新成交价
	Lever       string `bson:"Lever"`       // 杠杆倍数
	LiqPx       string `bson:"LiqPx"`       // 预估强平价格
	MarkPx      string `bson:"MarkPx"`      // 标记价格
	MgnRatio    string `bson:"MgnRatio"`    // 保证金率
	NotionalUsd string `bson:"NotionalUsd"` // 持仓数量
	Pos         string `bson:"Pos"`         // 持仓数量
	UTime       string `bson:"UTime"`       // 更新时间
	Upl         string `bson:"Upl"`         // 未实现收益
	UplRatio    string `bson:"UplRatio"`    // 未实现收益率
	Imr         string `bson:"Imr"`         // 初始保证金
}

// 查看持仓信息
func GetOKXPositions(OKXKey dbType.OkxKeyType) (resData []PositionsData, resErr error) {
	resData = []PositionsData{}
	resErr = nil

	if len(OKXKey.ApiKey) < 10 {
		resErr = fmt.Errorf("account.GetOKXPositions OKXKey.ApiKey 不能为空 Name:%+v", OKXKey.Name)
		global.LogErr(resErr)
		return
	}

	res, err := mOKX.FetchOKX(mOKX.OptFetchOKX{
		Path:   "/api/v5/account/positions",
		Method: "GET",
		OKXKey: mOKX.TypeOkxKey{
			ApiKey:     OKXKey.ApiKey,
			SecretKey:  OKXKey.SecretKey,
			Passphrase: OKXKey.Passphrase,
		},
	})
	if err != nil {
		resErr = fmt.Errorf("account.GetOKXPositions1 %+v Name:%+v", mStr.ToStr(err), OKXKey.Name)
		global.LogErr(resErr)
		return
	}
	var resObj mOKX.TypeReq
	jsoniter.Unmarshal(res, &resObj)
	if resObj.Code != "0" {
		resErr = fmt.Errorf("account.GetOKXPositions2 Name:%+v", OKXKey.Name)
		global.LogErr(resErr)
		return
	}

	var Data []PositionsData
	jsoniter.Unmarshal(mJson.ToJson(resObj.Data), &Data)

	Positions_file := mStr.Join(config.Dir.JsonData, "/Positions.json")
	mFile.Write(Positions_file, string(res))

	resData = Data

	return
}
