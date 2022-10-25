package api

import (
	"fmt"

	"CoinAI.net/server/global/config"
	"CoinAI.net/server/okxApi"
	"CoinAI.net/server/okxApi/restApi/account"
	"CoinAI.net/server/router/middle"
	"CoinAI.net/server/router/result"
	"github.com/EasyGolang/goTools/mFiber"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/gofiber/fiber/v2"
)

type AccountDetailParam struct {
	Index int
}

type AccountDetail struct {
	Positions []account.PositionsData
	Balance   []account.AccountBalance
}

func GetAccountDetail(c *fiber.Ctx) error {
	// 在这里请求数据
	var json HandleKeyParam
	mFiber.Parser(c, &json)

	UserID, err := middle.TokenAuth(c)
	if err != nil {
		return c.JSON(result.ErrToken.WithData(mStr.ToStr(err)))
	}

	ApiKeyList := config.AppEnv.ApiKeyList

	var ListErr error
	OkxKey := mOKX.TypeOkxKey{}
	for key, val := range ApiKeyList {
		if key == json.Index {
			if val.UserID != UserID {
				ListErr = fmt.Errorf("无权操作")
				break
			}
			OkxKey = val
		}
	}

	if ListErr != nil {
		return c.JSON(result.Fail.WithMsg(ListErr))
	}

	// 新建账户
	OKXAccount, err := okxApi.NewAccount(okxApi.AccountParam{
		OkxKey: OkxKey,
	})
	if err != nil {
		return c.JSON(result.ErrOKXAccount.WithMsg(err))
	}

	err = OKXAccount.GetPositions()
	if err != nil {
		return c.JSON(result.ErrOKXAccount.WithMsg(err))
	}
	err = OKXAccount.GetBalance()
	if err != nil {
		return c.JSON(result.ErrOKXAccount.WithMsg(err))
	}

	return c.JSON(result.Succeed.WithData(AccountDetail{
		Positions: OKXAccount.Positions,
		Balance:   OKXAccount.Balance,
	}))
}
