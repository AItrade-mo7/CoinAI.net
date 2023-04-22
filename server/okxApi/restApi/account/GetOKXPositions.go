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

// 查看持仓信息
func GetOKXPositions(OKXKey dbType.OkxKeyType) (resData []dbType.PositionsData, resErr error) {
	resData = []dbType.PositionsData{}
	resErr = nil

	if len(OKXKey.ApiKey) < 10 {
		resErr = fmt.Errorf("account.GetOKXPositions OKXKey.ApiKey 不能为空 Name:%+v", OKXKey.Name)
		global.LogErr("该错误已同步至用户邮箱", resErr)
		LogErr(OKXKey, resErr)
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
	// 打印接口日志
	global.OKXLogo.Println("account.GetOKXPositions",
		err,
		mStr.ToStr(res),
		OKXKey.Name,
	)

	if err != nil {
		resErr = fmt.Errorf("account.GetOKXPositions1 %+v Name:%+v", mStr.ToStr(err), OKXKey.Name)
		global.LogErr("该错误已同步至用户邮箱", resErr)
		LogErr(OKXKey, resErr)
		return
	}
	var resObj mOKX.TypeReq
	jsoniter.Unmarshal(res, &resObj)
	if resObj.Code != "0" {
		resErr = fmt.Errorf("account.GetOKXPositions2 Name:%+v", OKXKey.Name)
		global.LogErr("该错误已同步至用户邮箱", resErr)
		LogErr(OKXKey, resErr)
		return
	}

	var Data []dbType.PositionsData
	jsoniter.Unmarshal(mJson.ToJson(resObj.Data), &Data)

	Positions_file := mStr.Join(config.Dir.JsonData, "/Positions.json")
	mFile.Write(Positions_file, string(res))

	resData = Data

	return
}
