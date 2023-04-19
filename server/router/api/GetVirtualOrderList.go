package api

import (
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/global/dbType"
	"CoinAI.net/server/okxInfo"
	"CoinAI.net/server/router/result"
	"CoinAI.net/server/utils/dbSearch"
	"github.com/EasyGolang/goTools/mFiber"
	"github.com/EasyGolang/goTools/mMongo"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/gofiber/fiber/v2"
)

func GetVirtualOrderList(c *fiber.Ctx) error {
	var json dbSearch.FindParam
	mFiber.Parser(c, &json)

	HunterName := mStr.ToStr(json.Query["HunterName"])

	Hunter := okxInfo.HunterData{}
	for key, item := range okxInfo.NowHunterData {
		if key == HunterName {
			Hunter = item
		}
	}
	if len(Hunter.HunterName) < 1 {
		return c.JSON(result.Fail.WithMsg("该Hunter不存在!"))
	}

	db := mMongo.New(mMongo.Opt{
		UserName: config.SysEnv.MongoUserName,
		Password: config.SysEnv.MongoPassword,
		Address:  config.SysEnv.MongoAddress,
		DBName:   "AIServe",
	}).Connect().Collection("CoinOrder")
	defer db.Close()

	// 构建搜索参数

	resCur, err := dbSearch.GetCursor(dbSearch.CurOpt{
		Param: json,
		DB:    db,
	})
	if err != nil {
		return c.JSON(result.ErrDB.WithMsg(err))
	}

	// 提取数据
	var OrderArr []any
	for resCur.Cursor.Next(db.Ctx) {
		var result dbType.CoinOrderTable
		resCur.Cursor.Decode(&result)
		OrderArr = append(OrderArr, result)
	}

	// 生成返回数据
	returnData := resCur.GenerateData(OrderArr)

	return c.JSON(result.Succeed.WithData(returnData))
}
