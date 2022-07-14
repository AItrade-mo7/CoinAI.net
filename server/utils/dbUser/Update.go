package dbUser

import (
	"Hunter.net/server/global/apiType"
	"Hunter.net/server/global/dbType"
	"github.com/EasyGolang/goTools/mJson"
	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
)

func (dbObj *AccountType) Update() apiType.UserInfo {
	db := dbObj.DB
	var result dbType.AccountTable
	FK := bson.D{{
		Key:   "UserID",
		Value: dbObj.UserID,
	}}
	db.Table.FindOne(db.Ctx, FK).Decode(&result)
	dbObj.AccountData = result

	var userinfoData apiType.UserInfo
	jsonStr := mJson.ToJson(dbObj.AccountData)
	jsoniter.Unmarshal(jsonStr, &userinfoData)

	return userinfoData
}
